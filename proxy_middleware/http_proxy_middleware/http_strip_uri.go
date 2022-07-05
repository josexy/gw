package httpproxymiddleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/api"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/constants"
)

// HTTPStripUriMiddleware 将前缀匹配的uri删除
func HTTPStripUriMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exist := ctx.Get("service")
		if !exist {
			api.ResponseJsonErrorCode(ctx, 2001, errors.New("service not found"))
			ctx.Abort()
			return
		}
		serviceDetail := service.(*model.ServiceDetail)
		/*
			1. 处理前：http://127.0.0.1:8080/test_strip_uri/abcd
			删除匹配的http前缀uri规则: /test_strip_uri 保留后面的 /abcd
			2. 处理后(经过负载均衡)：http://127.0.0.1:2003/abcd
		*/
		// 开启 strip-uri 且 是前缀匹配
		if serviceDetail.HTTPRule.NeedStripUri == 1 &&
			serviceDetail.HTTPRule.RuleType == constants.HTTPRuleTypePrefixURL {
			// 前缀uri删除
			ctx.Request.URL.Path = strings.Replace(ctx.Request.URL.Path, serviceDetail.HTTPRule.Rule, "", 1)
		}
		ctx.Next()
	}
}
