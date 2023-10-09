package main

import (
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/kuangcp/gobase/pkg/ctk"
	"github.com/kuangcp/logger"
)

type filterFun = func(string) bool
type mapFun func(string) string

var (
	buildVersion string
	ignoreDirMap = make(map[string]int8)
	ignoreDirs   = [...]string{
		".git", ".svn", ".vscode", ".idea", ".gradle", "out", "build", "target", "log", "logs", "__pycache__",
	}
	ignoreFiles = [...]string{
		"README", "Readme", "Readme_CN", "readme", "SUMMARY", "Process", "License", "LICENSE",
	}
	handleSuffix = [...]string{
		".md", ".markdown", ".txt",
	}
	titleRemoveChar = []string{
		".", "【", "】", ":", "：", ",", "，", "/", "(", ")", "（", "）", "《", "》", "*", "＊", "。", "?", "？",
	}
)

var (
	splitTag       = "💠"
	headerFirst    = "---\n"
	headerLast     = "****************************************\n"
	headerTemplate = headerFirst + `title: %s
date: %s
tags: 
categories: 
---

%s
%s
` + headerLast
	catalogTemplate = `
%s
%s
` + headerLast
	headerFmt = splitTag + "\n\n%s\n" + splitTag + " %s\n"
)

var (
	help             bool
	refreshDir       string
	mindMapFile      string
	refreshChangeDir string
	appendFile       string
	printCatalog     string
	rmFile           string
	rmAppendFile     string

	titleReplace *strings.Replacer
)

var info = ctk.HelpInfo{
	Description:   "Format markdown file, generate catalog",
	Version:       "1.0.4",
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
		{Short: "-c", Value: "dir", Comment: "Refresh git repo dir changed file"},
		{Short: "-a", Value: "file", Comment: "Append catalog and title for file"},
		{Short: "-ra", Value: "file", Comment: "Remove then Append catalog and title for file"},
	},
}

func init() {
	flag.BoolVar(&help, "h", false, "")
	flag.StringVar(&refreshDir, "d", "", "")
	flag.StringVar(&mindMapFile, "mm", "", "")
	flag.StringVar(&refreshChangeDir, "c", "", "")
	flag.StringVar(&appendFile, "a", "", "")
	flag.StringVar(&rmFile, "r", "", "")
	flag.StringVar(&rmAppendFile, "ra", "", "")
	flag.StringVar(&printCatalog, "p", "", "")

	//logger.SetLogPathTrim("md-formatter/")
	flag.Usage = info.PrintHelp

	var replacePairList []string
	for i := range titleRemoveChar {
		replacePairList = append(replacePairList, titleRemoveChar[i], "")
	}
	replacePairList = append(replacePairList, " ", "-")
	titleReplace = strings.NewReplacer(replacePairList...)

	for _, dir := range ignoreDirs {
		ignoreDirMap[dir] = 1
	}
}

func main() {
	flag.Parse()
	if help {
		info.PrintHelp()
		return
	}

	invokeWhen(refreshDir, refreshDirAllFiles)
	invokeWhen(mindMapFile, printMindMap)
	invokeWhen(refreshChangeDir, refreshChangeFile)
	invokeWhen(appendFile, createCatalog)
	invokeWhen(rmFile, func(f string) {
		removeCatalog(f)
	})
	invokeWhen(rmAppendFile, ReplaceThenRefreshCatalog)
	invokeWhen(printCatalog, func(s string) {
		rows := generateCatalog(s)
		for _, r := range rows {
			fmt.Print(r)
		}
	})

	ReplaceThenRefreshCatalog(os.Args[1])
}

func invokeWhen(param string, action func(string)) {
	if param != "" {
		action(param)
		os.Exit(0)
	}
}

func readFileLines(filename string) []string {
	return readLinesWithFunc(filename,
		func(s string) bool {
			return true
		},
		func(s string) string {
			return s
		})
}

