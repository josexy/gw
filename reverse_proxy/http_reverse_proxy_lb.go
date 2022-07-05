package reverseproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/reverse_proxy/load_balance"
)

func NewHTTPLoadBalanceReverseProxy(lb load_balance.LoadBalance, transport *http.Transport) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		nextAddr, err := lb.Get(req.RemoteAddr)

		if err != nil || nextAddr == "" {
			logx.Error("the upstream address is invalid or not found")
		}

		target, err := url.Parse(nextAddr)
		if err != nil {
			logx.Error("url parse err: %v", err)
			return
		}

		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme // http/https
		req.URL.Host = target.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
	}
	modifyFunc := func(resp *http.Response) error {
		return nil
	}

	errorFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		w.Write([]byte(fmt.Sprintf("connect upstream address failed: %v", err)))
	}

	return &httputil.ReverseProxy{
		Transport:      transport,
		Director:       director,
		ModifyResponse: modifyFunc,
		ErrorHandler:   errorFunc,
	}
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
