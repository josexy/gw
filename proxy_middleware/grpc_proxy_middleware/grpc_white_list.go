package grpcproxymiddleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func GrpcWhiteListMiddleware(serviceDetail *model.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, streamHandler grpc.StreamHandler) error {
		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("not found grpc ip address from context")
		}
		peerAddr := peerCtx.Addr.String()
		// ip:port
		parts := strings.Split(peerAddr, ":")
		var clientIp string
		if len(parts) == 2 {
			clientIp = parts[0]
		}

		var ipList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		if serviceDetail.AccessControl.EnableAuth == 1 && len(ipList) > 0 {
			if !util.FindInList(ipList, clientIp) {
				return errors.New(fmt.Sprintf("%s not in white ip list", clientIp))
			}
		}

		if err := streamHandler(srv, ss); err != nil {
			logx.Debug("GRPC middleware call handler err: %v", err)
			return err
		}
		return nil
	}
}
