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
	"time"

	"github.com/go-git/go-git/v5"
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
	"README.md", "Readme.md", "Readme_CN.md", "readme.md", "SUMMARY.md", "Process.md", "License.md",
}
var handleSuffix = [...]string{
	".md", ".markdown", ".txt",
}
var deleteChar = [...]string{
	".", "【", "】", ":", "：", ",", "，", "/", "(", ")", "《", "》", "*", "＊", "。", "?", "？",
}

var info = cuibase.HelpInfo{
	Description: "Format markdown file, generate catalog",
	Version:     "1.0.1",
	VerbLen:     -3,
	ParamLen:    -5,
	Params: []cuibase.ParamInfo{
		{
			Verb:    "-h",
			Param:   "",
			Comment: "Help info",
		}, {
			Verb:    "",
			Param:   "file",
			Comment: "Refresh catalog for file",
			Handler: func(params []string) {
				cuibase.AssertParamCount(1, "must input filename ")
				refreshCatalog(params[1])
			},
		}, {
			Verb:    "-f",
			Param:   "file",
			Comment: "Refresh catalog for file",
			Handler: func(params []string) {
				cuibase.AssertParamCount(2, "must input filename ")
				refreshCatalog(params[2])
			},
		}, {
			Verb:    "-d",
			Param:   "dir",
			Comment: "Refresh catalog for file that recursive dir",
			Handler: func(_ []string) {
				refreshDirAllFiles("./")
			},
		}, {
			Verb:    "-mm",
			Param:   "file",
			Comment: "Print mind map",
			Handler: func(params []string) {
				cuibase.AssertParamCount(2, "must input filename ")
				printMindMap(params[2])
			},
		}, {
			Verb:    "-rc",
			Param:   "dir",
			Comment: "Refresh git repo dir changed file",
			Handler: func(params []string) {
				cuibase.AssertParamCount(2, "must input repo dir ")
				refreshChangeFile(params[2])
			},
		}, {
			Verb:    "-a",
			Param:   "file",
			Comment: "Append catalog and title for file",
			Handler: func(params []string) {
				cuibase.AssertParamCount(2, "must input filename")
				appendCatalogAndTitle(params[2])
			},
		},
	}}

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

// 递归更新当前目录下所有文件的目录
func refreshDirAllFiles(path string) {
	var fileList = list.New()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("occur error: ", err)
			return err
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
		return
	}

	for e := fileList.Front(); e != nil; e = e.Next() {
		fileName := e.Value.(string)
		if isNeedHandleFile(fileName) {
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
	return readLinesWithFunc(filename,
		func(s string) bool {
			return strings.HasPrefix(s, "#")
		},
		func(s string) string {
			title := strings.TrimSpace(strings.Replace(s, "#", "", -1))
			strings.Count(s, "#")
			temps := strings.Split(s, "# ")
			level := strings.Replace(temps[0], "#", "    ", -1)
			return fmt.Sprintf("%s1. [%s](#%s)\n", level, title, normalizeForTitle(title))
		})
}

// 更新指定文件的目录
func refreshCatalog(filename string) {
	logger.Info("refresh:", filename)

	titleBlock := ""
	titles := generateCatalog(filename)
	for t := range titles {
		titleBlock += titles[t]
	}

	startIdx := -1
	endIdx := -1
	var result = ""

	// replace title block
	lines := readFileLines(filename)
	for i, line := range lines {
		if strings.Contains(line, startTag) {
			startIdx = i
		}
		if strings.Contains(line, endTag) {
			endIdx = i
			timeStr := time.Now().Format("2006-01-02 15:04")
			result += startTag + "\n\n" + titleBlock + "\n" + endTag + "|_" + timeStr + "_|\n"
			continue
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
		logger.Error("write error", filename)
	}
}

// 打印 百度脑图支持的 MindMap 格式
func printMindMap(filename string) {
	cuibase.AssertParamCount(2, "must input filename ")

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

// 更新指定目录的Git仓库里发成变更的文件
func refreshChangeFile(dir string) {
	r, err := git.PlainOpen(dir)
	cuibase.CheckIfError(err)
	worktree, err := r.Worktree()
	cuibase.CheckIfError(err)
	status, err := worktree.Status()
	cuibase.CheckIfError(err)
	if status.IsClean() {
		return
	}

	for filePath := range status {
		fileStatus := status.File(filePath)
		if fileStatus.Staging == git.Modified || fileStatus.Worktree == git.Modified {
			if isNeedHandleFile(dir + filePath) {
				refreshCatalog(dir + filePath)
			}
		}
	}
}

func appendCatalogAndTitle(filename string) {
	lines := readLinesWithFunc(filename, nil, nil)
	var result = "---\ntitle: " + filename + "\ndate: " +
		time.Now().Format("2006-01-02 15:04:05") +
		"\ntags: \ncategories: \n---\n\n" + startTag + "\n" + endTag +
		"\n****************************************\n"
	for i := range lines {
		result += lines[i]
	}

	if ioutil.WriteFile(filename, []byte(result), 0644) != nil {
		logger.Error("write error", filename)
	}
	refreshCatalog(filename)
}

func main() {
	logger.SetLogPathTrim("/toolbox/")
	cuibase.RunActionFromInfo(info, nil)
}
