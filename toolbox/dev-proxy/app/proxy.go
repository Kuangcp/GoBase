package app

import (
	"crypto/tls"
	"fmt"
	"github.com/google/uuid"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/ouqiang/goproxy"
	"log"
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
	if clientIP, _, err := net.SplitHostPort(ctx.Req.RemoteAddr); err == nil {
		if prior, ok := ctx.Req.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		ctx.Req.Header.Set("X-Forwarded-For", clientIP)
	}
	bodyBt, body := core.CopyStream(ctx.Req.Body)
	ctx.Req.Body = body

	now := time.Now()
	id := uuid.New().String()
	query, _ := url.QueryUnescape(ctx.Req.URL.String())
	id = fmt.Sprintf("%v%v", id[0:8], now.UnixMilli()%1000)
	cacheId := fmt.Sprintf("%v  %v", now.Format("01-02 15:04:05.000"), id)
	reqMes := core.Message{Header: ctx.Req.Header, Body: bodyBt}
	ctx.Data["RLog"] = &core.ReqLog[core.Message]{Id: id, CacheId: cacheId, Url: query, Request: reqMes, ReqTime: now, Method: ctx.Req.Method}
}

func (e *EventHandler) BeforeResponse(ctx *goproxy.Context, resp *http.Response, err error) {
	if err != nil {
		return
	}
	v := ctx.Data["RLog"]
	reqLog := v.(*core.ReqLog[core.Message])
	bodyBt, body := core.CopyStream(resp.Body)
	resp.Body = body

	resMes := core.Message{Header: resp.Header, Body: bodyBt}

	core.HandleCompressed(&resMes, resp)

	reqLog.Response = resMes
	reqLog.ResTime = time.Now()
	reqLog.ElapsedTime = core.FmtDuration(reqLog.ResTime.Sub(reqLog.ReqTime))
	reqLog.Status = resp.Status
	reqLog.StatusCode = resp.StatusCode
	// TODO save leveldb redis
}

// ParentProxy 设置上级代理
func (e *EventHandler) ParentProxy(req *http.Request) (*url.URL, error) {
	//return url.Parse("http://localhost:1087")
	return nil, nil
}

func (e *EventHandler) Finish(ctx *goproxy.Context) {
	reqLog := ctx.Data["RLog"].(*core.ReqLog[core.Message])

	for k, v := range reqLog.Response.Header {
		fmt.Println(k, "<=>", v)
	}
	fmt.Println(string(reqLog.Response.Body))
	fmt.Printf("请求结束 URL:%s %v\n", ctx.Req.URL, reqLog)
}

// ErrorLog 记录错误日志
func (e *EventHandler) ErrorLog(err error) {
	log.Println(err)
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

func HttpsProxy() {
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
