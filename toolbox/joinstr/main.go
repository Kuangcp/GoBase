package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/ctool/stream"
	"github.com/kuangcp/logger"
	"strings"
)

var (
	buildVersion string

	help bool
	long bool
	str  bool
	strs bool
)

var info = ctool.HelpInfo{
	Description:   "read clip , join , write clip",
	Version:       "1.1.1",
	BuildVersion:  buildVersion,
	SingleFlagLen: -2,
	ValueLen:      -6,
	Flags: []ctool.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help"},
		{Short: "-l", BoolVar: &long, Comment: "join long"},
		{Short: "-s", BoolVar: &str, Comment: "join single quote string"},
		{Short: "-d", BoolVar: &strs, Comment: "join double quote string"},
	},
	Options: []ctool.ParamVO{}}

func main() {
	info.Parse()
	if help {
		info.PrintHelp()
		return
	}
	if long {
		concatStream(func(s stream.Stream) string {
			return stream.ToJoins(
				s.Map(func(item any) any {
					return item.(string) + "L"
				}), ",")
		})
		return
	}
	if str {
		concatStream(func(s stream.Stream) string {
			return stream.ToJoins(
				s.Map(func(item any) any {
					return "'" + item.(string) + "'"
				}), ",")
		})
		return
	}
	if strs {
		concatStream(func(s stream.Stream) string {
			return stream.ToJoins(
				s.Map(func(item any) any {
					return "\"" + item.(string) + "\""
				}), ",")
		})
		return
	}

	concatStream(func(s stream.Stream) string {
		return stream.ToJoins(s, ",")
	})
}

func concatStream(fun func(s stream.Stream) string) {
	lines := readClipList()
	from := stream.From(func(source chan<- any) {
		for _, s := range lines {
			source <- s
		}
	})
	result := fun(from)
	//result := stream.ToJoins(from, ",")
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
