package etcd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Observer interface {
	Update([]string)
}

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	mu           sync.RWMutex
	format       string
	serverList   map[string]string // 最新服务列表：保存相同prefix的服务, key-value
	ipWeightList map[string]string // 上游服务的ip列表信息
	client       *clientv3.Client
	observers    []Observer
}

func NewServiceDiscovery(endpoints []string, format string, ipWeightList map[string]string) (*ServiceDiscovery, error) {

	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}
	c, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	// github issue: https://github.com/etcd-io/etcd/issues/9877
	// clientv3: clientv3.New() won't return error when no endpoint is available
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()
	// 连接是否超时
	_, err = c.Status(ctx, cfg.Endpoints[0])
	if err != nil {
		return nil, err
	}

	discovery := &ServiceDiscovery{
		client:       c,
		serverList:   make(map[string]string),
		format:       format,
		ipWeightList: ipWeightList,
	}

	return discovery, nil
}

func (discovery *ServiceDiscovery) AddObserver(observer Observer) {
	discovery.observers = append(discovery.observers, observer)
}

// WatchService 首次初始化本地的服务列表+后续监听相同前缀的服务列表
func (discovery *ServiceDiscovery) WatchService(prefix string) ([]string, error) {
	resp, err := discovery.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range resp.Kvs {
		if kv.Value != nil {
			discovery.SetServiceList(string(kv.Key), string(kv.Value))
		}
	}

	addrs := discovery.GetServiceList()
	go discovery.watcher(prefix)
	return addrs, nil
}

func (discovery *ServiceDiscovery) watcher(prefix string) {
	watchChan := discovery.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				discovery.SetServiceList(string(event.Kv.Key), string(event.Kv.Value))
			case mvccpb.DELETE:
				discovery.DelServiceList(string(event.Kv.Key))
			}
		}
		discovery.notifyAll()
	}
}

func (discovery *ServiceDiscovery) notifyAll() {
	n := len(discovery.observers)
	if n == 0 {
		return
	}
	addrs := discovery.GetServiceList()
	for i := n - 1; i >= 0; i-- {
		discovery.observers[i].Update(addrs)
	}
}

func (discovery *ServiceDiscovery) GetServiceList() []string {
	discovery.mu.RLock()
	defer discovery.mu.RUnlock()

	addrs := make([]string, 0, len(discovery.serverList))
	for _, addr := range discovery.serverList {
		if _, ok := discovery.ipWeightList[addr]; ok {
			addrs = append(addrs, addr)
		}
	}

	list := make([]string, 0, len(addrs))
	// 构造 [ ip:port,weight ]
	for _, addr := range addrs {
		weight, ok := discovery.ipWeightList[addr]
		if !ok {
			weight = "50" //默认weight
		}
		list = append(list, fmt.Sprintf(discovery.format, addr)+","+weight)
	}

	return list
}

func (discovery *ServiceDiscovery) SetServiceList(key, value string) {
	discovery.mu.Lock()
	discovery.serverList[key] = value
	discovery.mu.Unlock()
}

func (discovery *ServiceDiscovery) DelServiceList(key string) {
	discovery.mu.Lock()
	delete(discovery.serverList, key)
	discovery.mu.Unlock()
}

func (discovery *ServiceDiscovery) Close() error {
	return discovery.client.Close()
}
