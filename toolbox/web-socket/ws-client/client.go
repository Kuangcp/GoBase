package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"time"

	_ "net/http/pprof"

	"github.com/kuangcp/logger"
	"github.com/kuangcp/sizedwaitgroup"

	"github.com/gorilla/websocket"
)

var (
	serverHost       = "localhost:7094" // ws-server config
	serverPath       = "/ws"
	maxCount         int
	debugMode        bool
	debugPort        string
	createNewDelayMs = 57 // 创建新连接之间的延迟
	randomMsgSec     = 0
)

const (
	printActiveCount = 30 * time.Second
	fillActiveClient = 10 * time.Second
	heartbeat        = 30 * time.Second
)

var activeClients int32 = 0
var hasExit = false
var latch *sizedwaitgroup.SizedWaitGroup
var quitAll = make(chan int)

func init() {
	flag.IntVar(&maxCount, "n", 5, "max count of client")

	flag.BoolVar(&debugMode, "D", false, "debug mode")
	flag.StringVar(&debugPort, "DP", "8891", "debug pprof port")

	flag.IntVar(&createNewDelayMs, "d", createNewDelayMs, "delay ms of create new connection")
	flag.IntVar(&randomMsgSec, "m", randomMsgSec, "random msg period(second)")

	flag.StringVar(&serverHost, "LH", serverHost, "direct host")
	flag.StringVar(&serverPath, "LP", serverPath, "direct path")

	logger.SetLogPathTrim("ws-client/")
	_ = logger.SetLoggerConfig(&logger.LogConfig{
		TimeFormat: logger.LogTimeDetailFormat,
	})
}

func main() {
	flag.Parse()
	if debugMode {
		go debug()
	}

	group, err := sizedwaitgroup.New(maxCount)
	if err != nil {
		logger.Error(err)
		return
	}
	latch = group

	if isInvalidServer() {
		return
	}

	go func() {
		createTotalClients()
		fillActiveClients()
	}()

	waitExitSignal()
}

func debug() {
	err := http.ListenAndServe("0.0.0.0:"+debugPort, nil)
	if err != nil {
		logger.Error(err)
	}
}

// 测试服务端是否可用
func isInvalidServer() bool {
	atomic.AddInt32(&activeClients, 1)
	client := CreateClient("0", "test")
	if client == nil {
		return true
	}

	client.close()
	return false
}

func createTotalClients() {
	for i := 0; i < maxCount; i++ {
		if hasExit {
			logger.Info("abort create new connection")
			break
		}
		id := strconv.Itoa(i + 1)
		startClient(id, "create")
		time.Sleep(time.Duration(createNewDelayMs) * time.Millisecond)
	}
}

func waitExitSignal() {
	ticker := time.NewTicker(printActiveCount)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-ticker.C:
			logger.Debug("active clients:", activeClients)
		case <-interrupt:
			hasExit = true
			close(quitAll)
			logger.Error("start exiting app")
			latch.Wait()
			return
		}
	}
}

func fillActiveClients() {
	ticker := time.NewTicker(fillActiveClient)
	for {
		select {
		case <-quitAll:
			return
		case <-ticker.C:
			for i := int(activeClients); i < maxCount; i++ {
				if hasExit {
					logger.Info("abort check live client")
					return
				}
				time.Sleep(time.Millisecond * time.Duration(createNewDelayMs))
				id := strconv.Itoa(i + 1)
				startClient(id+"F", "fill")
			}
		}
	}
}

func closeConnect(conn *websocket.Conn) {
	if conn != nil {
		err := conn.Close()
		if err != nil {
			logger.Error(err)
		}
	}
	atomic.AddInt32(&activeClients, -1)
}

func startClient(id, createWay string) {
	latch.Run(func() {
		client := CreateClient(id, createWay)
		if client == nil {
			return
		}

		atomic.AddInt32(&activeClients, 1)
		//defer closeConnect(client.conn)

		go client.HandleReceive()
		go client.SendRandomMsg()

		client.mainProcessLoop()
	})
}