func readLinesWithFunc(filename string, filterFunc filterFun, mapFunc mapFun) []string {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if stat.Size() == 0 {
		logger.Warn("file:%s is empty", filename)
		return nil
	}

	var result []string
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if filterFunc == nil || filterFunc(line) {
			if mapFunc != nil {
				result = append(result, mapFunc(line))
			} else {
				result = append(result, line)
			}
		}

		if err == nil {
			continue
		}
		if err == io.EOF {
			break
		}

		logger.Error("Read file error!", err)
		return nil
	}
	return result
}

func isNeedHandleFile(filename string) bool {
	for _, file := range ignoreFiles {
		if strings.Contains(filename, file) {
			return false
		}
	}
	for _, fileType := range handleSuffix {
		if strings.HasSuffix(filename, fileType) {
			return true
		}
	}
	return false
}

// 递归更新当前目录下所有文件的目录
func refreshDirAllFiles(path string) {
	var fileList = list.New()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("occur error: ", err)
			return err
		}

		if info.IsDir() {
			_, ok := ignoreDirMap[path]
			if ok {
				return filepath.SkipDir
			}
			return nil
		}
		fileList.PushBack(path)
		return nil
	})
	if err != nil {
		logger.Error(err)
		return
	}

	for e := fileList.Front(); e != nil; e = e.Next() {
		refreshCatalogWithCondition(e.Value.(string), isNeedHandleFile)
	}
}

func normalizeForTitle(title string) string {
	title = strings.ToLower(title)
	return titleReplace.Replace(title)
}

func generateCatalog(filename string) []string {

	var pPath []int

	rows := readLinesWithFunc(filename,
		func(s string) bool {
			return strings.HasPrefix(s, "#")
		},
		func(s string) string {
			if len(pPath) == 0 {
				pPath = []int{0}
			}
			level := strings.Count(s, "#")
			for len(pPath) < level {
				pPath = append(pPath, 0)
			}
			pPath[level-1] += 1
			if level < len(pPath) {
				for i := level; i < len(pPath); i++ {
					pPath[i] = 0
				}
			}

			title := strings.TrimSpace(strings.Replace(s, "#", "", -1))
			strings.Count(s, "#")
			temps := strings.Split(s, "# ")
			levelStr := strings.Replace(temps[0], "#", "    ", -1)
			return fmt.Sprintf("%s- %s. [%s](#%s)\n", levelStr, pathToString(pPath[:level]), title, normalizeForTitle(title))
		})
	return rows
}
func pathToString(path []int) string {
	var result []string
	for _, i := range path {
		result = append(result, fmt.Sprint(i))
	}
	return strings.Join(result, ".")
}

func refreshCatalogWithCondition(filename string, condition func(filename string) bool) {
	if !condition(filename) {
		return
	}

	//refreshCatalog(filename)
	ReplaceThenRefreshCatalog(filename)
}

// 删除原有目录数据，保留标签，更新索引目录
func ReplaceThenRefreshCatalog(filename string) {
	replaceCatalog(filename, removeCatalog(filename))
}

// 更新指定文件的目录
func refreshCatalog(filename string) {
	if refreshChangeDir != "" {
		logger.Info("refresh:", strings.TrimLeft(filename, refreshChangeDir))
	} else {
		logger.Info("refresh:", filename)
	}

	tocBlock := ""
	tocList := generateCatalog(filename)
	if tocList == nil {
		return
	}
	for t := range tocList {
		tocBlock += tocList[t]
	}

	startIdx := -1
	endIdx := -1
	var content = ""

	// replace title block
	lines := readFileLines(filename)
	if lines == nil {
		return
	}
	for i, line := range lines {
		if startIdx != -1 && strings.Contains(line, splitTag) {
			endIdx = i
			content += fmt.Sprintf(headerFmt, tocBlock, time.Now().Format("2006-01-02 15:04"))
			continue
		}

		if startIdx == -1 && strings.Contains(line, splitTag) {
			startIdx = i
		}

		if startIdx == -1 || (startIdx != -1 && endIdx != -1) {
			content += line
		}
	}

	if startIdx == -1 || endIdx == -1 {
		logger.Warn("Invalid catalog: ", filename, startIdx, endIdx)
		return
	}
	//logger.Info("index", startIdx, endIdx, content)
	if os.WriteFile(filename, []byte(content), 0644) != nil {
		logger.Error("Write error", filename)
	}
}

