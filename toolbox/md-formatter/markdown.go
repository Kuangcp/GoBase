package main

import (
	"bufio"
	"fmt"
	"github.com/kuangcp/gobase/cuibase"
	"github.com/wonderivan/logger"
	"io"
	"os"
	"strings"
)

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

func readLines(filename string, filterFunc func(string) bool, mapFunc func(string) string) []string {
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
		logger.Debug("file:%s is empty", filename)
		return nil
	}

	var result []string

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
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

func printMindMapFormat(params []string) {
	cuibase.AssertParamCount(2, "must input filename ")

	lines := readLines(params[2], func(s string) bool {
		return strings.HasPrefix(s, "#")
	}, func(s string) string {
		temp := strings.Split(s, " ")
		prefix := strings.Replace(temp[0], "#", "    ", -1)
		return prefix + temp[1]

	})

	if lines != nil {
		for i := range lines {
			fmt.Println(lines[i])
		}
	}
}

func main() {
	cuibase.RunAction(map[string]func(params []string){
		"-h":  help,
		"-mm": printMindMapFormat,
	}, help)
}
