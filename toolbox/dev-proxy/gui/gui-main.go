package main

import (
	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/app"

	//"github.com/kuangcp/gobase/toolbox/dev-proxy/app"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
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

	go core.StartQueryServer()

	systray.Run(OnReady, OnExit)
}
