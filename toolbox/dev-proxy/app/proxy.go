package app

import (
	"crypto/tls"
	"fmt"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/goproxy"
	"github.com/kuangcp/logger"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var ProxyHandler *goproxy.Proxy

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
	receive := time.Now()
	ctx.Data["ReqCtx"] = &ReqCtx{
		receiveReqMs: receive.UnixMilli(),
	}
}

func (e *EventHandler) Auth(ctx *goproxy.Context, rw http.ResponseWriter) {
	// 身份验证
}

func (e *EventHandler) BeforeRequest(ctx *goproxy.Context) {
	// 设置X-Forwarded-For
	proxyReq := ctx.Req
	if clientIP, _, err := net.SplitHostPort(proxyReq.RemoteAddr); err == nil {
		if prior, ok := proxyReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		proxyReq.Header.Set("X-Forwarded-For", clientIP)
	}

	headers := core.GetHeaders(proxyReq.Host)
	//logger.Info(proxyReq.Host, headers)
	if headers != nil {
		for k, v := range headers {
			if k == "" {
				continue
			}
			proxyReq.Header.Set(k, v)
		}
	}

	proxyLog := ""
	var reqLog *core.ReqLog[core.Message]
	var proxyType = core.Direct
	needStorage := false
	hostPath := proxyReq.Host + proxyReq.URL.Path
	if core.TrackAllType || !core.IsMatch(core.StaticUrlPattern, hostPath) {
		//logger.Debug("none static %v", hostPath)
		var findNewUrl *url.URL
		findNewUrl, proxyType = core.FindReplaceByRegexp(*proxyReq)
		needStorage = core.MatchNeedStorage(*proxyReq, proxyType)
		if findNewUrl != nil {
			proxyLog, reqLog = core.RewriteRequestAndBuildLog(findNewUrl, proxyReq, needStorage)
		}
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

	reqCtx := ctx.Data["ReqCtx"].(*ReqCtx)
	reqCtx.startReqMs = time.Now().UnixMilli()
	reqCtx.reqLog = reqLog
	reqCtx.proxyLog = proxyLog
	reqCtx.needStorage = needStorage
	reqCtx.proxyType = proxyType
	reqCtx.originUrl = hostPath
}

func (e *EventHandler) BeforeResponse(ctx *goproxy.Context, resp *http.Response, err error) {
	if err != nil {
		return
	}
	rspTick := time.Now().UnixMilli()
	reqCtx := ctx.Data["ReqCtx"].(*ReqCtx)
	startTick := reqCtx.startReqMs
	reqMs := rspTick - startTick
	reqLog := reqCtx.reqLog
	resp.Header.Add("Ack-Proxy", reqCtx.proxyType+"  "+ctx.Req.Host)

	if reqCtx.proxyLog != "" {
		if reqCtx.proxyType == core.Proxy {
			reqCtx.proxyLog += " SELF"
		}
		logger.Debug("%4vms %v", reqMs, reqCtx.proxyLog)
	}
	if reqCtx.needStorage {
		core.TrySaveLog(reqLog, resp)
	}

	// 当前时刻 - 建立连接时刻 = 请求前处理耗时 + 请求&响应耗时 + 读取响应处理耗时
	handleRspTick := time.Now().UnixMilli()
	allMs := handleRspTick - reqCtx.receiveReqMs

	// 代理处理逻辑的耗时
	proxyReqMs := reqCtx.startReqMs - reqCtx.receiveReqMs
	proxyRspMs := handleRspTick - rspTick

	// 代理耗时 = 总耗时 - 请求&响应耗时
	proxyMs := allMs - reqMs
	if proxyMs > 10 {
		logger.Warn("SlowProxy proxy:%3vms=%v+%v total:%4v req:%v %v", proxyMs, proxyReqMs, proxyRspMs,
			allMs, reqMs, reqCtx.originUrl)
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
	msg := err.Error()
	if strings.Contains(msg, "context canceled") {
		return
	}
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
	core.StartLog("HTTPS")

	core.Go(func() {
		http.ListenAndServe("0.0.0.0:1255", nil)
	})

	// TODO 优化高并发下 transport 锁竞争问题
	// TODO 刚启动时延迟很低，跑了几千个压测后 延迟很高，怀疑代理逻辑有资源泄漏
	ProxyHandler = goproxy.New(goproxy.WithDecryptHTTPS(&Cache{}), goproxy.WithDelegate(&EventHandler{}))
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", core.Port),
		Handler:      ProxyHandler,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	core.StartAndCloseHook(server, func() error {
		core.StoreByMemory(core.ProxyConfVar)
		return nil
	})

	logger.Info("exit")
	os.Exit(0)
}
