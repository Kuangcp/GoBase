package main

import (
	"flag"

	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/gobase/toolbox/hosts-group/app"
)

var info = cuibase.HelpInfo{
	Description:   "Hosts switch",
	Version:       "1.0.2",
	SingleFlagLen: -2,
	DoubleFlagLen: 0,
	ValueLen:      -5,
	Flags: []cuibase.ParamVO{
		{Short: "-h", Comment: "help info"},
		{Short: "-d", Comment: "debug"},
	},
	Options: []cuibase.ParamVO{
		{Short: "-l", Value: "path", Comment: "log path"},
	},
}

func init() {
	flag.BoolVar(&app.Debug, "d", false, "")
	flag.StringVar(&app.LogPath, "l", "", "")
	flag.Usage = info.PrintHelp
}

func main() {
	flag.Parse()

	app.InitConfigAndEnv()

	go func() {
		app.WebServer("8066")
	}()

	systray.Run(app.OnReady, app.OnExit)
}
