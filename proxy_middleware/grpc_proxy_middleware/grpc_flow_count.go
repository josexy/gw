package grpcproxymiddleware

import (
	"github.com/josexy/gw/handler"
	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/constants"
	"google.golang.org/grpc"
)

func GrpcFlowCountMiddleware(serviceDetail *model.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, streamHandler grpc.StreamHandler) error {
		totalCounter, err := handler.FlowCounterHandler.GetCounter(constants.FlowTotal)
		if err != nil {
			return err
		}
		// 增加计数
		totalCounter.Increase()
		serviceCounter, err := handler.FlowCounterHandler.GetCounter(constants.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			return err
		}
		// 增加计数
		serviceCounter.Increase()

		if err = streamHandler(srv, ss); err != nil {
			logx.Debug("GRPC middleware call handler err: %v", err)
			return err
		}
		return nil
	}
}
