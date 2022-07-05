package endpoint

import (
	"context"
	"net"

	"github.com/josexy/gw/logx"
)

// TcpConn 参考 http.conn
type TcpConn struct {
	server     *TcpServer
	rwc        net.Conn
	remoteAddr string
}

func (conn *TcpConn) close() error {
	return conn.rwc.Close()
}

func (conn *TcpConn) serve(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			logx.Debug("tcp: panic serving %v: %v", conn.remoteAddr, err)
		}
		conn.close()
	}()

	logx.Debug("remote client ip address: %v", conn.rwc.RemoteAddr().String())
	logx.Debug("local server ip address: %v", conn.rwc.LocalAddr().String())

	conn.remoteAddr = conn.rwc.RemoteAddr().String()
	ctx = context.WithValue(ctx, LocalAddrContextKey, conn.rwc.LocalAddr())

	if conn.server.Handler == nil {
		panic("server handler empty")
	}
	conn.server.Handler.ServeTCP(ctx, conn.rwc)
}
