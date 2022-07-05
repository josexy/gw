package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/pkg/constants"
	"github.com/josexy/gw/pkg/etcd"
	"github.com/josexy/gw/proxy_router/tcp_proxy_router/endpoint"
)

type RealServer struct {
	Addr       string
	ServiceReg *etcd.ServiceRegister
}

func (server *RealServer) ServeTCP(ctx context.Context, conn net.Conn) {
	logx.Debug("ServeTCP client address: %v", conn.RemoteAddr().String())
	//io.Copy(conn, conn)
	conn.Write([]byte("hello world"))
	conn.Close()
}

func (server *RealServer) Run() {
	logx.Debug("start server at: %v", server.Addr)

	svr := endpoint.TcpServer{
		Addr:    server.Addr,
		Handler: server,
	}
	go func() {
		err := svr.ListenAndServe()
		if err != nil && err != endpoint.ErrServerClosed {
			fmt.Println(err)
			return
		}
	}()

	// 服务注册
	server.registerService(server.Addr)
}

func (server *RealServer) unregisterService() {
	logx.Debug("revoke service and close server: %s", server.Addr)
	if server.ServiceReg != nil {
		_ = server.ServiceReg.Close()
	}
}

func (server *RealServer) registerService(addr string) {

	endpoints := []string{"127.0.0.1:2379"}
	prefix := constants.EtcdPrefix

	go func() {
		service, err := etcd.NewServiceRegister(endpoints, 5)
		if err != nil {
			logx.Error("connect server register: %v", err)
			return
		}
		server.ServiceReg = service

		key := prefix + addr
		logx.Debug("register service to etcd: %s", key)
		_ = service.RegisterService(key, addr)
	}()
}

func main() {

	var port1, port2 int
	var err error

	if len(os.Args) == 1 {
		port1 = 2003
		port2 = 2004
	} else {
		port1, err = strconv.Atoi(os.Args[1])
		if err != nil {
			logx.Fatal("%v", err)
			return
		}
		port2, err = strconv.Atoi(os.Args[2])
		if err != nil {
			logx.Fatal("%v", err)
			return
		}
	}

	rs1 := &RealServer{Addr: fmt.Sprintf("127.0.0.1:%d", port1)}
	rs2 := &RealServer{Addr: fmt.Sprintf("127.0.0.1:%d", port2)}

	rs1.Run()
	rs2.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT)
	<-interrupt

	rs1.unregisterService()
	rs2.unregisterService()
}
