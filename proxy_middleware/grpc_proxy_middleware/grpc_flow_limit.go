package grpcproxymiddleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/josexy/gw/handler"
	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/constants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func GrpcFlowLimitMiddleware(serviceDetail *model.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, streamHandler grpc.StreamHandler) error {

		if serviceDetail.AccessControl.ServiceFlowLimit > 0 {
			serviceLimiter, err := handler.FlowLimiterHandler.GetLimiter(serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				return err
			}

			// 限流
			if !serviceLimiter.Allow() {
				return errors.New(fmt.Sprintf("service flow limit %v", serviceDetail.AccessControl.ServiceFlowLimit))
			}
		}

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("not found grpc ip address from context")
		}

		peerAddr := peerCtx.Addr.String()
		splits := strings.Split(peerAddr, ":")
		var clientIP string
		if len(splits) == 2 {
			clientIP = splits[0]
		}

		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := handler.FlowLimiterHandler.GetLimiter(
				constants.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP,
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				return err
			}
			if !clientLimiter.Allow() {
				return errors.New(fmt.Sprintf("client ip flow limit %v",
					serviceDetail.AccessControl.ServiceFlowLimit))
			}
		}

		if err := streamHandler(srv, ss); err != nil {
			logx.Debug("GRPC middleware call handler err: %v", err)
			return err
		}
		return nil
	}
}
