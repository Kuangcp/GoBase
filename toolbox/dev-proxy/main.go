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

	// TODO 替换proxy实现
	go app.HttpsProxy()

	go core.StartQueryServer()
	core.StartMainServer()
}
