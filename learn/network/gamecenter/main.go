package main

import (
	"flag"
	"fmt"
	"gamecenter/ws"
	"net/http"

	"github.com/kuangcp/logger"
)

func init() {
	flag.BoolVar(&ws.SilentLogMode, "s", false, "silent log")

	logger.SetLogPathTrim("ws-server/")
	_ = logger.SetLoggerConfig(&logger.LogConfig{
		TimeFormat: logger.LogTimeDetailFormat,
	})
}

func main() {
	flag.Parse()
	addrPort := 6062

	ws.NewSimpleServer()

	fmt.Println("Start on port:", addrPort)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", addrPort), nil)
}