// 打印 百度脑图支持的 MindMap 格式
func printMindMap(filename string) {
	if filename == "" {
		return
	}
	lines := readLinesWithFunc(filename,
		func(s string) bool {
			return strings.HasPrefix(s, "#")
		},
		func(s string) string {
			temp := strings.Split(s, "# ")
			prefix := strings.Replace(temp[0], "#", "    ", -1)
			return prefix + temp[1]
		})

	if lines == nil {
		return
	}
	for i := range lines {
		fmt.Print(lines[i])
	}
}

// 更新指定目录的Git仓库中 发生变更 的文件
func refreshChangeFile(dir string) {
	r, err := git.PlainOpen(dir)
	ctk.CheckIfError(err)
	worktree, err := r.Worktree()
	ctk.CheckIfError(err)
	status, err := worktree.Status()
	ctk.CheckIfError(err)
	if status.IsClean() {
		return
	}

	showChange := false
	for filePath := range status {
		fileStatus := status.File(filePath)

		careStatus := fileStatus.Staging == git.Modified || fileStatus.Worktree == git.Modified
		if careStatus && !showChange {
			logger.Info("Repository:", refreshChangeDir)
			showChange = true
		}
		if careStatus {
			refreshCatalogWithCondition(dir+filePath, isNeedHandleFile)
		}
	}
}

func buildTitle(filename string) string {
	if strings.HasSuffix(filename, ".md") {
		filename = filename[:len(filename)-3]
	}
	index := strings.LastIndex(filename, "/")
	if index != -1 && index < len(filename)-1 {
		filename = filename[index+1:]
	}

	return filename
}

func replaceCatalog(filename, tag string) {
	lines := readLinesWithFunc(filename, nil, nil)
	if len(lines) > 0 && lines[0] != headerFirst {
		var content = fmt.Sprintf(catalogTemplate, splitTag, splitTag)
		content = tag + content
		for _, line := range lines {
			content += line
		}

		if os.WriteFile(filename, []byte(content), 0644) != nil {
			logger.Error("write error", filename)
		}
	}

	refreshCatalog(filename)
}

func createCatalog(filename string) {
	lines := readLinesWithFunc(filename, nil, nil)
	if len(lines) > 0 && lines[0] != headerFirst {
		var content = fmt.Sprintf(headerTemplate, buildTitle(filename),
			time.Now().Format("2006-01-02 15:04:05"), splitTag, splitTag)
		for _, line := range lines {
			content += line
		}

		if os.WriteFile(filename, []byte(content), 0644) != nil {
			logger.Error("write error", filename)
		}
	}
	refreshCatalog(filename)
}

func removeCatalog(filename string) string {
	tagBlock := ""
	lines := readLinesWithFunc(filename, nil, nil)
	content := ""
	header := false
	tagEnd := false
	for _, line := range lines {
		line = strings.Replace(line, "\r\n", "\n", 1)
		if line == headerLast {
			header = false
			continue
		}
		if line == headerFirst && header {
			tagEnd = true
			tagBlock += line
			continue
		}
		if line == headerFirst {
			header = true
		}
		if header {
			if !tagEnd {
				tagBlock += line
			}
			continue
		}

		content += line
	}

	if os.WriteFile(filename, []byte(content), 0644) != nil {
		logger.Error("write error", filename)
	}
	return tagBlock
}
