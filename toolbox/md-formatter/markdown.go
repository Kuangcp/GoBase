package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kuangcp/gobase/cuibase"
	"github.com/wonderivan/logger"
)

type filterFun = func(string) bool
type mapFun func(string) string

func help(_ []string) {
	info := cuibase.HelpInfo{
		Description: "Start simple http server on current path",
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
				Verb:    "-a",
				Param:   "file",
				Comment: "append catalog",
			},
			{
				Verb:    "-at",
				Param:   "file",
				Comment: "append title and catalog",
			},
			{
				Verb:    "-mm",
				Param:   "file",
				Comment: "show mind map",
			},
		}}
	cuibase.Help(info)
}

func readLines(filename string, filterFunc filterFun, mapFunc mapFun) []string {
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
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

func normalizeForTitle(title string) string {
	title = strings.Replace(title, " ", "-", -1)
	title = strings.ToLower(title)
	return title
}

func printCategory(filename string) {
	lines := readLines(filename, func(s string) bool {
		return strings.HasPrefix(s, "#")
	}, func(s string) string {
		title := strings.TrimSpace(strings.Replace(s, "#", "", -1))
		strings.Count(s, "#")
		temps := strings.Split(s, "# ")
		level := strings.Replace(temps[0], "#", "    ", -1)
		return fmt.Sprintf("%s1. [%s](#%s)\n", level, title, normalizeForTitle(title))
	})

	if lines != nil {
		for i := range lines {
			fmt.Print(lines[i])
		}
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

func main() {
	cuibase.RunAction(map[string]func(params []string){
		"-h": help,
		"-mm": func(params []string) {
			cuibase.AssertParamCount(2, "must input filename ")
			printMindMap(params[2])
		},
		"-f": func(params []string) {
			cuibase.AssertParamCount(2, "must input filename ")
			printCategory(params[2])
		},
	}, help)
}
