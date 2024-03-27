package core

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		handleHttps(w, r)
		return
	}

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
	findNewUrl, proxyType := FindReplaceByRegexp(*proxyReq)
	needStorage := MatchNeedStorage(*proxyReq, proxyType)
	if findNewUrl != nil {
		proxyLog, reqLog = RewriteRequestAndBuildLog(findNewUrl, proxyReq, needStorage)
	}

	w.Header().Add("Ack-Proxy", proxyType+"  "+r.Host)

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
		HandleRespError(w, r, err, reqLog, proxyLog, waste)
		return
	}

	if proxyLog != "" {
		if proxyType == Proxy {
			proxyLog += " SELF"
		}
		logger.Debug("%4vms %v", waste, proxyLog)
	}

	CopyResponseHeader(w, res)

	if needStorage {
		TrySaveLog(reqLog, res)
	}

	if res.Body != nil {
		written, err := io.Copy(w, res.Body)
		if err != nil {
			logger.Error("%3vms %v %v", waste, written, err)
		}
	}
}

func CopyResponseHeader(w http.ResponseWriter, res *http.Response) {
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

func HandleCompressed(msg *Message, res *http.Response) {
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

func FillReqLogResponse(reqLog *ReqLog[Message], res *http.Response) {
	if reqLog == nil {
		return
	}
	bodyBts, body := ctool.CopyStream(res.Body)
	res.Body = body
	resMes := Message{Header: res.Header, Body: bodyBts}
	HandleCompressed(&resMes, res)

	reqLog.Response = resMes
	reqLog.ResTime = time.Now()
	reqLog.ElapsedTime = FmtDuration(reqLog.ResTime.Sub(reqLog.ReqTime))
	reqLog.Status = res.Status
	reqLog.StatusCode = res.StatusCode
}

func HandleRespError(w http.ResponseWriter, r *http.Request, err error, reqLog *ReqLog[Message], proxyLog string, waste int64) {
	if proxyLog == "" {
		logger.Error("%4vms %v proxy error %v", waste, r.URL.String(), err)
	} else {
		logger.Error("%4vms %v proxy error %v", waste, proxyLog, err)
	}

	if reqLog != nil {
		if strings.Contains(err.Error(), "connect: connection refused") {
			logger.Error("%v proxy error %v", r.URL.String(), "down")
			reqLog.Status = fmt.Sprint(http.StatusServiceUnavailable, " server refused")
			reqLog.StatusCode = 98
			reqLog.ResTime = time.Now()
			reqLog.ElapsedTime = FmtDuration(reqLog.ResTime.Sub(reqLog.ReqTime))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		reqLog.Status = fmt.Sprint(http.StatusInternalServerError, " server error")
		reqLog.StatusCode = 99
		reqLog.ResTime = time.Now()
		reqLog.ElapsedTime = FmtDuration(reqLog.ResTime.Sub(reqLog.ReqTime))
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func FmtDuration(d time.Duration) string {
	ms := d.Milliseconds()
	d = d.Round(time.Millisecond)
	if ms < 10_000 {
		return fmt.Sprintf("%vms", ms)
	}
	return d.String()
}

func RewriteRequestAndBuildLog(newUrl *url.URL, proxyReq *http.Request, needStorage bool) (string, *ReqLog[Message]) {
	var reqLog *ReqLog[Message]
	var logStr string
	if needStorage {
		now := time.Now()
		id := uuid.New().String()
		if newUrl.Path == proxyReq.URL.Path {
			logStr = fmt.Sprintf("%v %s => %s", id, proxyReq.Host+proxyReq.URL.Path, newUrl.Host+" .")
		} else {
			logStr = fmt.Sprintf("%v %s => %s", id, proxyReq.Host+proxyReq.URL.Path, newUrl.Host+newUrl.Path)
		}

		bodyBt, body := ctool.CopyStream(proxyReq.Body)
		proxyReq.Body = body
		query, _ := url.QueryUnescape(proxyReq.URL.String())
		reqMes := Message{Header: proxyReq.Header, Body: FilterFormType(bodyBt)}

		id = fmt.Sprintf("%v%v", id[0:8], now.UnixMilli()%1000)
		cacheId := fmt.Sprintf("%v  %v", now.Format("01-02 15:04:05.000"), id)
		reqLog = &ReqLog[Message]{Id: id, CacheId: cacheId, Url: query, Request: reqMes, ReqTime: now, Method: proxyReq.Method}
	}

	proxyReq.Host = newUrl.Host
	//proxyReq.URL.Scheme = newUrl.Scheme
	proxyReq.URL.Host = newUrl.Host
	proxyReq.URL.Path = newUrl.Path
	//proxyReq.URL.RawQuery = newUrl.RawQuery
	return logStr, reqLog
}

// request body : start with "------"
func FilterFormType(s []byte) []byte {
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

func StartLog(mode string) {
	logger.Info("Redis key: ", RequestList, RequestUrlList)
	logger.Info("Start %v proxy server on 127.0.0.1:%d", mode, Port)
	logger.Warn("Pac: http://127.0.0.1:%d%v", ApiPort, PacUrl)
}

// HttpProxy HTTP代理和修改 HTTPS转发
func HttpProxy() {
	StartLog("HTTP")
	cert, err := GenCertificate()
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Addr:      fmt.Sprintf("0.0.0.0:%d", Port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProxyHandler(w, r)
		}),
	}

	StartAndCloseHook(server, func() error {
		StoreByMemory(ProxyConfVar)
		return nil
	})

	logger.Info("exit")
}

func StartAndCloseHook(server *http.Server, fns ...func() error) {
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err)
			os.Exit(0)
		}
	}()
	if fns == nil {
		fns = []func() error{}
	}
	fns = append(fns, func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	})
	Watch(fns...)
}

// Watch 监听进程收到的信号量
func Watch(fns ...func() error) {
	// 程序无法捕获信号 SIGKILL 和 SIGSTOP （终止和暂停进程），因此 os/signal 包对这两个信号无效。
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 阻塞
	s := <-ch
	close(ch)
	logger.Warn("catch signal", s.String())
	for i := range fns {
		if err := fns[i](); err != nil {
			log.Println(err)
		}
	}
}
