package httpproxymiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/api"
	"github.com/josexy/gw/handler"
)

func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		serviceDetail, err := handler.ServiceManagerHandler.HTTPAccessMode(ctx.Request.Host, ctx.Request.URL.Path)
		if err != nil {
			api.ResponseJsonErrorCode(ctx, 2001, err)
			ctx.Abort()
			return
		}
		ctx.Set("service", serviceDetail)
		ctx.Next()
	}
}
