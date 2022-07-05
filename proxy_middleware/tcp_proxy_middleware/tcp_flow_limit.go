package tcpproxymiddleware

import (
	"fmt"
	"strings"

	"github.com/josexy/gw/handler"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/constants"
	"github.com/josexy/gw/proxy_router/tcp_proxy_router/router"
)

func TCPFlowLimitMiddleware() router.TcpHandlerFunc {
	return func(c *router.TcpRouterContext) {
		service := c.Get("service")
		if service == nil {
			c.Conn.Write([]byte("service not found"))
			c.Abort()
			return
		}

		serviceDetail := service.(*model.ServiceDetail)
		if serviceDetail.AccessControl.ServiceFlowLimit > 0 {
			serviceLimiter, err := handler.FlowLimiterHandler.GetLimiter(
				serviceDetail.Info.ServiceName, float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				c.Conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}

			if !serviceLimiter.Allow() {
				c.Conn.Write([]byte(fmt.Sprintf("service flow limit %v",
					serviceDetail.AccessControl.ServiceFlowLimit)))
				c.Abort()
				return
			}
		}

		splits := strings.Split(c.Conn.RemoteAddr().String(), ":")
		var clientIP string
		if len(splits) == 2 {
			clientIP = splits[0]
		}

		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := handler.FlowLimiterHandler.GetLimiter(
				constants.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP,
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				c.Conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				c.Conn.Write([]byte(fmt.Sprintf("%v flow limit %v",
					clientIP, serviceDetail.AccessControl.ClientIPFlowLimit)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
