package reverseproxy

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/proxy_router/tcp_proxy_router/router"
	"github.com/josexy/gw/reverse_proxy/load_balance"
)

func NewTCPLoadBalanceReverseProxy(c *router.TcpRouterContext, lb load_balance.LoadBalance) *TcpReverseProxy {

	director := func(proxy *TcpReverseProxy) {
		nextAddr, err := lb.Get(c.Conn.RemoteAddr().String())
		if err != nil {
			logx.Error("the upstream address is invalid or not found")
		}
		proxy.Addr = nextAddr
	}

	return &TcpReverseProxy{
		ctx:             c.Ctx,
		director:        director,
		KeepAlivePeriod: time.Second,
		DialTimeout:     time.Second,
	}
}

type TcpReverseProxy struct {
	ctx                  context.Context
	director             func(proxy *TcpReverseProxy)
	Addr                 string
	KeepAlivePeriod      time.Duration
	DialTimeout          time.Duration
	DialContext          func(ctx context.Context, network, address string) (net.Conn, error)
	OnDialError          func(src net.Conn, dstDialErr error)
	ProxyProtocolVersion int
}

func (dp *TcpReverseProxy) dialTimeout() time.Duration {
	if dp.DialTimeout > 0 {
		return dp.DialTimeout
	}
	return 10 * time.Second
}

// dialContext 封装 	net.Dial()
func (dp *TcpReverseProxy) dialContext() func(ctx context.Context, network, address string) (net.Conn, error) {
	if dp.DialContext != nil {
		return dp.DialContext
	}
	var d = net.Dialer{
		Timeout:   dp.DialTimeout,
		KeepAlive: dp.KeepAlivePeriod,
	}
	// net.Dial
	return d.DialContext
}

func (dp *TcpReverseProxy) keepAlivePeriod() time.Duration {
	if dp.KeepAlivePeriod != 0 {
		return dp.KeepAlivePeriod
	}
	return time.Minute
}

func (dp *TcpReverseProxy) ServeTCP(ctx context.Context, src net.Conn) {
	var cancel context.CancelFunc
	if dp.DialTimeout >= 0 {
		ctx, cancel = context.WithTimeout(ctx, dp.dialTimeout())
	}

	dp.director(dp)

	// src: client -> [ reverse proxy ]
	// dst: [ reverse proxy ] -> internal server
	dst, err := dp.dialContext()(ctx, "tcp", dp.Addr)
	if cancel != nil {
		cancel()
	}
	if err != nil {
		dp.onDialError()(src, err)
		return
	}

	defer func() { go dst.Close() }()

	//设置dst的 keepAlive 参数,在数据请求之前
	if ka := dp.keepAlivePeriod(); ka > 0 {
		if c, ok := dst.(*net.TCPConn); ok {
			c.SetKeepAlive(true)
			c.SetKeepAlivePeriod(ka)
		}
	}
	errChan := make(chan error, 1)
	// 数据传输
	// dest -> src
	go dp.proxyCopy(errChan, src, dst)
	// src  -> dest
	go dp.proxyCopy(errChan, dst, src)
	<-errChan
}

func (dp *TcpReverseProxy) onDialError() func(src net.Conn, dstDialErr error) {
	if dp.OnDialError != nil {
		return dp.OnDialError
	}
	return func(src net.Conn, dstDialErr error) {
		logx.Error("tcpproxy: for incoming conn %v, error dialing %q: %v", src.RemoteAddr().String(), dp.Addr, dstDialErr)
		src.Close()
	}
}

func (dp *TcpReverseProxy) proxyCopy(errChan chan<- error, dst, src net.Conn) {
	_, err := io.Copy(dst, src)
	errChan <- err
}
