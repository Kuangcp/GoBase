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

	if core.HttpProxy {
		go core.StartMainServer()
	} else {
		go app.HttpsProxy()
	}
}
