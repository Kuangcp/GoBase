package core

import (
	"bytes"
	"compress/gzip"
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

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		handleHttps(w, r)
		return
	}

	w.Header().Add("Server", "dev-proxy")

	defer func() {
		re := recover()
		if re != nil {
			logger.Error("代理异常: ", re)
			w.WriteHeader(911)
		}
	}()

	proxyReq := new(http.Request)
	*proxyReq = *r

	// replace, if not use proxy, log will be nil
	proxyLog := ""
	var reqLog *ReqLog[Message]
	findNewUrl, proxyType := findReplaceByRegexp(*proxyReq)
	ignoreStorage := matchIgnoreStorage(*proxyReq)
	if findNewUrl != nil {
		proxyLog, reqLog = rewriteRequestAndBuildLog(findNewUrl, proxyReq, ignoreStorage)
		if !ignoreStorage {
			defer saveReqLog(reqLog)
		}
	}

	// TODO websocket
	//if websocketHandler(w, r, proxyReq) {
	//	return
	//}

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

	transport := http.DefaultTransport.(*http.Transport)

	startMs := time.Now().UnixMilli()
	res, err := transport.RoundTrip(proxyReq)
	endMs := time.Now().UnixMilli()
	waste := endMs - startMs
	if err != nil {
		handleError(w, r, err, reqLog, proxyLog, waste)
		return
	}

	if proxyLog != "" {
		if proxyType == Proxy {
			proxyLog += " SELF"
		}
		logger.Debug("%4vms %v", waste, proxyLog)
	}

	copyResponseHeader(w, res)

	if !ignoreStorage {
		fillReqLogResponse(reqLog, res)
	}

	if res.Body != nil {
		written, err := io.Copy(w, res.Body)
		if err != nil {
			logger.Error("%3vms %v %v", waste, written, err)
		}
	}
}

func copyResponseHeader(w http.ResponseWriter, res *http.Response) {
	header := w.Header()
	for k, vv := range res.Header {
		for _, v := range vv {
			header.Add(k, v)
		}
	}
	for _, c := range res.Cookies() {
		header.Add("Set-Cookie", c.Raw)
	}
	w.WriteHeader(res.StatusCode)
}

func handleCompressed(msg *Message, res *http.Response) {
	encoding := res.Header.Get("Content-Encoding")
	if encoding == "" {
		return
	}
	if encoding == "gzip" {
		reader, err := gzip.NewReader(bytes.NewBuffer(msg.Body))
		if err != nil {
			logger.Error(err)
			return
		}
		defer reader.Close()
		var buff bytes.Buffer

		_, err = io.Copy(&buff, reader)
		if err != nil {
			logger.Error(err)
			return
		}
		msg.Body = buff.Bytes()
	}
}

func fillReqLogResponse(reqLog *ReqLog[Message], res *http.Response) {
	if reqLog == nil {
		return
	}
	bodyBts, body := copyStream(res.Body)
	res.Body = body
	resMes := Message{Header: res.Header, Body: bodyBts}
	handleCompressed(&resMes, res)

	reqLog.Response = resMes
	reqLog.ResTime = time.Now()
	reqLog.ElapsedTime = fmtDuration(reqLog.ResTime.Sub(reqLog.ReqTime))
	reqLog.Status = res.Status
	reqLog.StatusCode = res.StatusCode
}

func handleError(w http.ResponseWriter, r *http.Request, err error, reqLog *ReqLog[Message], proxyLog string, waste int64) {
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
		logger.Error("%4vms %v proxy error %v", waste, r.URL.String(), err)
	} else {
		logger.Error("%4vms %v proxy error %v", waste, proxyLog, err)
	}
	if reqLog != nil {
		reqLog.Status = fmt.Sprint(http.StatusInternalServerError, " server error")
		reqLog.StatusCode = 99
		reqLog.ResTime = time.Now()
		reqLog.ElapsedTime = fmtDuration(reqLog.ResTime.Sub(reqLog.ReqTime))
	}
	w.WriteHeader(http.StatusInternalServerError)
}

func fmtDuration(d time.Duration) string {
	ms := d.Milliseconds()
	d = d.Round(time.Millisecond)
	if ms < 10_000 {
		return fmt.Sprintf("%vms", ms)
	}
	return d.String()
}

func rewriteRequestAndBuildLog(newUrl *url.URL, proxyReq *http.Request, ignoreStorage bool) (string, *ReqLog[Message]) {
	now := time.Now()
	id := uuid.New().String()

	bodyBt, body := copyStream(proxyReq.Body)
	query, _ := url.QueryUnescape(proxyReq.URL.String())
	reqMes := Message{Header: proxyReq.Header, Body: filterFormType(bodyBt)}

	id = fmt.Sprintf("%v%v", id[0:8], now.UnixMilli()%1000)
	cacheId := fmt.Sprintf("%v  %v", now.Format("01-02 15:04:05.000"), id)
	reqLog := &ReqLog[Message]{Id: id, CacheId: cacheId, Url: query, Request: reqMes, ReqTime: now, Method: proxyReq.Method}

	if !ignoreStorage {
		// redis cache
		connection.ZAdd(RequestList, redis.Z{Member: cacheId, Score: float64(reqLog.ReqTime.UnixNano())})
		connection.HSet(RequestUrlList, id, proxyReq.URL.String())
	}

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

// request body : start with "------"
func filterFormType(s []byte) []byte {
	if len(s) > 7 && s[0] == 45 && s[1] == 45 && s[2] == 45 &&
		s[3] == 45 && s[4] == 45 && s[5] == 45 && s[6] == 45 {
		var r []byte
		for _, i := range s {
			if i == 10 {
				return r
			}
			r = append(r, i)
		}
		return r
	}
	return s
}
