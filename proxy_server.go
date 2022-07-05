package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/josexy/gw/handler"
	"github.com/josexy/gw/logx"
	grpcproxyrouter "github.com/josexy/gw/proxy_router/grpc_proxy_router"
	httpproxyrouter "github.com/josexy/gw/proxy_router/http_proxy_router"
	tcpproxyrouter "github.com/josexy/gw/proxy_router/tcp_proxy_router"
)

type ProxyServer struct{}

func (svr *ProxyServer) Run() {

	if err := handler.ServiceManagerHandler.LoadOnce(); err != nil {
		logx.Fatal("service manager load once err: %v", err)
	}

	go func() { httpproxyrouter.HttpServerRun() }()
	go func() { httpproxyrouter.HttpsServerRun() }()
	go func() { tcpproxyrouter.TcpServerRun() }()
	go func() { grpcproxyrouter.GrpcServerRun() }()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	httpproxyrouter.HttpServerStop()
	httpproxyrouter.HttpsServerStop()
	tcpproxyrouter.TcpServerStop()
	grpcproxyrouter.GrpcServerStop()
}
