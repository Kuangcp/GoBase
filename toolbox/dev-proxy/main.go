package main

import "github.com/kuangcp/gobase/toolbox/dev-proxy/core"

func main() {
	core.HelpInfo.Parse()

	go core.StartQueryServer()
	core.StartMainServer()
}
