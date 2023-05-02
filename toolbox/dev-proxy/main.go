package main

import "github.com/kuangcp/gobase/toolbox/dev-proxy/core"

func main() {
	go core.StartQueryServer()
	core.StartMainServer()
}
