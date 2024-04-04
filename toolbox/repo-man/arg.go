package main

import "github.com/kuangcp/gobase/pkg/ctool"

var (
	list    bool
	fetch   bool
	push    bool
	pushAll bool // pushAll all remote
	status  bool
	repo    string
)
var helpInfo = ctool.HelpInfo{
	Description: "Repo management",
	//BuildVersion: BuildVersion,
	Version: "1.0.0",
	Flags: []ctool.ParamVO{
		{Short: "-l", BoolVar: &list, Comment: "list all repo"},
		{Short: "-f", BoolVar: &fetch, Comment: "fetch remote upstream"},
		{Short: "-p", BoolVar: &push, Comment: "push to remote"},
		{Short: "-pa", BoolVar: &pushAll, Comment: "push to remote"},
		{Short: "-s", BoolVar: &status, Comment: "show status"},
	},
	Options: []ctool.ParamVO{
		{Short: "-r", StringVar: &repo, String: "", Value: "repo", Comment: "target repo"},
	},
}
