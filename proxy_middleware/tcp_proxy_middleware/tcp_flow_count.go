package tcpproxymiddleware

import (
	"github.com/josexy/gw/handler"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/constants"
	"github.com/josexy/gw/proxy_router/tcp_proxy_router/router"
)

func TCPFlowCountMiddleware() router.TcpHandlerFunc {
	return func(c *router.TcpRouterContext) {
		service := c.Get("service")
		if service == nil {
			c.Conn.Write([]byte("service not found"))
			c.Abort()
			return
		}

		serviceDetail := service.(*model.ServiceDetail)
		totalCounter, err := handler.FlowCounterHandler.GetCounter(constants.FlowTotal)
		if err != nil {
			c.Conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		totalCounter.Increase()

		serviceName := constants.FlowServicePrefix + serviceDetail.Info.ServiceName
		serviceCounter, err := handler.FlowCounterHandler.GetCounter(serviceName)
		if err != nil {
			c.Conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		serviceCounter.Increase()
		c.Next()
	}
}
