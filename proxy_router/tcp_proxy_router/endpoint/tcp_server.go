package endpoint

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrServerClosed = errors.New("tcp: Server closed")
	ErrAbortHandler = errors.New("tcp: abort TCPHandler")

	ServerContextKey    = &contextKey{name: "tcp-server"}
	LocalAddrContextKey = &contextKey{name: "local-addr"}
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "tcp_proxy context value " + k.name
}

// onceCloseListener 参考 http.onceCloseListener
type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() {
	oc.closeErr = oc.Listener.Close()
}

type TcpHandler interface {
	ServeTCP(ctx context.Context, conn net.Conn)
}

type TcpServer struct {
	Addr    string
	Handler TcpHandler

	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	KeepAliveTimeout time.Duration

	BaseContext context.Context

	listener   *onceCloseListener
	mu         sync.Mutex
	err        error
	inShutdown int32
	doneChan   chan struct{}
}

func (srv *TcpServer) Close() error {
	atomic.StoreInt32(&srv.inShutdown, 1)
	close(srv.doneChan)
	return srv.listener.Close()
}

func (srv *TcpServer) shuttingDown() bool {
	return atomic.LoadInt32(&srv.inShutdown) != 0
}

func (srv *TcpServer) ListenAndServe() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}

	addr := srv.Addr
	if addr == "" {
		panic("tcp server need address")
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	tcpListener := ln.(*net.TCPListener)
	srv.listener = &onceCloseListener{Listener: tcpListener}
	return srv.Serve(ln)
}

func (srv *TcpServer) Serve(listener net.Listener) error {
	defer srv.listener.Close()

	if srv.BaseContext == nil {
		srv.BaseContext = context.Background()
	}
	ctx := context.WithValue(srv.BaseContext, ServerContextKey, srv)
	for {
		rw, err := listener.Accept()
		if err != nil {
			select {
			// 服务器退出
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
			continue
		}
		conn := srv.newConn(rw)
		go conn.serve(ctx)
	}
}

func (srv *TcpServer) newConn(rwc net.Conn) *TcpConn {
	conn := &TcpConn{
		rwc:    rwc,
		server: srv,
	}
	if d := conn.server.ReadTimeout; d != 0 {
		// 读超时
		conn.rwc.SetReadDeadline(time.Now().Add(d))
	}
	if d := conn.server.WriteTimeout; d != 0 {
		// 写超时
		conn.rwc.SetWriteDeadline(time.Now().Add(d))
	}
	if d := conn.server.KeepAliveTimeout; d != 0 {
		// keepalive心跳包
		if tcpConn, ok := conn.rwc.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(d)
		}
	}
	return conn
}

func (srv *TcpServer) getDoneChan() <-chan struct{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	return srv.doneChan
}

func ListenAndServe(addr string, handler TcpHandler) error {
	server := &TcpServer{
		Addr:     addr,
		Handler:  handler,
		doneChan: make(chan struct{}),
	}
	return server.ListenAndServe()
}
