package main

import (
	"flag"
	"github.com/kuangcp/gobase/pkg/ctool"
)

var (
	buildVersion string
	addRepo      string
	delRepo      string

	help     bool
	push     bool
	pull     bool
	allRepo  bool
	listRepo bool
)

var (
	cfgFile = ".ggi.ini"
)
var info = ctool.HelpInfo{
	Description:   "Manage multiple repository",
	Version:       "1.0.0",
	BuildVersion:  buildVersion,
	SingleFlagLen: -2,
	ValueLen:      -6,
	Flags: []ctool.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help"},
		{Short: "-pu", BoolVar: &push, Comment: "push repo"},
		{Short: "-pl", BoolVar: &pull, Comment: "pull repo"},
		{Short: "-l", BoolVar: &listRepo, Comment: "list repo"},
		{Short: "-all", BoolVar: &allRepo, Comment: "all repo"},
	},
	Options: []ctool.ParamVO{
		{Short: "-a", Value: "add", Comment: "add repo"},
		{Short: "-d", Value: "del", Comment: "del repo"},
	}}

func init() {

	flag.StringVar(&addRepo, "a", "", "")
	flag.StringVar(&delRepo, "d", "", "")

}
