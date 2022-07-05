package tcpproxymiddleware

import (
	"fmt"
	"strings"

	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/util"
	"github.com/josexy/gw/proxy_router/tcp_proxy_router/router"
)

func TCPWhiteListMiddleware() router.TcpHandlerFunc {
	return func(c *router.TcpRouterContext) {
		service := c.Get("service")
		if service == nil {
			c.Conn.Write([]byte("service not found"))
			c.Abort()
			return
		}

		serviceDetail := service.(*model.ServiceDetail)

		// ip:port
		parts := strings.Split(c.Conn.RemoteAddr().String(), ":")
		var clientIp string
		if len(parts) == 2 {
			clientIp = parts[0]
		}

		var ipList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		if serviceDetail.AccessControl.EnableAuth == 1 && len(ipList) > 0 {
			if !util.FindInList(ipList, clientIp) {
				c.Conn.Write([]byte(fmt.Sprintf("%s not in white ip list", clientIp)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
