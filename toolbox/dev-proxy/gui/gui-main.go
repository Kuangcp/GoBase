package main

import (
	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/app"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/web"
)

func main() {
	core.HelpInfo.Parse()
	core.InitConfig()
	core.MarkGuiMode()

	core.InitConnection()
	defer core.CloseConnection()

	if core.HttpMode {
		go core.HttpProxy()
	} else {
		go app.HttpsProxy()
	}

	go web.StartQueryServer()

	systray.Run(OnReady, OnExit)
}
