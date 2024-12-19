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

	help     bool
	push     bool
	pull     bool
	allRepo  bool
	listRepo bool
)

var info = ctool.HelpInfo{
	Description:   "Manage multiple repository",
	Version:       "1.0.0",
	BuildVersion:  buildVersion,
	SingleFlagLen: -4,
	ValueLen:      -6,
	Flags: []ctool.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help"},
		{Short: "-pu", BoolVar: &push, Comment: "push repo"},
		{Short: "-pl", BoolVar: &pull, Comment: "pull repo"},
		{Short: "-l", BoolVar: &listRepo, Comment: "list repo"},
		{Short: "-all", BoolVar: &allRepo, Comment: "all repo"},
	},
	Options: []ctool.ParamVO{
		{Short: "-a", Value: "alias", Comment: "add repo"},
		{Short: "-d", Value: "alias", Comment: "del repo"},
		{Short: "-j", Value: "alias", Comment: "jump repo"},
	}}

func init() {
	flag.StringVar(&addRepo, "a", "", "")
	flag.StringVar(&delRepo, "d", "", "")
	flag.StringVar(&jumpRepo, "j", "", "")

}
