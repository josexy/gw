package httpproxyrouter

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/certs"
	"github.com/josexy/gw/global"
	"github.com/josexy/gw/logx"
	httpproxymiddleware "github.com/josexy/gw/proxy_middleware/http_proxy_middleware"
)

var (
	httpServer  *http.Server
	httpsServer *http.Server
)

func newHttpServerWithRouter() *http.Server {

	gin.DefaultWriter = ioutil.Discard
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(
		httpproxymiddleware.HTTPAccessModeMiddleware(),

		httpproxymiddleware.HTTPFlowCountMiddleware(),
		httpproxymiddleware.HTTPFlowLimitMiddleware(),

		httpproxymiddleware.HTTPWhiteListMiddleware(),
		httpproxymiddleware.HTTPBlackListMiddleware(),

		httpproxymiddleware.HTTPHeaderTransferMiddleware(),
		httpproxymiddleware.HTTPStripUriMiddleware(),
		httpproxymiddleware.HTTPURLRewriteMiddleware(),

		httpproxymiddleware.HTTPReverseProxyMiddleware(),
	)

	addr := fmt.Sprintf("%s:%d", global.ProxyConfig.Common.Addr, global.ProxyConfig.Http.Port)
	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    time.Second * time.Duration(global.ProxyConfig.Http.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(global.ProxyConfig.Http.WriteTimeout),
		MaxHeaderBytes: global.ProxyConfig.Http.MaxHeaderBytes,
	}
	return server
}

func HttpServerRun() {

	httpServer = newHttpServerWithRouter()

	logx.Info("run http proxy server successfully!")
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logx.Fatal("http proxy server err: %v", err)
	}
}

func HttpsServerRun() {

	httpsServer = newHttpServerWithRouter()
	addr := fmt.Sprintf("%s:%d", global.ProxyConfig.Common.Addr, global.ProxyConfig.Https.Port)
	httpsServer.Addr = addr
	httpsServer.ReadTimeout = time.Second * time.Duration(global.ProxyConfig.Https.ReadTimeout)
	httpsServer.WriteTimeout = time.Second * time.Duration(global.ProxyConfig.Https.WriteTimeout)
	httpsServer.MaxHeaderBytes = global.ProxyConfig.Https.MaxHeaderBytes

	logx.Info("run https proxy server successfully!")
	if err := httpsServer.ListenAndServeTLS(certs.Path("public.crt"), certs.Path("server.pem")); err != nil && err != http.ErrServerClosed {
		logx.Fatal("https proxy server err: %v", err)
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logx.Error("stop http proxy server err: %v", err)
	}
	logx.Warn("stop http proxy server %v", httpServer.Addr)
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := httpsServer.Shutdown(ctx); err != nil {
		logx.Error("stop https proxy server err: %v", err)
	}
	logx.Warn("stop https proxy server %v", httpsServer.Addr)
}
