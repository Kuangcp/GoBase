package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/kuangcp/logger"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	port       int
	reloadConf bool
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		handleHttps(w, r)
		return
	}

	proxyReq := new(http.Request)
	*proxyReq = *r

	// replace
	newUrl := findReplaceByRegexp(*proxyReq)
	proxyLog := ""
	if newUrl != nil {
		if newUrl.Path == proxyReq.URL.Path {
			proxyLog = fmt.Sprintf("%s => %s", proxyReq.Host+proxyReq.URL.Path, newUrl.Host+" .")
		} else {
			proxyLog = fmt.Sprintf("%s => %s", proxyReq.Host+proxyReq.URL.Path, newUrl.Host+newUrl.Path)
		}

		proxyReq.Host = newUrl.Host
		//proxyReq.URL.Scheme = newUrl.Scheme
		proxyReq.URL.Host = newUrl.Host
		proxyReq.URL.Path = newUrl.Path
		//proxyReq.URL.RawQuery = newUrl.RawQuery
	}

	// rebuild
	if q := proxyReq.URL.RawQuery; q != "" {
		proxyReq.URL.RawPath = proxyReq.URL.Path + "?" + q
	} else {
		proxyReq.URL.RawPath = proxyReq.URL.Path
	}
	proxyReq.Proto = "HTTP/1.1"
	proxyReq.ProtoMajor = 1
	proxyReq.ProtoMinor = 1
	proxyReq.Close = false

	// TODO websocket
	//if websocketHandler(w, r, proxyReq) {
	//	return
	//}

	transport := http.DefaultTransport
	startMs := time.Now().UnixMilli()
	res, err := transport.RoundTrip(proxyReq)
	endMs := time.Now().UnixMilli()
	if err != nil {
		logger.Error("%4v %v proxy error %v", endMs-startMs, proxyLog, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if proxyLog != "" {
		logger.Debug("%4v %v", endMs-startMs, proxyLog)
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
			logger.Error("%3vms %v %v", endMs-startMs, written, err)
		}
	}
}

func findReplaceByRegexp(proxyReq http.Request) *url.URL {
	lock.RLock()
	defer lock.RUnlock()

	for k, v := range proxyValMap {
		fullUrl := proxyReq.URL.Scheme + "://" + proxyReq.URL.Host + proxyReq.URL.Path
		tryResult := tryToReplacePath(k, v, fullUrl)
		if tryResult == "" {
			continue
		}

		parse, err := url.Parse(tryResult)
		if err != nil {
			logger.Error(err)
		}

		return parse
	}

	return nil
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
