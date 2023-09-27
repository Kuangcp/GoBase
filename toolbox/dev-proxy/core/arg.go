package core

import (
	"github.com/kuangcp/gobase/pkg/ctool"
)

var (
	Port         int
	ApiPort      int
	ReloadConf   bool
	TrackAllType bool
	Debug        bool
	BuildVersion string
	HttpMode     bool
	JsonPath     string
	PacPath      string
)

var HelpInfo = ctool.HelpInfo{
	Description:  "Http proxy for reroute and trace",
	BuildVersion: BuildVersion,
	Version:      "1.0.4",
	Flags: []ctool.ParamVO{
		{Short: "-r", BoolVar: &ReloadConf, Comment: "auto reload changed config"},
		{Short: "-d", BoolVar: &Debug, Comment: "debug mode"},
		{Short: "-x", BoolVar: &HttpMode, Comment: "track or modify http, capture https. (default https mode, need install cert)"},
		{Short: "-A", BoolVar: &TrackAllType, Comment: "track all request: js,css,html,image"},
	},
	Options: []ctool.ParamVO{
		{Short: "-w", IntVar: &ApiPort, Int: 1235, Value: "port", Comment: "web api port"},
		{Short: "-p", IntVar: &Port, Int: 1234, Value: "port", Comment: "proxy port"},
		{Short: "-j", StringVar: &JsonPath, String: "", Value: "path", Comment: "json config file abs path"},
		{Short: "-a", StringVar: &PacPath, String: "", Value: "path", Comment: "pac file abs path"},
	},
}
