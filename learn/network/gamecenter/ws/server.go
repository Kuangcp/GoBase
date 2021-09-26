package ws

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"

	"github.com/kuangcp/logger"

	"github.com/gorilla/websocket"
)

const (
	// 允许等待的写入时间
	writeWait = 10 * time.Second

	// 客户端心跳超时阈值
	heartbeatTimeout = 45 * time.Second
	maxConnect       = 100000

	// Maximum message size allowed from peer.
	maxMessageSize = 5120000
)

var (
	SilentLogMode bool // 静默日志
)
var (
	maxServerId int64                            // 连接ID，每次连接都加1
	allWsMap    = make(map[int64]*ServerSession) // ws 的所有连接 可用于广播
	allMapLock  = &sync.Mutex{}
)

var upgradeCnf = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有的CORS 跨域请求，正式环境可以关闭
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func simpleServerHandler(writer http.ResponseWriter, req *http.Request) {
	err := ActionRoundLock(allMapLock, func() error {
		if len(allWsMap) > maxConnect {
			return fmt.Errorf("out of limit")
		}
		return nil
	})
	if err != nil {
		return
	}

	// 应答客户端告知升级连接为websocket
	wsSocket, err := upgradeCnf.Upgrade(writer, req, nil)
	if err != nil {
		logger.Error("升级为websocket失败", err.Error())
		return
	}

	var wsConn *ServerSession
	SilentActionRoundLock(allMapLock, func() {
		maxServerId++
		wsConn = BuildNewSession(wsSocket, req)
		allWsMap[maxServerId] = wsConn
	})

	wsSocket.SetPingHandler(func(_ string) error {
		wsConn.lastHeartBeat = time.Now().UnixNano()
		return wsConn.wsWrite(websocket.PongMessage, nil)
	})

	// 业务处理
	go wsConn.processLoop()
	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()
}

func printTotal() {
	ticker := time.NewTicker(time.Second * 30)
	for range ticker.C {
		logger.Info("Online count: %5d", len(allWsMap))
	}
}

func startDebug() {
	fmt.Println("Debug: http://localhost:8891/debug/pprof/")
	err := http.ListenAndServe("0.0.0.0:8891", nil)
	if err != nil {
		logger.Error(err)
	}
}

func checkTimeOut() {
	ticker := time.NewTicker(time.Second * 5)
	for now := range ticker.C {
		SilentActionRoundLock(allMapLock, func() {
			now := now.UnixNano()
			for _, session := range allWsMap {
				session.hasTimeOut(now)
			}
		})
	}
}

func NewSimpleServer() {
	http.HandleFunc("/ws", simpleServerHandler)
	http.HandleFunc("/count",
		func(writer http.ResponseWriter, _ *http.Request) {
			writer.Write([]byte(fmt.Sprintf("%v", len(allWsMap))))
		})
	http.HandleFunc("/gc",
		func(writer http.ResponseWriter, _ *http.Request) {
			logger.Debug("start gc")
			runtime.GC()
			logger.Debug("end gc")
			writer.Write([]byte(fmt.Sprintf("OK")))
		})

	go startDebug()
	go printTotal()
	go checkTimeOut()
}
