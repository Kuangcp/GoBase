package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/kuangcp/logger"
	"github.com/kuangcp/gobase/pkg/ctool"
	"strings"
)


var (
	buildVersion string

	help         bool
	long         bool
)

var info = ctool.HelpInfo{
	Description:   "Start static file web server on current path",
	Version:       "1.1.0",
	BuildVersion:  buildVersion,
	SingleFlagLen: -2,
	ValueLen:      -6,
	Flags: []ctool.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help"},
		{Short: "-l", BoolVar: &long, Comment: "join long"},
	},
	Options: []ctool.ParamVO{
		// {Short: "-p", Value: "port", Comment: "web server port"},
		// {Short: "-d", Value: "folder", Comment: "folder pair. like -d x=y "},
	}}

func main() {
	info.Parse()
	if help {
		info.PrintHelp()
		return
	}
	if long {
		concatNumberLong()
		return
	}

	concatNumber()
}

func concatNumberLong() {
	lines := readClipList()
	result := strings.Join(lines, "L,")
	result += "L"
	fmt.Println(result)
	err := clipboard.WriteAll(result)
	if err != nil {
		logger.Fatal(err)
	}
}

func concatNumber() {
	lines := readClipList()
	result := strings.Join(lines, ",")
	fmt.Println(result)
	err := clipboard.WriteAll(result)
	if err != nil {
		logger.Fatal(err)
	}
}

func readClipList() []string {
	last, err := clipboard.ReadAll()
	if err != nil {
		logger.Fatal(err)
	}

	last = strings.TrimSpace(last)
	if last == "" {
		logger.Fatal("Empty")
	}
	lines := strings.Split(last, "\n")
	return lines
}
