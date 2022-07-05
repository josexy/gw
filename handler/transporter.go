package handler

import (
	"net"
	"net/http"
	"time"

	"github.com/josexy/gw/model"
)

var TransporterHandler *Transporter

type Transporter struct {
	handler MapBaseHandler[*http.Transport]
}

func init() {
	TransporterHandler = NewTransporter()
}

func NewTransporter() *Transporter {
	return &Transporter{
		handler: NewMapBaseHandler[*http.Transport](),
	}
}

func (t *Transporter) GetTransporter(serviceDetail *model.ServiceDetail) (*http.Transport, error) {
	t.handler.RLock()
	if transport, ok := t.handler.Cache[serviceDetail.Info.ServiceName]; ok {
		t.handler.RUnlock()
		return transport, nil
	}
	t.handler.RUnlock()

	// 连接上游服务器设置
	// 连接超时时间
	if serviceDetail.LoadBalance.UpstreamConnectTimeout == 0 {
		serviceDetail.LoadBalance.UpstreamConnectTimeout = 30
	}
	// 最大空闲连接数量
	if serviceDetail.LoadBalance.UpstreamMaxIdle == 0 {
		serviceDetail.LoadBalance.UpstreamMaxIdle = 100
	}
	// 空闲超时时间，超过则自动关闭
	if serviceDetail.LoadBalance.UpstreamIdleTimeout == 0 {
		serviceDetail.LoadBalance.UpstreamIdleTimeout = 90
	}
	// 响应header超时时间
	if serviceDetail.LoadBalance.UpstreamHeaderTimeout == 0 {
		serviceDetail.LoadBalance.UpstreamHeaderTimeout = 30
	}

	dialer := &net.Dialer{
		Timeout:   time.Duration(serviceDetail.LoadBalance.UpstreamConnectTimeout) * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          serviceDetail.LoadBalance.UpstreamMaxIdle,
		IdleConnTimeout:       time.Duration(serviceDetail.LoadBalance.UpstreamIdleTimeout) * time.Second,
		ResponseHeaderTimeout: time.Duration(serviceDetail.LoadBalance.UpstreamHeaderTimeout) * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
	}

	t.handler.Lock()
	defer t.handler.Unlock()

	t.handler.List = append(t.handler.List, transport)
	t.handler.Cache[serviceDetail.Info.ServiceName] = transport

	return transport, nil
}
