package cuibase

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var Red = "\033[0;31m"
var Green = "\033[0;32m"
var Yellow = "\033[0;33m"
var Blue = "\033[0;34m"
var Purple = "\033[0;35m"
var Cyan = "\033[0;36m"
var White = "\033[0;37m"
var End = "\033[0m"

var Format = ""

// ParamInfo one line struct
type ParamInfo struct {
	Verb    string
	Param   string
	Comment string
}

// AssertParamCount os.Args 参数构成: 0 文件 1 参数 2 参数
func AssertParamCount(count int, msg string) {
	param := os.Args
	flag := enoughCount(param, count)
	if !flag {
		log.Printf("param count less than %v \n", count)
		log.Fatal(msg)
	}
}

func BuildFormat(verbLen int, paramLen int) string {
	return "    %v %" + strconv.Itoa(verbLen) + "v %v %" + strconv.Itoa(paramLen) + "v %v %v\n"
}

func PrintParam(format string, verb string, param string, comment string) {
	fmt.Printf(format, Green, verb, Yellow, param, End, comment)
}

func PrintParams(format string, params []ParamInfo) {
	for _, param := range params {
		PrintParam(format, param.Verb, param.Param, param.Comment)
	}
}

func PrintTitle(command string, description string) {
	fmt.Printf("  usage: %v %v <verb> %v <param> %v\n\n", command, Green, Yellow, End)
	fmt.Printf("  %v\n\n", description)
}

func enoughCount(param []string, count int) bool {
	return len(param) > count
}
