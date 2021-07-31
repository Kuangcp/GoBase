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
	flag.BoolVar(&app.DebugHostFile, "d", false, "")
	flag.BoolVar(&app.DebugStatic, "D", false, "")
	flag.IntVar(&app.Port, "p", 8066, "")
	flag.BoolVar(&app.Version, "v", false, "")
	flag.StringVar(&app.LogPath, "l", "", "")
	flag.StringVar(&app.FinalHostFile, "f", "", "")
	flag.StringVar(&app.MainPath, "m", "", "")
	flag.StringVar(&app.ChangeFileHook, "hook", "", "")
	flag.StringVar(&app.SupportMode, "mode", "", "")
	flag.Usage = app.Info.PrintHelp
}

func main() {
	flag.Parse()
	app.PortStr = fmt.Sprint(app.Port)
	if app.Version {
		fmt.Println(app.Info.Version)
		return
	}

	app.InitConfigBuildEnv()

	go func() {
		app.WebServer(staticFS, app.PortStr)
	}()

	systray.Run(app.OnReady, app.OnExit)
}
