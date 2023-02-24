package ctool

import "testing"

var (
	buildVersion string
	debug        bool
	port         int
)

var helpInfo = HelpInfo{
	Description:  " ",
	BuildVersion: buildVersion,
	Version:      "1.0.3",
	Flags: []ParamVO{
		{Short: "-d", BoolVar: &debug, Comment: "debug mode"},
	},
	Options: []ParamVO{
		{Short: "-p", IntVar: &port, Value: "port", Comment: "port"},
	},
}

func TestRun(t *testing.T) {
	helpInfo.Parse()
}
