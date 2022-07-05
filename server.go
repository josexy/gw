package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/logx"
)

type Server struct {
	address string
	app     *gin.Engine
}

func NewServer(port int, app *gin.Engine) *Server {
	return &Server{
		address: fmt.Sprintf(":%d", port),
		app:     app,
	}
}

func (svr *Server) Run() {
	server := &http.Server{
		Addr:           svr.address,
		Handler:        svr.app,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 9102,
	}
	go func() {
		logx.Debug("start server at: %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
