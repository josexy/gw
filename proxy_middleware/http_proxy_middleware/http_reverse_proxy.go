package httpproxymiddleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/api"
	"github.com/josexy/gw/handler"
	"github.com/josexy/gw/model"
)

func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service, exist := ctx.Get("service")
		if !exist {
			api.ResponseJsonErrorCode(ctx, 2001, errors.New("service not found"))
			ctx.Abort()
			return
		}
		serviceDetail := service.(*model.ServiceDetail)
		lb, err := handler.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			api.ResponseJsonErrorCode(ctx, 2002, err)
			ctx.Abort()
			return
		}
		// 获取连接对象
		transporter, err := handler.TransporterHandler.GetTransporter(serviceDetail)
		if err != nil {
			api.ResponseJsonErrorCode(ctx, 2003, err)
			ctx.Abort()
			return
		}
		// 负载均衡反向代理
		proxy, err := handler.ReverseProxyHandler.GetHTTPReverseProxy(serviceDetail, lb, transporter)
		if err != nil {
			api.ResponseJsonErrorCode(ctx, 2004, err)
			ctx.Abort()
			return
		}
		// 转发请求
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
		// 不再传递中间件
		ctx.Abort()
	}
}
