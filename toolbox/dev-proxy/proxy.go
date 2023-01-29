package main

import (
	"bytes"
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

	// replace
	proxyLog := ""
	findNewUrl := findReplaceByRegexp(*proxyReq)
	if findNewUrl != nil {
		proxyLog = rewriteRequest(findNewUrl, proxyReq)
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

func rewriteRequest(newUrl *url.URL, proxyReq *http.Request) string {
	var proxyLog string
	var id string
	if newUrl.Path == proxyReq.URL.Path {
		proxyLog = fmt.Sprintf("%s => %s", proxyReq.Host+proxyReq.URL.Path, newUrl.Host+" .")
	} else {
		proxyLog = fmt.Sprintf("%s => %s", proxyReq.Host+proxyReq.URL.Path, newUrl.Host+newUrl.Path)
	}

	// 记录请求
	id = uuid.New().String()
	bodyBt, _ := io.ReadAll(proxyReq.Body)
	//fmt.Println(string(bodyBt), err)
	saveRequest(ReqLog{Id: id, Url: proxyReq.Host + proxyReq.URL.Path, Header: proxyReq.Header, Body: string(bodyBt), Time: time.Now()})
	// 回写流
	proxyReq.Body = io.NopCloser(bytes.NewBuffer(bodyBt))

	proxyReq.Host = newUrl.Host
	//proxyReq.URL.Scheme = newUrl.Scheme
	proxyReq.URL.Host = newUrl.Host
	proxyReq.URL.Path = newUrl.Path
	//proxyReq.URL.RawQuery = newUrl.RawQuery
	return proxyLog
}
