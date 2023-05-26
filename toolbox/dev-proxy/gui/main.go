package main

import (
	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/app"
)

func main() {
	core.HelpInfo.Parse()
	core.InitConfig()

	// TODO 替换proxy实现
	go app.HttpsProxy()

	go core.StartQueryServer()
	go core.StartMainServer()

	systray.Run(OnReady, OnExit)
}
