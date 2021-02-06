package main

import (
	"flag"
	"github.com/kuangcp/gobase/toolbox/hosts-group/app"
)

func init() {
	flag.BoolVar(&app.Debug, "d", false, "")
}

func main() {
	flag.Parse()

	app.InitPrepare()

	app.WebServer("8066")
}
