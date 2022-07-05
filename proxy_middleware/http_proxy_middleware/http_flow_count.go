package httpproxymiddleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/api"
	"github.com/josexy/gw/handler"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/constants"
)

func HTTPFlowCountMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exist := ctx.Get("service")
		if !exist {
			api.ResponseJsonErrorCode(ctx, 2001, errors.New("service not found"))
			ctx.Abort()
			return
		}
		serviceDetail := service.(*model.ServiceDetail)

		totalCounter, err := handler.FlowCounterHandler.GetCounter(constants.FlowTotal)
		if err != nil {
			api.ResponseJsonErrorCode(ctx, 4001, err)
			ctx.Abort()
			return
		}
		totalCounter.Increase()

		serviceName := constants.FlowServicePrefix + serviceDetail.Info.ServiceName
		serviceCounter, err := handler.FlowCounterHandler.GetCounter(serviceName)
		if err != nil {
			api.ResponseJsonErrorCode(ctx, 4001, err)
			ctx.Abort()
			return
		}
		serviceCounter.Increase()

		ctx.Next()
	}
}
