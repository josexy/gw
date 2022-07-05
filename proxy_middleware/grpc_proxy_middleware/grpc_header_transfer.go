package grpcproxymiddleware

import (
	"errors"
	"strings"

	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GrpcHeaderTransferMiddleware(serviceDetail *model.ServiceDetail) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, streamHandler grpc.StreamHandler) error {
		// 获取上下文context传递过来的metadata数据
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}

		for _, item := range strings.Split(serviceDetail.GRPCRule.HeaderTransfer, ",") {
			items := strings.Split(item, " ")
			if len(items) != 3 {
				continue
			}
			if items[0] == "add" || items[0] == "edit" {
				md.Set(items[1], items[2])
			} else if items[0] == "del" {
				delete(md, items[1])
			}
		}

		if err := ss.SetHeader(md); err != nil {
			return err
		}
		if err := streamHandler(srv, ss); err != nil {
			logx.Error("GRPC middleware call handler err: %v", err)
			return err
		}
		return nil
	}
}
