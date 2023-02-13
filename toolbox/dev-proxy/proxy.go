package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/kuangcp/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	var reqLog *ReqLog[Message]
	findNewUrl, proxyType := findReplaceByRegexp(*proxyReq)
	if findNewUrl != nil {
		proxyLog, reqLog = rewriteRequestAndBuildLog(findNewUrl, proxyReq)
		defer saveReqLog(reqLog)
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
		if strings.Contains(err.Error(), "connect: connection refused") {
			logger.Error("%v proxy error %v", r.URL.String(), "down")
			reqLog.Status = fmt.Sprint(http.StatusServiceUnavailable, " server refused")
			reqLog.StatusCode = 98
			reqLog.ResTime = time.Now()
			reqLog.ElapsedTime = fmtDuration(reqLog.ResTime.Sub(reqLog.ReqTime))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		if proxyLog == "" {
			logger.Error("%4vms %v proxy error %v", endMs-startMs, r.URL.String(), err)
		} else {
			logger.Error("%4vms %v proxy error %v", endMs-startMs, proxyLog, err)
		}
		if reqLog != nil {
			reqLog.Status = fmt.Sprint(http.StatusInternalServerError, " server error")
			reqLog.StatusCode = 99
			reqLog.ResTime = time.Now()
			reqLog.ElapsedTime = fmtDuration(reqLog.ResTime.Sub(reqLog.ReqTime))
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if proxyLog != "" {
		if proxyType == Proxy {
			proxyLog += " SELF"
		}
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
		resMes := Message{Header: res.Header, Body: bytes}
		reqLog.Response = resMes
		reqLog.ResTime = time.Now()
		reqLog.ElapsedTime = fmtDuration(reqLog.ResTime.Sub(reqLog.ReqTime))
		reqLog.Status = res.Status
		reqLog.StatusCode = res.StatusCode
	}

	w.WriteHeader(res.StatusCode)
	if res.Body != nil {
		written, err := io.Copy(w, res.Body)
		if err != nil {
			logger.Error("%3vms %v %v", endMs-startMs, written, err)
		}
	}
}

func fmtDuration(d time.Duration) string {
	ms := d.Milliseconds()
	d = d.Round(time.Millisecond)
	if ms < 10_000 {
		return fmt.Sprintf("%vms", ms)
	}
	return d.String()
}

func rewriteRequestAndBuildLog(newUrl *url.URL, proxyReq *http.Request) (string, *ReqLog[Message]) {
	now := time.Now()
	id := uuid.New().String()

	bodyBt, body := copyStream(proxyReq.Body)
	query, _ := url.QueryUnescape(proxyReq.URL.String())
	reqMes := Message{Header: proxyReq.Header, Body: filterFileType(bodyBt)}
	id = fmt.Sprintf("%v%v", id[0:8], now.UnixMilli()%1000)
	cacheId := fmt.Sprintf("%v  %v", now.Format("01-02 15:04:05.000"), id)
	reqLog := &ReqLog[Message]{Id: id, CacheId: cacheId, Url: query, Request: reqMes, ReqTime: now, Method: proxyReq.Method}
	connection.ZAdd(RequestList, redis.Z{Member: cacheId, Score: float64(reqLog.ReqTime.UnixNano())})

	var logStr string
	if newUrl.Path == proxyReq.URL.Path {
		logStr = fmt.Sprintf("%v %s => %s", id, proxyReq.Host+proxyReq.URL.Path, newUrl.Host+" .")
	} else {
		logStr = fmt.Sprintf("%v %s => %s", id, proxyReq.Host+proxyReq.URL.Path, newUrl.Host+newUrl.Path)
	}

	proxyReq.Body = body
	proxyReq.Host = newUrl.Host
	//proxyReq.URL.Scheme = newUrl.Scheme
	proxyReq.URL.Host = newUrl.Host
	proxyReq.URL.Path = newUrl.Path
	//proxyReq.URL.RawQuery = newUrl.RawQuery
	return logStr, reqLog
}
func filterFileType(body []byte) []byte {
	// TODO 比较字节数组
	str := string(body)
	if strings.HasPrefix(str, "------") {
		endIdx := strings.Index(str, "Content-Type:")
		return []byte(str[:endIdx])
	}
	return body
}
