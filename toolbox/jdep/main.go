package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kuangcp/gobase/pkg/ctool"
)

var (
	buildVersion string
	help         bool
	noDep        bool
	lessDep      bool
	projectRoot  string
)
var info = ctool.HelpInfo{
	Description:   "Maven 项目依赖分析：找出无依赖类、少依赖类等",
	Version:       "1.1.0",
	BuildVersion:  buildVersion,
	SingleFlagLen: -2,
	ValueLen:      -6,
	Flags: []ctool.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help"},
		{Short: "-n", BoolVar: &noDep, Comment: "找出未被任何其他类 import 或引用的类"},
		{Short: "-l", BoolVar: &lessDep, Comment: "lessDep"},
	},
	Options: []ctool.ParamVO{
		{Short: "-r", Value: "path", StringVar: &projectRoot, Comment: "Maven 项目根目录（默认当前目录）"},
	}}

func init() {
	//flag.StringVar(&projectRoot, "r", "", "")
}

func findNoDepClass() {
	root := projectRoot
	if root == "" {
		root = "."
	}
	list, err := FindNoDepClasses(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "扫描失败:", err)
		os.Exit(1)
	}
	if len(list) == 0 {
		fmt.Println("未找到无依赖类，或当前目录不是 Maven 项目根（无 pom.xml）。")
		return
	}
	fmt.Printf("共 %d 个类未被任何其他类 import 或引用：\n", len(list))
	for i, jf := range list {
		//fmt.Printf("  %s\n    %s\n", jf.FullClass, jf.Path)
		fmt.Printf("%4v %-80s\n\t\t%s\n", i+1, jf.FullClass, ctool.Yellow.Print(strings.Replace(jf.Path, root, "", -1)))
	}
}

func main() {
	info.Parse()
	if help {
		info.PrintHelp()
		return
	}
	if noDep {
		findNoDepClass()
		return
	}
	if lessDep {

	}
}
