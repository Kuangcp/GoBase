package main

import (
	"github.com/kuangcp/gobase/toolbox/dev-proxy/app"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
)

func main() {
	core.HelpInfo.Parse()
	core.InitConfig()

	go app.HttpsProxy()

	go core.StartQueryServer()
	core.StartMainServer()
}
