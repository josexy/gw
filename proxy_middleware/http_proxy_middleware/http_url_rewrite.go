package httpproxymiddleware

import (
	"errors"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/api"
	"github.com/josexy/gw/model"
)

// HTTPURLRewriteMiddleware URL重写
func HTTPURLRewriteMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exist := ctx.Get("service")
		if !exist {
			api.ResponseJsonErrorCode(ctx, 2001, errors.New("service not found"))
			ctx.Abort()
			return
		}
		serviceDetail := service.(*model.ServiceDetail)
		/*
			可使用 "分组引用符" 获取相应的分组内容
			/aaa/(.*) /bbb/$1
			前：/test_strip_uri/aaa/xxx
			后：/test_strip_uri/bbb/xxx
		*/
		for _, item := range strings.Split(serviceDetail.HTTPRule.UrlRewrite, ",") {
			items := strings.Split(item, " ")
			if len(items) != 2 {
				continue
			}
			rgx, err := regexp.Compile(items[0])
			if err != nil {
				continue
			}
			rewriteUrl := rgx.ReplaceAllString(ctx.Request.URL.Path, items[1])
			ctx.Request.URL.Path = rewriteUrl
		}
		ctx.Next()
	}
}
