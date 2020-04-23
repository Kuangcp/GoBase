package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kuangcp/gobase/cuibase"
	"github.com/wonderivan/logger"
)

type filterFun = func(string) bool
type mapFun func(string) string

var startTag = "**目录 start**"
var endTag = "**目录 end**"

var ignoreDirs = [...]string{
	".git", ".svn", ".vscode", ".idea", ".gradle",
	"out", "build", "target", "log", "logs", "__pycache__", "ARTS",
}
var ignoreFiles = [...]string{
	"README.md", "Readme.md", "readme.md", "SUMMARY.md", "Process.md", "License.md",
}
var handleSuffix = [...]string{
	".md", ".markdown", ".txt",
}
var deleteChar = [...]string{
	".", "【", "】", ":", "：", ",", "，", "/", "(", ")", "《", "》", "*", "。", "?", "？",
}

func help(_ []string) {
	info := cuibase.HelpInfo{
		Description: "Format markdown file, generate catalog",
		VerbLen:     -5,
		ParamLen:    -5,
		Params: []cuibase.ParamInfo{
			{
				Verb:    "-h",
				Param:   "",
				Comment: "help",
			},
			{
				Verb:    "",
				Param:   "file",
				Comment: "refresh catalog",
			},
			{
				Verb:    "-f",
				Param:   "file",
				Comment: "refresh catalog",
			},
			{
				Verb:    "-d",
				Param:   "dir",
				Comment: "refresh catalog with dir",
			},
			{
				Verb:    "-mm",
				Param:   "file",
				Comment: "show mind map",
			},
		}}
	cuibase.Help(info)
}

func readFileLines(filename string) []string {
	return readLines(filename, func(s string) bool {
		return true
	}, func(s string) string {
		return s
	})
}

func readLines(filename string, filterFunc filterFun, mapFunc mapFun) []string {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		logger.Error("Open file error!", err)
		return nil
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if stat.Size() == 0 {
		logger.Info("file:%s is empty", filename)
		return nil
	}

	var result []string

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if filterFunc(line) {
			result = append(result, mapFunc(line))
		}

		if err != nil {
			if err == io.EOF {
				break
			} else {
				logger.Error("Read file error!", err)
				return nil
			}
		}
	}
	return result
}

func isFileNeedHandle(filename string) bool {
	for _, file := range ignoreFiles {
		if strings.HasSuffix(filename, file) {
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

func refreshDirAllFiles(path string) {
	var fileList = list.New()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("occur error: ", err)
			return nil
		}

		if info.IsDir() {
			for _, dir := range ignoreDirs {
				if path == dir {
					return filepath.SkipDir
				}
			}
			return nil
		}
		fileList.PushBack(path)
		return nil
	})
	if err != nil {
		logger.Error(err)
	}

	for e := fileList.Front(); e != nil; e = e.Next() {
		fileName := e.Value.(string)
		if isFileNeedHandle(fileName) {
			logger.Info(fileName)
			refreshCatalog(fileName)
		}
	}
}

func normalizeForTitle(title string) string {
	title = strings.Replace(title, " ", "-", -1)
	title = strings.ToLower(title)

	for _, char := range deleteChar {
		title = strings.Replace(title, char, "", -1)
	}

	return title
}

func generateCatalog(filename string) []string {
	return readLines(filename, func(s string) bool {
		return strings.HasPrefix(s, "#")
	}, func(s string) string {
		title := strings.TrimSpace(strings.Replace(s, "#", "", -1))
		strings.Count(s, "#")
		temps := strings.Split(s, "# ")
		level := strings.Replace(temps[0], "#", "    ", -1)
		return fmt.Sprintf("%s1. [%s](#%s)\n", level, title, normalizeForTitle(title))
	})
}

func refreshCatalog(filename string) {
	titles := generateCatalog(filename)
	lines := readFileLines(filename)

	startIdx := -1
	endIdx := -1
	var result = ""
	for i, line := range lines {
		if strings.Contains(line, startTag) {
			startIdx = i
		}
		if strings.Contains(line, endTag) {
			endIdx = i
			result += startTag + "\n\n"
			for t := range titles {
				result += titles[t]
			}
			result += "\n"
		}
		if startIdx == -1 || (startIdx != -1 && endIdx != -1) {
			result += line
		}
	}

	if startIdx == -1 || endIdx == -1 {
		logger.Warn("Invalid catalog: ", filename, startIdx, endIdx)
		return
	}
	//logger.Info("index", startIdx, endIdx, result)
	if ioutil.WriteFile(filename, []byte(result), 0644) != nil {
		logger.Error("write error")
	}
}

func printMindMap(filename string) {
	cuibase.AssertParamCount(2, "must input filename ")

	lines := readLines(filename, func(s string) bool {
		return strings.HasPrefix(s, "#")
	}, func(s string) string {
		temp := strings.Split(s, "# ")
		prefix := strings.Replace(temp[0], "#", "    ", -1)
		return prefix + temp[1]
	})

	if lines != nil {
		for i := range lines {
			fmt.Print(lines[i])
		}
	}
}

func refreshCatalogFromParam(params []string) {
	cuibase.AssertParamCount(2, "must input filename ")
	refreshCatalog(params[2])
}

func main() {
	logger.SetLogPathTrim("/toolbox/")
	cuibase.RunAction(map[string]func(params []string){
		"-h": help,
		"-mm": func(params []string) {
			cuibase.AssertParamCount(2, "must input filename ")
			printMindMap(params[2])
		},
		"-f": func(params []string) {
			cuibase.AssertParamCount(2, "must input filename ")
			refreshCatalog(params[2])
		},
		"-d": func(params []string) {
			refreshDirAllFiles("./")
		},
	}, func(params []string) {
		cuibase.AssertParamCount(1, "must input filename ")
		refreshCatalog(params[1])
	})
}
