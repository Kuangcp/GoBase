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
		{Short: "-p", Value: "file", Comment: "Print catalog"},
		{Short: "-r", Value: "file", Comment: "Remove catalog"},
		{Short: "-a", Value: "file", Comment: "Append catalog and title for file"},
		{Short: "-ra", Value: "file", Comment: "Remove then Append catalog and title for file. default options"},
		{Short: "-c", Value: "dir", Comment: "Refresh git repo dir changed file. same to -ra"},
		{Short: "-d", Value: "dir", Comment: "Refresh file catalog that recursive dir, default current dir. same to -ra"},
		{Short: "-xt", Value: "url", Comment: "Extract title from article url"},
	},
}

func init() {
	flag.BoolVar(&help, "h", false, "")
	flag.StringVar(&extractTitleUrl, "xt", "", "")

	// 组合使用
	optionToFunction("p", &printCatalog, PrintCatalog)

	optionToFunction("r", &rmFile, RemoveCatalog)
	optionToFunction("ra", &rmAppendFile, RefreshTagAndCatalog)

	optionToFunction("c", &refreshChangeDir, RefreshChangeFile)
	optionToFunction("d", &refreshDir, RefreshDirAllFiles)

	flag.Usage = info.PrintHelp
}

func optionToFunction(name string, val *string, action func(string)) {
	flag.StringVar(val, name, "", "")

	acts = append(acts, Act{val: val, action: action})
}
