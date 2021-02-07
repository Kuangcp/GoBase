package main

import (
	"flag"

	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/hosts-group/app"
)

func init() {
	flag.BoolVar(&app.Debug, "d", false, "")
}

func main() {
	flag.Parse()

	app.InitPrepare()

	go func() {
		app.WebServer("8066")
	}()

	systray.Run(app.OnReady, app.OnExit)
}
