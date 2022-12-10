package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/kuangcp/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	port       int
	reloadConf bool
)

func concatIgnoreSlash(left, right string) string {
	aslash := strings.HasSuffix(left, "/")
	bslash := strings.HasPrefix(right, "/")
	switch {
	case aslash && bslash:
		return left + right[1:]
	case !aslash && !bslash:
		return left + "/" + right
	}
	return left + right
}

func handlePath(origin, target *url.URL, path string) string {
	// 原始路径去前缀
	if origin.Path != "" && strings.HasPrefix(path, origin.Path) {
		path = path[len(origin.Path):]
	}
	// 原始路径加前缀
	if target.Path != "" {
		path = concatIgnoreSlash(target.Path, path)
	}
	return path
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		handleHttps(w, r)
		return
	}

	proxyReq := new(http.Request)
	*proxyReq = *r

	log := proxyReplace(proxyReq)

	startMs := time.Now().UnixMilli()
	if q := proxyReq.URL.RawQuery; q != "" {
		proxyReq.URL.RawPath = proxyReq.URL.Path + "?" + q
	} else {
		proxyReq.URL.RawPath = proxyReq.URL.Path
	}
	proxyReq.Proto = "HTTP/1.1"
	proxyReq.ProtoMajor = 1
	proxyReq.ProtoMinor = 1
	proxyReq.Close = false

	//if websocketHandler(w, r, proxyReq) {
	//	return
	//}

	transport := http.DefaultTransport
	res, err := transport.RoundTrip(proxyReq)
	endMs := time.Now().UnixMilli()
	if err != nil {
		logger.Error("%vms %v proxy error %v", endMs-startMs, log, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if log != "" {
		logger.Debug("%vms %v", endMs-startMs, log)
	}

	hdr := w.Header()
	for k, vv := range res.Header {
		for _, v := range vv {
			hdr.Add(k, v)
		}
	}
	for _, c := range res.Cookies() {
		w.Header().Add("Set-Cookie", c.Raw)
	}

	w.WriteHeader(res.StatusCode)
	if res.Body != nil {
		written, err := io.Copy(w, res.Body)
		if err != nil {
			logger.Error(written, err)
		}
	}
}

func proxyReplace(proxyReq *http.Request) string {
	originUrl, targetUrl := findTargetReplace(proxyReq)
	if targetUrl != nil {
		path := handlePath(originUrl, targetUrl, proxyReq.URL.Path)
		log := ""
		if path == proxyReq.URL.Path {
			log = fmt.Sprintf("%s => %s", proxyReq.Host+proxyReq.URL.Path, targetUrl.Host+" .")
		} else {
			log = fmt.Sprintf("%s => %s", proxyReq.Host+proxyReq.URL.Path, targetUrl.Host+path)
		}

		proxyReq.Host = targetUrl.Host
		//o.URL.Scheme = targetURL.Scheme
		proxyReq.URL.Host = targetUrl.Host
		proxyReq.URL.Path = path
		//o.URL.RawQuery = targetUrl.RawQuery
		return log
	} else {
		//logger.Debug("direct:", proxyReq.Host)
		return ""
	}
}

func main() {
	flag.IntVar(&port, "p", 1234, "port")
	flag.BoolVar(&reloadConf, "r", false, "auto reload changed config")
	flag.Parse()

	initConfig()

	logger.Info("Start serving on 127.0.0.1:%d", port)

	cert, err := genCertificate()
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Addr:      fmt.Sprintf("0.0.0.0:%d", port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxyHandler(w, r)
		}),
	}

	logger.Fatal(server.ListenAndServe())

	//err := http.ListenAndServe(fmt.Sprintf(":%d", port),
	//	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		proxyHandler(w, r)
	//	}))
	//if err != nil {
	//	logger.Error(err)
	//}
	//os.Exit(0)
}
