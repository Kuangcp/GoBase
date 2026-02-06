package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/kuangcp/gobase/pkg/ctool"
)

var (
	buildVersion      string
	help              bool
	noDep             bool
	lessDep           bool
	duplicateClass    bool
	duplicateDiff     bool
	rmDuplicateClass  bool
	projectRoot       string
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
		{Short: "-d", BoolVar: &duplicateClass, Comment: "重复类名"},
		{Short: "-dd", BoolVar: &duplicateDiff, Comment: "重复类名并输出与首文件的 git 风格 diff（-红 +绿）"},
		{Short: "-rd", BoolVar: &rmDuplicateClass, Comment: "删除重复类名"},
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

func findLessDepClass() {

}

func findDuplicateClass() {
	root := projectRoot
	if root == "" {
		root = "."
	}
	dup, err := FindDuplicateClasses(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "扫描失败:", err)
		os.Exit(1)
	}
	if len(dup) == 0 {
		fmt.Println("未发现重名类，或当前目录不是 Maven 项目根（无 pom.xml）。")
		return
	}
	fmt.Printf("共 %d 个重名类：\n", len(dup))
	names := make([]string, 0, len(dup))
	for name := range dup {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		list := dup[name]
		fmt.Printf("\n  %s（%d 处）\n", name, len(list))
		for _, jf := range list {
			path := strings.Replace(jf.Path, root, "", 1)
			if path == jf.Path {
				path = jf.Path
			}
			fmt.Printf("    %s\n", ctool.Yellow.Print(path))
		}
		identical, diffPcts := CompareDuplicateContents(list)
		if identical {
			fmt.Printf("    %s\n", ctool.Green.Print("类定义完全一致"))
		} else {
			parts := make([]string, 0, len(diffPcts)-1)
			for i := 1; i < len(diffPcts); i++ {
				parts = append(parts, fmt.Sprintf("%.0f%%", diffPcts[i]))
			}
			fmt.Printf("    %s %s\n", ctool.Yellow.Print("类定义不一致，与首文件差异约"), strings.Join(parts, "、"))
		}
	}
}

func findDuplicateClassWithDiff() {
	root := projectRoot
	if root == "" {
		root = "."
	}
	dup, err := FindDuplicateClasses(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "扫描失败:", err)
		os.Exit(1)
	}
	if len(dup) == 0 {
		fmt.Println("未发现重名类，或当前目录不是 Maven 项目根（无 pom.xml）。")
		return
	}
	names := make([]string, 0, len(dup))
	for name := range dup {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		list := dup[name]
		fmt.Printf("\n  %s（%d 处）\n", name, len(list))
		for _, jf := range list {
			path := strings.Replace(jf.Path, root, "", 1)
			if path == jf.Path {
				path = jf.Path
			}
			fmt.Printf("    %s\n", ctool.Yellow.Print(path))
		}
		identical, diffPcts := CompareDuplicateContents(list)
		if identical {
			fmt.Printf("    %s\n", ctool.Green.Print("类定义完全一致"))
		} else {
			parts := make([]string, 0, len(diffPcts)-1)
			for i := 1; i < len(diffPcts); i++ {
				parts = append(parts, fmt.Sprintf("%.0f%%", diffPcts[i]))
			}
			fmt.Printf("    %s %s\n", ctool.Yellow.Print("类定义不一致，与首文件差异约"), strings.Join(parts, "、"))
		}
		// 与首文件的 git 风格 diff：第 2、3、4… 个文件分别对第一个做 diff
		firstPath := list[0].Path
		for i := 1; i < len(list); i++ {
			other := list[i]
			diffs, err := DiffTwoFiles(firstPath, other.Path)
			if err != nil {
				fmt.Printf("    %s\n", ctool.Red.Print("diff 读取失败: "+err.Error()))
				continue
			}
			relFirst := strings.Replace(firstPath, root, "", 1)
			if relFirst == firstPath {
				relFirst = firstPath
			}
			relOther := strings.Replace(other.Path, root, "", 1)
			if relOther == other.Path {
				relOther = other.Path
			}
			fmt.Printf("    --- %s\n", relFirst)
			fmt.Printf("    +++ %s\n", relOther)
			for _, d := range diffs {
				switch d.Op {
				case "-":
					fmt.Printf("    %s\n", ctool.Red.Print("- "+d.Line))
				case "+":
					fmt.Printf("    %s\n", ctool.Green.Print("+ "+d.Line))
				default:
					fmt.Printf("      %s\n", d.Line)
				}
			}
		}
	}
}

func deleteDuplicateClass() {
	root := projectRoot
	if root == "" {
		root = "."
	}
	dup, err := FindDuplicateClasses(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "扫描失败:", err)
		os.Exit(1)
	}
	if len(dup) == 0 {
		fmt.Println("未发现重名类，或当前目录不是 Maven 项目根（无 pom.xml）。")
		return
	}
	names := make([]string, 0, len(dup))
	for name := range dup {
		names = append(names, name)
	}
	sort.Strings(names)

	var allModified []string
	var allDeleted []string
	for _, name := range names {
		list := dup[name]
		identical, _ := CompareDuplicateContents(list)
		if !identical {
			continue
		}
		sorted := SortDuplicateListForKeep(root, list)
		kept := sorted[0]
		toDelete := sorted[1:]
		excludePaths := make(map[string]struct{})
		for _, jf := range toDelete {
			excludePaths[jf.Path] = struct{}{}
		}
		for _, jf := range toDelete {
			modified, err := ReplaceImportInProject(root, jf.FullClass, kept.FullClass, excludePaths)
			if err != nil {
				fmt.Fprintln(os.Stderr, "替换 import 失败:", err)
				os.Exit(1)
			}
			allModified = append(allModified, modified...)
		}
		for _, jf := range toDelete {
			// 暂时注释，先不删文件
			fmt.Println("-- 删除： ", jf.Path)
			if err := os.Remove(jf.Path); err != nil {
				fmt.Fprintln(os.Stderr, "删除失败:", jf.Path, err)
				os.Exit(1)
			}
			allDeleted = append(allDeleted, jf.Path)
		}
		fmt.Printf("- %s: 保留 %s，已删除 %d 个重复文件，并已将引用改为保留类\n",
			name, ctool.Green.Print(kept.FullClass), len(toDelete))
	}

	if len(allDeleted) == 0 {
		fmt.Println("没有类定义完全一致的重名类，未执行删除。")
		return
	}
	fmt.Printf("\n共删除 %d 个重复类文件，修改 %d 个引用文件。\n", len(allDeleted), len(allModified))
	if len(allModified) > 0 {
		fmt.Println("被修改 import 的文件：")
		seen := make(map[string]struct{})
		for _, p := range allModified {
			if _, ok := seen[p]; ok {
				continue
			}
			seen[p] = struct{}{}
			rel := strings.Replace(p, root, "", 1)
			if rel == p {
				rel = p
			}
			fmt.Printf("    %s\n", ctool.Yellow.Print(rel))
		}
	}
}

func main() {
	info.Parse()
	if help {
		info.PrintHelp()
		return
	}

	// 如果projectRoot没有值，默认取当前目录的完整路径
	if projectRoot == "" {
		projectRoot, _ = os.Getwd()
	}
	if duplicateClass {
		findDuplicateClass()
		return
	}
	if duplicateDiff {
		findDuplicateClassWithDiff()
		return
	}
	if rmDuplicateClass {
		deleteDuplicateClass()
		return
	}
	if noDep {
		findNoDepClass()
		return
	}
	if lessDep {
		findLessDepClass()
		return
	}

	findNoDepClass()
}
