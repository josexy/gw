package httpproxymiddleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/api"
	"github.com/josexy/gw/model"
)

func HTTPHeaderTransferMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exist := ctx.Get("service")
		if !exist {
			api.ResponseJsonErrorCode(ctx, 2001, errors.New("service not found"))
			ctx.Abort()
			return
		}
		serviceDetail := service.(*model.ServiceDetail)
		for _, item := range strings.Split(serviceDetail.HTTPRule.HeaderTransfer, ",") {
			items := strings.Split(item, " ")
			if len(items) != 3 {
				continue
			}
			if items[0] == "add" || items[0] == "edit" {
				ctx.Request.Header.Set(items[1], items[2])
			} else if items[0] == "del" {
				ctx.Request.Header.Del(items[1])
			}
		}
		ctx.Next()
	}
}
