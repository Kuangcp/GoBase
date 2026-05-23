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

	"github.com/atotto/clipboard"
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
		"out", "build", "target", "log", "logs", "__pycache__",
	}
	ignoreFiles = [...]string{
		"README", "Readme", "Readme_CN", "readme", "SUMMARY", "Process", "License", "LICENSE",
	}
	handleSuffix = [...]string{
		".md", ".markdown", ".txt",
	}
	// 制作标题跳转时需要去除的符号
	titleRemoveChar = []string{
		".", "【", "】", ":", "：", ",", "，", "/", "(", ")", "（", "）", "《", "》", "*", "＊", "。", "?", "？", "+", "＋",
	}
)

var (
	splitTag    = "💠"
	headerFirst = "---\n"
	headerLast  = "****************************************\n"
	codeBlock   = "```"
)

var tagTemplate = headerFirst + `title: %s
date: %s
tags: 
categories: 
---
`

var (
	help             bool
	refreshDir       string
	mindMapFile      string
	refreshChangeDir string
	appendFile       string
	printCatalog     string
	rmFile           string
	rmAppendFile     string
	extractTitleUrl  string

	titleReplace *strings.Replacer
)

func main() {
	flag.Parse()
	if help {
		info.PrintHelp()
		return
	}

	if extractTitleUrl != "" {
		CopyWebsiteTitle()
		return
	}

	if len(os.Args) < 2 {
		logger.Fatal("Usage: md-formatter <command> [<args>]")
	}
	prepareContext()

	// action
	for _, a := range acts {
		a.tryInvoke()
	}

	filename := os.Args[1]
	RefreshTagAndCatalog(filename)
}

func CopyWebsiteTitle() {
	if extractTitleUrl == "auto" {
		//"github.com/go-vgo/robotgo"
		//robotgo.KeyTap("y")
		//robotgo.KeyTap("y")
		//time.Sleep(time.Second)

		last, err := clipboard.ReadAll()
		if err != nil {
			logger.Error(err)
			return
		}
		last = strings.TrimSpace(last)
		if last == "" || !strings.Contains(last, "http") {
			return
		}
		extractTitleUrl = last
	}

	//fmt.Println(extractTitleUrl)

	// 库主要是为了获取站点主页面的信息，所以要禁用重定向才能获取当前页信息
	s, err := Scrape(extractTitleUrl, 0)
	if err != nil {
		logger.Error(err)
		return
	}
	title := s.Preview.Title
	if title != "" {
		//fmt.Printf("Title : %s\n", title)
		clipboard.WriteAll("> [" + title + "](" + extractTitleUrl + ")  ")
	}
}

func prepareContext() {
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
func RefreshDirAllFiles(path string) {
	var fileList = list.New()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("occur error: ", err)
			return err
		}

		if info.IsDir() {
			last := strings.Split(path, "/")
			if strings.HasPrefix(last[len(last)-1], ".") {
				return filepath.SkipDir
			}
			if strings.HasPrefix(path, ".") {
				return filepath.SkipDir
			}
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

func PrintCatalog(filename string) {
	article := BuildArticle(filename)
	if article == nil {
		logger.Error(filename + " 格式有误，未包含定位行： " + headerLast)
		return
	}
	article.Refresh()

	for _, r := range article.catalog {
		fmt.Print(r)
	}
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

	RefreshTagAndCatalog(filename)
}

func RefreshTagAndCatalog(filename string) {
	if refreshChangeDir != "" {
		logger.Info("refresh:", strings.TrimLeft(filename, refreshChangeDir))
	} else {
		logger.Info("refresh:", filename)
	}

	article := BuildArticle(filename)
	if article == nil {
		logger.Error(filename + " 格式有误，未包含定位行： " + headerLast)
		return
	}
	article.Refresh()
	article.writeToDisk(false)
	//logger.Info("\n" + strings.Join(article.tag, ""))
	//logger.Info("\n" + strings.Join(article.catalog, ""))
}

// RefreshChangeFile 更新指定目录的Git仓库中 发生变更 的文件
func RefreshChangeFile(dir string) {
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
		if strings.HasPrefix(filePath, ".") {
			continue
		}
		fileStatus := status.File(filePath)

		careStatus := fileStatus.Staging == git.Added || fileStatus.Staging == git.Modified ||
			fileStatus.Worktree == git.Modified || fileStatus.Worktree == git.Untracked
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

func RemoveCatalog(filename string) {
	article := BuildArticle(filename)
	if article != nil {
		article.writeToDisk(true)
	}
}
