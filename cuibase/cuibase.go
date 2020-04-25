package cuibase

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Color
var Red = "\x1b[0;31m"
var Green = "\x1b[0;32m"
var Yellow = "\x1b[0;33m"
var Blue = "\x1b[0;34m"
var Purple = "\x1b[0;35m"
var Cyan = "\x1b[0;36m"
var White = "\x1b[0;37m"
var End = "\x1b[0m"

// ParamInfo one line struct
type (
	ParamInfo struct {
		Verb    string
		Param   string
		Comment string
	}

	HelpInfo struct {
		Description string
		VerbLen     int
		ParamLen    int
		Params      []ParamInfo
	}
)

// AssertParamCount os.Args 参数构成: 0 go源文件; 1 参数1; 2 参数2; count 必填参数个数
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

func PrintTitleDefault(description string) {
	PrintTitle(os.Args[0], description)
}

func PrintTitle(command string, description string) {
	fmt.Printf("  usage: %v %v <verb> %v <param> %v\n\n", command, Green, Yellow, End)
	fmt.Printf("  %v\n\n", description)
}

func RunAction(actions map[string]func(params []string), defaultAction func(params []string)) {
	runAction(os.Args, actions, defaultAction)
}

func Help(helpInfo HelpInfo) {
	PrintTitleDefault(helpInfo.Description)
	format := BuildFormat(helpInfo.VerbLen, helpInfo.ParamLen)
	PrintParams(format, helpInfo.Params)
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func runAction(params []string, actions map[string]func(params []string), defaultAction func(params []string)) {
	if len(params) < 2 {
		defaultAction(os.Args)
		return
	}

	verb := params[1]
	action := actions[verb]
	if action == nil {
		defaultAction(os.Args)
	} else {
		action(os.Args)
	}
}

func enoughCount(params []string, count int) bool {
	return len(params) > count
}
