package main

import (
	"github.com/kuangcp/gobase/toolbox/dev-proxy/app"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/web"
)

func main() {
	core.HelpInfo.Parse()
	core.InitConfig()
	web.InitClient()

	core.InitConnection()

	core.Go(web.StartQueryServer)

	if core.HttpMode {
		core.HttpProxy()
	} else {
		app.HttpsProxy()
	}
}
