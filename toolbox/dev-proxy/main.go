package main

import "github.com/kuangcp/gobase/toolbox/dev-proxy/core"

func main() {
	core.HelpInfo.Parse()
	core.InitConfig()

	go core.StartQueryServer()
	core.StartMainServer()
}
