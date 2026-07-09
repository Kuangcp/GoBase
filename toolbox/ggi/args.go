package main

import (
	"flag"
	"github.com/kuangcp/gobase/pkg/ctool"
)

var (
	buildVersion string
	addRepo      string
	delRepo      string
	jumpRepo     string
	listRepo     bool
	help         bool

	push        bool
	pushCur     bool
	pull        bool
	allRepo     bool
	checkChange bool

	lod        bool
	lods       bool
	afterDate  string
	beforeDate string
)

var info = ctool.HelpInfo{
	Description:   "Manage multiple repository",
	Version:       "1.0.0",
	BuildVersion:  buildVersion,
	SingleFlagLen: -4,
	ValueLen:      -6,
	Flags: []ctool.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help"},
		{Short: "-p", BoolVar: &pushCur, Comment: "push current repo to all remotes"},
		{Short: "-pu", BoolVar: &push, Comment: "push ahead repos"},
		{Short: "-pl", BoolVar: &pull, Comment: "pull all config repo"},
		{Short: "-s", BoolVar: &checkChange, Comment: "check repo change (default)"},
		{Short: "-l", BoolVar: &listRepo, Comment: "list all repo"},
		{Short: "-all", BoolVar: &allRepo, Comment: "list all repo with status"},
		{Short: "-lod", BoolVar: &lod, Comment: "list commit counts per subdir"},
		{Short: "-lods", BoolVar: &lods, Comment: "list commit counts per subdir sorted"},
	},
	Options: []ctool.ParamVO{
		{Short: "-a", Value: "alias", Comment: "add repo"},
		{Short: "-d", Value: "alias", Comment: "del repo"},
		{Short: "-j", Value: "alias", Comment: "jump repo"},
		{Short: "--after", Value: "date", Comment: "filter commits after date"},
		{Short: "--before", Value: "date", Comment: "filter commits before date"},
	}}

func init() {
	flag.StringVar(&addRepo, "a", "", "")
	flag.StringVar(&delRepo, "d", "", "")
	flag.StringVar(&jumpRepo, "j", "", "")
	flag.StringVar(&afterDate, "after", "", "")
	flag.StringVar(&beforeDate, "before", "", "")

}
