package grpcproxyrouter

import (
	"fmt"
	"net"

	"github.com/josexy/gw/handler"
	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/model"
	grpcproxymiddleware "github.com/josexy/gw/proxy_middleware/grpc_proxy_middleware"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
)

var grpcServerList []*GrpcServer

type GrpcServer struct {
	Addr string
	*grpc.Server
}

func GrpcServerRun() {
	// 获取grpc服务列表
	grpcServiceList := handler.ServiceManagerHandler.GetGrpcServiceList()
	for _, service := range grpcServiceList {
		go func(serviceDetail *model.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.GRPCRule.Port)

			// 为每一个服务创建独立的负载均衡器
			lb, err := handler.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				logx.Debug("get load balancer err: %v", err)
				return
			}
			// 启动grpc服务器
			tcpSvr, err := net.Listen("tcp", addr)
			if err != nil {
				logx.Debug("run tcp grpc server err: %v", err)
				return
			}

			reverseProxy, _ := handler.ReverseProxyHandler.GetGRPCReverseProxy(serviceDetail, lb)
			grpcSvr := grpc.NewServer(
				grpc.ChainStreamInterceptor(
					grpcproxymiddleware.GrpcFlowCountMiddleware(serviceDetail),
					grpcproxymiddleware.GrpcFlowLimitMiddleware(serviceDetail),

					grpcproxymiddleware.GrpcWhiteListMiddleware(serviceDetail),
					grpcproxymiddleware.GrpcBlackListMiddleware(serviceDetail),

					grpcproxymiddleware.GrpcHeaderTransferMiddleware(serviceDetail),
				),
				grpc.CustomCodec(proxy.Codec()),
				// 由 grpcHandler 处理未注册的服务
				grpc.UnknownServiceHandler(reverseProxy),
			)

			grpcServer := &GrpcServer{
				Addr:   addr,
				Server: grpcSvr,
			}

			grpcServerList = append(grpcServerList, grpcServer)

			logx.Info("GrpcServerRun run grpc server: %v", addr)
			if err = grpcSvr.Serve(tcpSvr); err != nil {
				logx.Fatal("GrpcServerRun: Grpc Serve err: %v", err)
			}
		}(service)
	}
}

func GrpcServerStop() {
	for _, server := range grpcServerList {
		// 优雅退出grpc服务器
		server.GracefulStop()
		logx.Warn("close and stop grpc server: %v", server.Addr)
	}
}
