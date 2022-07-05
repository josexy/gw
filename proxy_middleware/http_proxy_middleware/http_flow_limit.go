package httpproxymiddleware

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/api"
	"github.com/josexy/gw/handler"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/constants"
)

func HTTPFlowLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exist := ctx.Get("service")
		if !exist {
			api.ResponseJsonErrorCode(ctx, 2001, errors.New("service not found"))
			ctx.Abort()
			return
		}
		serviceDetail := service.(*model.ServiceDetail)

		if serviceDetail.AccessControl.ServiceFlowLimit > 0 {
			serviceLimiter, err := handler.FlowLimiterHandler.GetLimiter(
				serviceDetail.Info.ServiceName, float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				api.ResponseJsonErrorCode(ctx, 5001, err)
				ctx.Abort()
				return
			}

			if !serviceLimiter.Allow() {
				api.ResponseJsonErrorCode(ctx, 5002, errors.New(fmt.Sprintf("service flow limit %v",
					serviceDetail.AccessControl.ServiceFlowLimit)))
				ctx.Abort()
				return
			}
		}

		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := handler.FlowLimiterHandler.GetLimiter(
				constants.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+ctx.ClientIP(),
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				api.ResponseJsonErrorCode(ctx, 5004, err)
				ctx.Abort()
				return
			}
			if !clientLimiter.Allow() {
				api.ResponseJsonErrorCode(ctx, 5005, errors.New(fmt.Sprintf("%v flow limit %v",
					ctx.ClientIP(), serviceDetail.AccessControl.ClientIPFlowLimit)))
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
