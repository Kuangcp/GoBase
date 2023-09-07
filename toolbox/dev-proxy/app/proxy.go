package app

import (
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/logger"
	"github.com/ouqiang/goproxy"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type EventHandler struct{}

func (e *EventHandler) Connect(ctx *goproxy.Context, rw http.ResponseWriter) {
	// 保存的数据可以在后面的回调方法中获取
	//ctx.Data["req_id"] = "uuid"

	// 禁止访问某个域名
	//if strings.Contains(ctx.Req.URL.Host, "example.com") {
	//	rw.WriteHeader(http.StatusForbidden)
	//	ctx.Abort()
	//	return
	//}
}

func (e *EventHandler) Auth(ctx *goproxy.Context, rw http.ResponseWriter) {
	// 身份验证
}

func (e *EventHandler) BeforeRequest(ctx *goproxy.Context) {
	// 修改header
	//ctx.Req.Header.Add("X-Request-Id", ctx.Data["req_id"].(string))

	// 设置X-Forwarded-For
	proxyReq := ctx.Req
	if clientIP, _, err := net.SplitHostPort(proxyReq.RemoteAddr); err == nil {
		if prior, ok := proxyReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		proxyReq.Header.Set("X-Forwarded-For", clientIP)
	}
	proxyLog := ""
	var reqLog *core.ReqLog[core.Message]
	findNewUrl, proxyType := core.FindReplaceByRegexp(*proxyReq)
	needStorage := core.MatchNeedStorage(*proxyReq)
	if findNewUrl != nil {
		proxyLog, reqLog = core.RewriteRequestAndBuildLog(findNewUrl, proxyReq, needStorage)
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

	now := time.Now()
	ctx.Data["ReqCtx"] = &ReqCtx{
		reqLog:      reqLog,
		proxyLog:    proxyLog,
		needStorage: needStorage,
		proxyType:   proxyType,
		startMs:     now.UnixMilli(),
	}
}

func (e *EventHandler) BeforeResponse(ctx *goproxy.Context, resp *http.Response, err error) {
	if err != nil {
		return
	}
	reqCtx := ctx.Data["ReqCtx"].(*ReqCtx)
	reqLog := reqCtx.reqLog

	resp.Header.Add("Ack", "dev-proxy: "+reqCtx.proxyType)

	startMs := reqCtx.startMs

	bodyBt, body := core.CopyStream(resp.Body)
	resp.Body = body

	resMes := core.Message{Header: resp.Header, Body: bodyBt}

	endMs := time.Now().UnixMilli()
	waste := endMs - startMs

	if reqCtx.proxyLog != "" {
		if reqCtx.proxyType == core.Proxy {
			reqCtx.proxyLog += " SELF"
		}
		logger.Debug("%4vms %v", waste, reqCtx.proxyLog)
	}

	core.HandleCompressed(&resMes, resp)

	if reqCtx.needStorage && reqCtx.proxyType != core.Direct && reqLog != nil {
		core.FillReqLogResponse(reqLog, resp)
		// redis cache
		core.Conn.ZAdd(core.RequestList, redis.Z{Member: reqLog.CacheId, Score: float64(reqLog.ReqTime.UnixNano())})
		core.Conn.HSet(core.RequestUrlList, reqLog.Id, reqLog.Url)

		core.SaveReqLog(reqLog)
	}
}

// ParentProxy 设置上级代理
func (e *EventHandler) ParentProxy(req *http.Request) (*url.URL, error) {
	//return url.Parse("http://localhost:1087")
	return nil, nil
}

func (e *EventHandler) Finish(ctx *goproxy.Context) {
}

// ErrorLog 记录错误日志
func (e *EventHandler) ErrorLog(err error) {
	logger.Error(err)
}

func (e *EventHandler) WebSocketSendMessage(ctx *goproxy.Context, messageType *int, p *[]byte) {
	//TODO implement me
}

func (e *EventHandler) WebSocketReceiveMessage(ctx *goproxy.Context, messageType *int, p *[]byte) {
	//TODO implement me
}

// Cache 实现证书缓存接口
type Cache struct {
	m sync.Map
}

func (c *Cache) Set(host string, cert *tls.Certificate) {
	c.m.Store(host, cert)
}
func (c *Cache) Get(host string) *tls.Certificate {
	v, ok := c.m.Load(host)
	if !ok {
		return nil
	}

	return v.(*tls.Certificate)
}

// HttpsProxy replaced core.HttpProxy HTTP HTTPS 代理修改，密文解密
func HttpsProxy() {
	logger.Info("list key: ", core.RequestList)
	logger.Info("Start HTTPS proxy server on 127.0.0.1:%d", core.Port)
	logger.Warn("Pac: http://127.0.0.1:%d%v", core.ApiPort, core.PacUrl)

	proxy := goproxy.New(goproxy.WithDecryptHTTPS(&Cache{}), goproxy.WithDelegate(&EventHandler{}))
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", core.Port),
		Handler:      proxy,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
