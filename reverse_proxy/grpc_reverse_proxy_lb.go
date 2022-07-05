package reverseproxy

import (
	"context"

	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/reverse_proxy/load_balance"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func NewGRPCLoadBalanceReverseProxy(lb load_balance.LoadBalance) grpc.StreamHandler {

	director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {

		var clientIp string
		peerCtx, ok := peer.FromContext(ctx)
		if !ok {
			logx.Error("not found grpc ip address from context")
		} else {
			clientIp = peerCtx.Addr.String()
		}

		nextAddr, err := lb.Get(clientIp)
		if err != nil {
			logx.Error("the upstream address is invalid or not found")
		}

		conn, err := grpc.DialContext(ctx,
			nextAddr,
			grpc.WithCodec(proxy.Codec()),
			grpc.WithInsecure(),
		)
		// 转发grpc的metadata header头信息
		md, _ := metadata.FromIncomingContext(ctx)
		outCtx, _ := context.WithCancel(ctx)
		outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
		// 返回新的上下文
		return outCtx, conn, err
	}

	// grpc反向代理透明传输
	return proxy.TransparentHandler(director)
}
