package tcpproxyrouter

import (
	"context"
	"fmt"

	"github.com/josexy/gw/handler"
	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/model"
	tcpproxymiddleware "github.com/josexy/gw/proxy_middleware/tcp_proxy_middleware"
	"github.com/josexy/gw/proxy_router/tcp_proxy_router/endpoint"
	"github.com/josexy/gw/proxy_router/tcp_proxy_router/router"
)

var tcpServerList []*endpoint.TcpServer

func TcpServerRun() {
	tcpServiceList := handler.ServiceManagerHandler.GetTcpServiceList()
	for _, service := range tcpServiceList {
		go func(serviceDetail *model.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.TCPRule.Port)

			lb, err := handler.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				logx.Debug("get load balancer err: %v", err)
				return
			}

			r := router.NewTcpRouter()
			r.Use(
				tcpproxymiddleware.TCPFlowCountMiddleware(),
				tcpproxymiddleware.TCPFlowLimitMiddleware(),

				tcpproxymiddleware.TCPWhiteListMiddleware(),
				tcpproxymiddleware.TCPBlackListMiddleware(),
			)

			r.Serve(func(c *router.TcpRouterContext) endpoint.TcpHandler {
				reverseProxy, _ := handler.ReverseProxyHandler.GetTCPReverseProxy(serviceDetail, c, lb)
				return reverseProxy
			})

			baseCtx := context.WithValue(context.Background(), "service", serviceDetail)

			tcpServer := &endpoint.TcpServer{
				BaseContext: baseCtx,
				Addr:        addr,
				Handler:     r,
			}
			tcpServerList = append(tcpServerList, tcpServer)

			logx.Info("TcpServerRun: run tcp server: %v", addr)
			if err = tcpServer.ListenAndServe(); err != nil && err != endpoint.ErrServerClosed {
				logx.Fatal("TcpServerRun: ListenAndServe err: %v", err)
			}
		}(service)
	}
}

func TcpServerStop() {
	for _, tcpService := range tcpServerList {
		_ = tcpService.Close()
		logx.Warn("close and stop tcp server: %v", tcpService.Addr)
	}
}
