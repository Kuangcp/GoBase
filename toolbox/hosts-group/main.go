package main

import (
	"flag"
	"fmt"

	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/hosts-group/app"
)

func init() {
	flag.BoolVar(&app.Debug, "d", false, "")
	flag.BoolVar(&app.Version, "v", false, "")
	flag.StringVar(&app.LogPath, "l", "", "")
	flag.Usage = app.Info.PrintHelp
}

func main() {
	flag.Parse()
	if app.Version {
		fmt.Println(app.Info.Version)
		return
	}

	app.InitConfigAndEnv()

	go func() {
		app.WebServer("8066")
	}()

	systray.Run(app.OnReady, app.OnExit)
}
