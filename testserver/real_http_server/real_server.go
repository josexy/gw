package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/pkg/constants"
	"github.com/josexy/gw/pkg/etcd"
)

type RealServer struct {
	Addr       string
	ServiceReg *etcd.ServiceRegister
}

func (server *RealServer) Run() {
	logx.Debug("start server at: %v", server.Addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.HelloHandler)
	mux.HandleFunc("/error", server.ErrorHandler)
	mux.HandleFunc("/test_strip_uri/timeout", server.TimeoutHandler)

	svr := &http.Server{
		Addr:         server.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		err := svr.ListenAndServe()
		// err := svr.ListenAndServeTLS(certs.Path("public.crt"), certs.Path("server.pem"))
		if err != nil && err != http.ErrServerClosed {
			logx.Fatal("%v", err)
			return
		}
	}()

	// 服务注册
	server.registerService(server.Addr)
}

func (server *RealServer) unregisterService() {
	logx.Debug("revoke service and close server: %s", server.Addr)
	if server.ServiceReg != nil {
		_ = server.ServiceReg.Close()
	}
}

func (server *RealServer) registerService(addr string) {

	endpoints := []string{"127.0.0.1:2379"}
	prefix := constants.EtcdPrefix

	go func() {
		service, err := etcd.NewServiceRegister(endpoints, 5)
		if err != nil {
			logx.Error("connect server register: %v", err)
			return
		}
		server.ServiceReg = service

		key := prefix + addr
		logx.Debug("register service to etcd: %s", key)
		_ = service.RegisterService(key, addr)
	}()
}

func (server *RealServer) HelloHandler(w http.ResponseWriter, r *http.Request) {
	upath := fmt.Sprintf("full path: http://%s%s\n", server.Addr, r.URL.Path)
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v\n",
		r.RemoteAddr, r.Header.Get("X-Forwarded-For"), r.Header.Get("X-Real-Ip"))

	logx.Debug(upath)
	logx.Debug("%v", r.Header)
	header := fmt.Sprintf("headers =%v\n", r.Header)
	_, _ = io.WriteString(w, upath)
	_, _ = io.WriteString(w, realIP)
	_, _ = io.WriteString(w, header)
}

func (server *RealServer) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = io.WriteString(w, "error handler")
}

func (server *RealServer) TimeoutHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 5)
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "timeout handler")
}

func main() {
	var port1, port2 int
	var err error

	if len(os.Args) == 1 {
		port1 = 2003
		port2 = 2004
	} else {
		port1, err = strconv.Atoi(os.Args[1])
		if err != nil {
			logx.Fatal("%v", err)
			return
		}
		port2, err = strconv.Atoi(os.Args[2])
		if err != nil {
			logx.Fatal("%v", err)
			return
		}
	}

	rs1 := &RealServer{Addr: fmt.Sprintf("127.0.0.1:%d", port1)}
	rs2 := &RealServer{Addr: fmt.Sprintf("127.0.0.1:%d", port2)}

	rs1.Run()
	rs2.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT)
	<-interrupt

	rs1.unregisterService()
	rs2.unregisterService()
}
