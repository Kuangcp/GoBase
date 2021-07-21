package main

import (
	"embed"
	"flag"
	"fmt"

	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/hosts-group/app"
)

//go:embed static
var staticFS embed.FS

func init() {
	flag.BoolVar(&app.Debug, "d", false, "")
	flag.BoolVar(&app.DebugStatic, "D", false, "")
	flag.BoolVar(&app.Version, "v", false, "")
	flag.StringVar(&app.LogPath, "l", "", "")
	flag.StringVar(&app.FinalHostFile, "f", "", "")
	flag.StringVar(&app.MainPath, "m", "", "")
	flag.StringVar(&app.GenerateAfterCmd, "cmd", "", "")
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
		app.WebServer(staticFS, "8066")
	}()

	systray.Run(app.OnReady, app.OnExit)
}
