package main

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
	"web_socket/common"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kuangcp/logger"
)

type ClientSession struct {
	id        string // client session id
	conn      *websocket.Conn
	writeLock *sync.Mutex
	quit      chan byte
	msgId     *int32 // msgId
}

func CreateClient(id, createWay string) *ClientSession {
	serverURL, finalID := buildServerURL(id)
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
	}
	conn, _, err := dialer.Dial(serverURL.String(), nil)
	if err != nil || conn == nil {
		logger.Warn(id+" [connect error]", err)
		return nil
	}

	logger.Info("%4s %s %s", id, createWay, ctool.Yellow.Print(finalID))
	if debugMode {
		conn.SetPongHandler(func(_ string) error {
			logger.Info("%4s [receive] pong-msg", id)
			return nil
		})
	}
	var msgId int32 = 0
	return &ClientSession{conn: conn, id: id, writeLock: &sync.Mutex{}, msgId: &msgId, quit: make(chan byte)}
}

func (client *ClientSession) close() {
	close(client.quit)
	closeConnect(client.conn)
}

func (client *ClientSession) SendRandomMsg() {
	if randomMsgSec <= 0 {
		return
	}

	randomMsg := time.NewTicker(time.Second * time.Duration(randomMsgSec))
	for {
		select {
		case <-quitAll:
			return
		case <-client.quit:
			return
		case t := <-randomMsg.C:
			err := common.SyncRun(client.writeLock, func() error {
				atomic.AddInt32(client.msgId, 1)
				return client.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(*client.msgId)))
			})
			if err != nil {
				logger.Info(client.id+" write:", err, t)
				return
			}
		}
	}
}

func (client *ClientSession) HandleReceive() {
	defer client.close()
	for {
		//err := conn.SetReadDeadline(time.Now().Add(time.Second * 60))
		//if err != nil {
		//	logger.Warn(id+" read:", err)
		//	return
		//}

		_, message, err := client.conn.ReadMessage()
		if err != nil {
			return
		}

		if debugMode {
			logger.Info("%4s [receive]\n%s", client.id, message)
		}
	}
}

// 心跳和退出信号
func (client *ClientSession) mainProcessLoop() {
	sum := 0
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	heartbeatTicker := time.NewTicker(heartbeat)
	defer heartbeatTicker.Stop()
	for {
		select {
		case <-quitAll:
			return
		case <-client.quit:
			return
		case t := <-heartbeatTicker.C:
			if debugMode {
				sum++
				logger.Info("%4s ping. counter=%d", client.id, sum)
			}
			err := common.SyncRun(client.writeLock, func() error {
				return client.conn.WriteMessage(websocket.PingMessage, nil)
			})
			if err != nil {
				logger.Error(client.id+" write:", err, t)
				return
			}
		case <-interrupt:
			logger.Debug(client.id + " close connection")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := client.conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Debug("%4s write close:", client.id, err)
				return
			}
		}
	}
}

func buildServerURL(id string) (url.URL, string) {
	finalID := id + "_" + uuid.New().String()[24:]

	wsServer := url.URL{
		Host:     serverHost,
		Path:     serverPath,
		Scheme:   "ws",
		RawQuery: "id=" + finalID + "&name=" + finalID,
	}

	return wsServer, finalID
}
