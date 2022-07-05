package handler

import (
	"net/http"
	"net/http/httputil"

	"github.com/josexy/gw/model"
	"github.com/josexy/gw/proxy_router/tcp_proxy_router/router"
	reverseproxy "github.com/josexy/gw/reverse_proxy"
	"github.com/josexy/gw/reverse_proxy/load_balance"
	"google.golang.org/grpc"
)

var ReverseProxyHandler *ReverseProxy

type ReverseProxy struct {
	httpHandler MapBaseHandler[*httputil.ReverseProxy]
	tcpHandler  MapBaseHandler[*reverseproxy.TcpReverseProxy]
	grpcHandler MapBaseHandler[grpc.StreamHandler]
}

func init() {
	ReverseProxyHandler = NewReverseProxy()
}

func NewReverseProxy() *ReverseProxy {
	return &ReverseProxy{
		httpHandler: NewMapBaseHandler[*httputil.ReverseProxy](),
		tcpHandler:  NewMapBaseHandler[*reverseproxy.TcpReverseProxy](),
		grpcHandler: NewMapBaseHandler[grpc.StreamHandler](),
	}
}

func (p *ReverseProxy) GetHTTPReverseProxy(serviceDetail *model.ServiceDetail,
	lb load_balance.LoadBalance, transport *http.Transport) (*httputil.ReverseProxy, error) {
	p.httpHandler.RLock()
	if proxy, ok := p.httpHandler.Cache[serviceDetail.Info.ServiceName]; ok {
		p.httpHandler.RUnlock()
		return proxy, nil
	}
	p.httpHandler.RUnlock()

	// 创建反向代理对象
	proxy := reverseproxy.NewHTTPLoadBalanceReverseProxy(lb, transport)

	p.httpHandler.Lock()
	defer p.httpHandler.Unlock()
	p.httpHandler.List = append(p.httpHandler.List, proxy)
	p.httpHandler.Cache[serviceDetail.Info.ServiceName] = proxy

	return proxy, nil
}

func (p *ReverseProxy) GetTCPReverseProxy(serviceDetail *model.ServiceDetail,
	c *router.TcpRouterContext, lb load_balance.LoadBalance) (*reverseproxy.TcpReverseProxy, error) {
	p.tcpHandler.RLock()
	if proxy, ok := p.tcpHandler.Cache[serviceDetail.Info.ServiceName]; ok {
		p.tcpHandler.RUnlock()
		return proxy, nil
	}
	p.tcpHandler.RUnlock()

	// 创建反向代理对象
	proxy := reverseproxy.NewTCPLoadBalanceReverseProxy(c, lb)

	p.tcpHandler.Lock()
	defer p.tcpHandler.Unlock()
	p.tcpHandler.List = append(p.tcpHandler.List, proxy)
	p.tcpHandler.Cache[serviceDetail.Info.ServiceName] = proxy

	return proxy, nil
}

func (p *ReverseProxy) GetGRPCReverseProxy(serviceDetail *model.ServiceDetail,
	lb load_balance.LoadBalance) (grpc.StreamHandler, error) {
	p.grpcHandler.RLock()
	if proxy, ok := p.grpcHandler.Cache[serviceDetail.Info.ServiceName]; ok {
		p.grpcHandler.RUnlock()
		return proxy, nil
	}
	p.grpcHandler.RUnlock()

	// 创建反向代理对象
	proxy := reverseproxy.NewGRPCLoadBalanceReverseProxy(lb)

	p.grpcHandler.Lock()
	defer p.grpcHandler.Unlock()
	p.grpcHandler.List = append(p.grpcHandler.List, proxy)
	p.grpcHandler.Cache[serviceDetail.Info.ServiceName] = proxy

	return proxy, nil
}
