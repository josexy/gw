package etcd

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// ServiceRegister 服务注册
type ServiceRegister struct {
	Client        *clientv3.Client
	leaseID       clientv3.LeaseID // 租约ID
	cancelFunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse // 租约健康检查
}

func NewServiceRegister(endpoints []string, leaseTTL int64) (*ServiceRegister, error) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	lease := clientv3.NewLease(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	leaseResp, err := lease.Grant(ctx, leaseTTL)
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		cancelFunc()
		return nil, err
	}

	service := &ServiceRegister{
		Client:        c,
		leaseID:       leaseResp.ID,
		cancelFunc:    cancelFunc,
		keepAliveChan: leaseRespChan,
	}
	go service.listenLeaseRespChan()
	return service, nil
}

func (service *ServiceRegister) listenLeaseRespChan() {
	for resp := range service.keepAliveChan {
		_ = resp
	}
}

func (service *ServiceRegister) RegisterService(key, value string) error {
	_, err := service.Client.Put(context.Background(), key, value, clientv3.WithLease(service.leaseID))
	return err
}

func (service *ServiceRegister) RevokeService() error {
	service.cancelFunc()
	time.Sleep(1 * time.Second)
	_, err := service.Client.Revoke(context.Background(), service.leaseID)
	return err
}

func (service *ServiceRegister) Close() error {
	if err := service.RevokeService(); err != nil {
		return err
	}
	return service.Client.Close()
}
