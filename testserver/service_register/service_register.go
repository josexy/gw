package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/pkg/constants"
	"github.com/josexy/gw/pkg/etcd"
)

// 模拟服务注册
func main() {
	if len(os.Args) == 1 {
		logx.Debug("need service register address")
		return
	}
	addr := os.Args[1]

	endpoints := []string{"127.0.0.1:2379"}
	prefix := constants.EtcdPrefix

	service, err := etcd.NewServiceRegister(endpoints, 5)
	if err != nil {
		logx.Error("connect server register: %v", err)
		return
	}

	key := prefix + addr
	logx.Debug("register service to etcd: %s", key)
	_ = service.RegisterService(key, addr)

	// 结束程序前取消服务注册
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT)
	<-interrupt

	_ = service.Close()
}
