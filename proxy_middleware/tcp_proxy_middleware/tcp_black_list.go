package tcpproxymiddleware

import (
	"fmt"
	"strings"

	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/util"
	"github.com/josexy/gw/proxy_router/tcp_proxy_router/router"
)

func TCPBlackListMiddleware() router.TcpHandlerFunc {
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

		var whiteIpList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			whiteIpList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		var blackIpList []string
		if serviceDetail.AccessControl.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}

		if serviceDetail.AccessControl.EnableAuth == 1 && len(whiteIpList) == 0 && len(blackIpList) > 0 {
			if util.FindInList(blackIpList, clientIp) {
				c.Conn.Write([]byte(fmt.Sprintf("%s in black ip list", clientIp)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
