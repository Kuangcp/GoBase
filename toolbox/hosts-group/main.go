package main

import (
	"embed"
	"flag"
	"fmt"

	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/hosts-group/app"
	"github.com/webview/webview"
)

//go:embed static
var staticFS embed.FS

func init() {
	flag.BoolVar(&app.Win, "w", false, "")
	
	flag.BoolVar(&app.DebugHostFile, "d", false, "")
	flag.BoolVar(&app.DebugStatic, "D", false, "")
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
	if app.Version {
		fmt.Println(app.Info.Version)
		return
	}

	app.InitConfigBuildEnv()

	if app.Win {
		w := webview.New(false)
		defer w.Destroy()
		w.SetTitle("Hosts group")
		w.SetSize(1035, 650, webview.HintFixed)
		w.Navigate("http://localhost:8066/")
		w.Run()
		return
	}

	go func() {
		app.WebServer(staticFS, "8066")
	}()

	systray.Run(app.OnReady, app.OnExit)
}
