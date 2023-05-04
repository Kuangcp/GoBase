package main

import (
	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
)

func main() {
	core.HelpInfo.Parse()

	go core.StartQueryServer()
	go core.StartMainServer()

	systray.Run(OnReady, OnExit)
}
