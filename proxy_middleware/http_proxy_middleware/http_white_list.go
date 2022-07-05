package httpproxymiddleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/api"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/util"
)

func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exist := ctx.Get("service")
		if !exist {
			api.ResponseJsonErrorCode(ctx, 2001, errors.New("service not found"))
			ctx.Abort()
			return
		}
		serviceDetail := service.(*model.ServiceDetail)

		var ipList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.EnableAuth == 1 && len(ipList) > 0 {
			// 如果不在白名单中，则拒绝访问
			if !util.FindInList(ipList, ctx.ClientIP()) {
				api.ResponseJsonErrorCode(ctx, 2002,
					errors.New(fmt.Sprintf("%s not in white ip list", ctx.ClientIP())))
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
