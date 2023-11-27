package main

import (
	"flag"
	"github.com/kuangcp/gobase/pkg/ctk"
	"os"
)

type Act struct {
	val    *string
	action func(string)
}

var acts []Act

func (a *Act) tryInvoke() {
	if *a.val != "" {
		a.action(*a.val)
		os.Exit(0)
	}
}

var info = ctk.HelpInfo{
	Description:   "Format markdown file, generate catalog",
	Version:       "1.0.5",
	BuildVersion:  buildVersion,
	SingleFlagLen: -3,
	DoubleFlagLen: -3,
	ValueLen:      -5,
	Flags: []ctk.ParamVO{
		{Short: "-h", Comment: "Help info"},
	},
	Options: []ctk.ParamVO{
		{Short: "", Value: "file", Comment: "Refresh file catalog"},
		{Short: "-d", Value: "dir", Comment: "Refresh file catalog that recursive dir, default current dir"},
		{Short: "-mm", Value: "file", Comment: "Print mind map"},
		{Short: "-r", Value: "file", Comment: "Remove catalog"},
		{Short: "-c", Value: "dir", Comment: "Refresh git repo dir changed file. same to -ra"},
		{Short: "-a", Value: "file", Comment: "Append catalog and title for file"},
		{Short: "-ra", Value: "file", Comment: "Remove then Append catalog and title for file. default options"},
	},
}

func init() {
	flag.BoolVar(&help, "h", false, "")

	optionToFunction("d", &refreshDir, refreshDirAllFiles)
	optionToFunction("mm", &mindMapFile, printMindMap)
	optionToFunction("c", &refreshChangeDir, refreshChangeFile)
	optionToFunction("a", &appendFile, createCatalog)
	optionToFunction("r", &rmFile, ignoreReturn(removeCatalog))
	optionToFunction("ra", &rmAppendFile, ReplaceThenRefreshCatalog)
	optionToFunction("p", &printCatalog, PrintCatalog)

	flag.Usage = info.PrintHelp
}

func ignoreReturn(act func(string) string) func(string) {
	return func(s string) {
		act(s)
	}
}

func optionToFunction(name string, val *string, action func(string)) {
	flag.StringVar(val, name, "", "")

	acts = append(acts, Act{val: val, action: action})
}
