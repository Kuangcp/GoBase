package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
	"web_socket/common"

	"github.com/gorilla/websocket"
	"github.com/kuangcp/logger"
)

// 客户端读写消息
type websocketMsg struct {
	// websocket.TextMessage 消息类型
	messageType int
	data        []byte
}

// ServerSession 服务端会话连接
type ServerSession struct {
	wsSocket *websocket.Conn    // 底层websocket
	inChan   chan *websocketMsg // 读队列
	outChan  chan *websocketMsg // 写队列

	mutex         sync.Mutex // 避免重复关闭管道,加锁处理
	isClosed      bool
	closeChan     chan byte // 连接关闭的通知
	serverId      int64
	clientId      string
	lastHeartBeat int64
}

func BuildNewSession(wsSocket *websocket.Conn, req *http.Request) *ServerSession {
	wxId := req.FormValue("wxId")
	if wxId == "" {
		wxId = fmt.Sprint(maxServerId)
	}
	wsConn := &ServerSession{
		serverId:      maxServerId,
		clientId:      wxId,
		wsSocket:      wsSocket,
		inChan:        make(chan *websocketMsg, 500),
		outChan:       make(chan *websocketMsg, 500),
		closeChan:     make(chan byte),
		lastHeartBeat: time.Now().UnixNano(),
		isClosed:      false,
	}
	if !silentLogMode {
		logger.Info("%s connected", wsConn.ID())
	}
	return wsConn
}

func (wsConn *ServerSession) ID() string {
	return fmt.Sprint(wsConn.serverId, "#", wsConn.clientId)
}

// 处理消息队列中的消息
func (wsConn *ServerSession) ReadLoop() {
	// 设置消息的最大长度
	wsConn.wsSocket.SetReadLimit(maxMessageSize)
	//wsConn.wsSocket.SetReadDeadline(time.Now().Add(pongWait))
	for {
		// 读一个message
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
			logger.Info("%s [read] %v", wsConn.ID(), err.Error())
			wsConn.closeThenClean()
			return
		}
		//fmt.Println(wsConn.serverId, string(data))
		req := &websocketMsg{msgType, data}

		// 放入消息队列
		select {
		case wsConn.inChan <- req:
		case <-wsConn.closeChan:
			return
		}
	}
}

// 发送消息给客户端
func (wsConn *ServerSession) WriteLoop() {
	//ticker := time.NewTicker(pingPeriod)
	//defer func() {
	//	ticker.Stop()
	//}()
	for {
		select {
		case msg := <-wsConn.outChan:
			// 写给websocket
			if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
				logger.Error("%s send msg to client error %v", wsConn.ID(), err.Error())
				wsConn.closeThenClean()
				return
			}
		case <-wsConn.closeChan:
			// 获取到关闭通知
			return
			//case <-ticker.C:
			//	// 出现超时情况
			//	wsConn.wsSocket.SetWriteDeadline(time.Now().Add(writeWait))
			//	if err := wsConn.wsSocket.WriteMessage(websocket.PingMessage, nil); err != nil {
			//		return
			//	}
		}
	}
}

// 处理队列中的消息
func (wsConn *ServerSession) ProcessLoop() {
	for {
		msg, err := wsConn.wsRead()
		if err != nil {
			//logger.Error("%s process msg error %v", wsConn.ID(), err.Error())
			break
		}

		//log.Println(wsConn.serverId, "接收到消息", msgStr)
		wsConn.handleReceiveMsg(msg)
	}
}

// 写入消息到队列中
func (wsConn *ServerSession) wsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &websocketMsg{messageType, data}:
	case <-wsConn.closeChan:
		return errors.New("[write] connection has closed")
	}
	return nil
}

// 读取消息队列中的消息
func (wsConn *ServerSession) wsRead() (*websocketMsg, error) {
	select {
	case msg := <-wsConn.inChan:
		// 获取到消息队列中的消息
		return msg, nil
	case <-wsConn.closeChan:
		return nil, errors.New("[read] connection has closed")
	}
}

// 关闭连接
func (wsConn *ServerSession) closeThenClean() {
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()

	if wsConn.wsSocket != nil {
		err := wsConn.wsSocket.Close()
		if err != nil {
			logger.Error(err)
		}
	}

	if wsConn.isClosed == false {
		wsConn.isClosed = true

		// 删除这个连接的变量
		common.SyncRuns(allMapLock, func() {
			delete(allWsMap, wsConn.serverId)
		})

		close(wsConn.closeChan)
	}
}

func (wsConn *ServerSession) handleReceiveMsg(msg *websocketMsg) {
	msgStr := string(msg.data)
	intVal, err := strconv.Atoi(msgStr)
	var rsp []byte
	if err == nil {
		rsp = []byte(fmt.Sprint(intVal + 1))
	} else {
		rsp = []byte(buildVirtualRsp(msgStr))
	}

	wsConn.outChan <- &websocketMsg{messageType: msg.messageType, data: rsp}
}

func buildVirtualRsp(msg string) string {
	result := ""
	for i := 0; i < 3; i++ {
		sum := md5.Sum([]byte(fmt.Sprint(msg, i)))
		result += msg + "  " + fmt.Sprintf("%x", sum) + "\n"
	}
	return result
}

func (wsConn *ServerSession) hasTimeOut(now int64) {
	if wsConn.isClosed {
		return
	}
	if time.Duration(now-wsConn.lastHeartBeat) > heartbeatTimeout {
		logger.Warn("%s client heartbeat has bean timeout !!!", wsConn.ID())
		wsConn.closeThenClean()
	}
}
