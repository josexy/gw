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

func HTTPBlackListMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exist := ctx.Get("service")
		if !exist {
			api.ResponseJsonErrorCode(ctx, 2001, errors.New("service not found"))
			ctx.Abort()
			return
		}
		serviceDetail := service.(*model.ServiceDetail)

		var whiteIpList []string
		var blackIpList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			whiteIpList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}

		// 1. 白名单不为空，只处理白名单
		// 2. 白名单为空，则处理黑名单列表
		if serviceDetail.AccessControl.EnableAuth == 1 && len(whiteIpList) == 0 && len(blackIpList) > 0 {
			// 如果在黑名单中，则拒绝访问
			if util.FindInList(blackIpList, ctx.ClientIP()) {
				api.ResponseJsonErrorCode(ctx, 2002,
					errors.New(fmt.Sprintf("%s in black ip list", ctx.ClientIP())))
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
