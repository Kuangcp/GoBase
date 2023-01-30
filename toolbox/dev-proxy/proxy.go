package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kuangcp/logger"
	"io"
	"net/http"
	"net/url"
	"time"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		handleHttps(w, r)
		return
	}

	proxyReq := new(http.Request)
	*proxyReq = *r

	// replace, if not use proxy, log will be nil
	proxyLog := ""
	var reqLog *ReqLog
	findNewUrl := findReplaceByRegexp(*proxyReq)
	if findNewUrl != nil {
		proxyLog, reqLog = rewriteRequest(findNewUrl, proxyReq)
		defer saveRequest(reqLog)
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
		logger.Error("%4vms %v proxy error %v", endMs-startMs, proxyLog, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if proxyLog != "" {
		logger.Debug("%4vms %v", endMs-startMs, proxyLog)
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

	if reqLog != nil {
		bytes, body := copyStream(res.Body)
		res.Body = body
		resMes := Message{Header: res.Header, Body: string(bytes)}
		reqLog.Response = resMes
	}

	w.WriteHeader(res.StatusCode)
	if res.Body != nil {
		written, err := io.Copy(w, res.Body)
		if err != nil {
			logger.Error("%3vms %v %v", endMs-startMs, written, err)
		}
	}
}

func rewriteRequest(newUrl *url.URL, proxyReq *http.Request) (string, *ReqLog) {
	var proxyLog string
	var id string
	if newUrl.Path == proxyReq.URL.Path {
		proxyLog = fmt.Sprintf("%s => %s", proxyReq.Host+proxyReq.URL.Path, newUrl.Host+" same path")
	} else {
		proxyLog = fmt.Sprintf("%s => %s", proxyReq.Host+proxyReq.URL.Path, newUrl.Host+newUrl.Path)
	}

	// 记录请求
	id = uuid.New().String()
	bodyBt, body := copyStream(proxyReq.Body)

	query, _ := url.QueryUnescape(proxyReq.URL.String())
	reqMes := Message{Header: proxyReq.Header, Body: string(bodyBt)}
	log := &ReqLog{Id: id, Url: query, Request: reqMes, Time: time.Now()}

	proxyReq.Body = body
	proxyReq.Host = newUrl.Host
	//proxyReq.URL.Scheme = newUrl.Scheme
	proxyReq.URL.Host = newUrl.Host
	proxyReq.URL.Path = newUrl.Path
	//proxyReq.URL.RawQuery = newUrl.RawQuery
	return proxyLog, log
}
