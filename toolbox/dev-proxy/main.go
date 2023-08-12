package main

import (
	"github.com/kuangcp/gobase/toolbox/dev-proxy/app"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
)

func main() {
	core.HelpInfo.Parse()
	core.InitConfig()

	core.InitConnection()
	defer core.CloseConnection()

	go core.StartQueryServer()

	if core.HttpMode {
		core.HttpProxy()
	} else {
		app.HttpsProxy()
	}
}
