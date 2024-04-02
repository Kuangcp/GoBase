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

	web.InitClient()
	core.InitConnection()

	if core.HttpMode {
		core.Go(core.HttpProxy)
	} else {
		core.Go(app.HttpsProxy)
	}

	core.Go(web.StartQueryServer)

	systray.Run(OnReady, OnExit)
}
