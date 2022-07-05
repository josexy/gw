package router

import (
	"context"
	"math"
	"net"
	"sync"

	"github.com/josexy/gw/proxy_router/tcp_proxy_router/endpoint"
)

const abortIndex int8 = math.MaxInt8 / 2 //最多 63 个中间件

type TcpHandlerFunc func(*TcpRouterContext)

type TcpRouterContext struct {
	Conn   net.Conn
	Ctx    context.Context
	index  int8
	router *TcpRouter
}

type TcpRouter struct {
	handlers []TcpHandlerFunc
	coreFunc func(*TcpRouterContext) endpoint.TcpHandler
	pool     sync.Pool
}

func NewTcpRouter() *TcpRouter {
	router := &TcpRouter{}
	router.pool.New = func() interface{} {
		return &TcpRouterContext{index: -1, router: router}
	}
	return router
}

func (r *TcpRouter) Serve(fn func(*TcpRouterContext) endpoint.TcpHandler) {
	r.coreFunc = fn
}

func (r *TcpRouter) ServeTCP(ctx context.Context, conn net.Conn) {
	c := r.pool.Get().(*TcpRouterContext)
	c.Conn = conn
	c.Ctx = ctx

	c.Reset()
	c.Next()
	if r.coreFunc != nil {
		r.coreFunc(c).ServeTCP(ctx, conn)
	}
	r.pool.Put(c)
}

func (r *TcpRouter) Use(middlewares ...TcpHandlerFunc) *TcpRouter {
	r.handlers = append(r.handlers, middlewares...)
	return r
}

func (c *TcpRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

func (c *TcpRouterContext) Set(key, val interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, val)
}

func (c *TcpRouterContext) Next() {
	c.index++
	for c.index < int8(len(c.router.handlers)) {
		c.router.handlers[c.index](c)
		c.index++
	}
}

func (c *TcpRouterContext) Abort() {
	c.index = abortIndex
}

func (c *TcpRouterContext) IsAborted() bool {
	return c.index >= abortIndex
}

func (c *TcpRouterContext) Reset() {
	c.index = -1
}
