package cuibase

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
)

// Color
var Red = "\033[0;31m"
var Green = "\033[0;32m"
var Yellow = "\033[0;33m"
var Blue = "\033[0;34m"
var Purple = "\033[0;35m"
var Cyan = "\033[0;36m"
var White = "\033[0;37m"
var End = "\033[0m"

var LightRed = "\033[0;91m"
var LightGreen = "\033[0;92m"
var LightYellow = "\033[0;93m"
var LightBlue = "\033[0;94m"
var LightPurple = "\033[0;95m"
var LightCyan = "\033[0;96m"
var LightWhite = "\033[0;97m"


// ParamInfo one line struct
type (
	ParamInfo struct {
		Verb    string
		Param   string
		Comment string
	}

	HelpInfo struct {
		Version     string
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

func PrintTitle(command string, description string) {
	fmt.Printf("%sUsage:%s\n\n  %v %v <verb> %v <param> %v\n\n", LightGreen, End, command, Green, Yellow, End)
	fmt.Printf("%sDescription:%s\n\n  %v\n\n", LightGreen, End, description)
}

func RunAction(actions map[string]func(params []string), defaultAction func(params []string)) {
	runAction(os.Args, actions, defaultAction)
}

func Help(helpInfo HelpInfo) {
	printTitleDefault(helpInfo.Description)
	format := BuildFormat(helpInfo.VerbLen, helpInfo.ParamLen)
	PrintParams(format, helpInfo.Params)
	if helpInfo.Version != "" {
		fmt.Printf("\n%sVersion:%s  %v\n\n", LightGreen, End, helpInfo.Version)
	}
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func ReadFileLines(filename string, filterFunc func(string) bool, mapFunc func(string) interface{}) []interface{} {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		log.Println("Open file error!", err)
		return nil
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if stat.Size() == 0 {
		log.Printf("file:%s is empty", filename)
		return nil
	}

	var result []interface{}

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
				log.Println("Read file error!", err)
				return nil
			}
		}
	}
	return result
}

func Home() (string, error) {
	curUser, err := user.Current()
	if nil == err {
		return curUser.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func PrintWithColorful() {
	for i := 0; i < 255; i++ {
		fmt.Printf("\x1b[48;5;%dm%3d\u001B[0m", i, i)
		if i == 15 || (i > 15 && ((i-15)%6 == 0)) {
			println()
		}
	}
}

func printTitleDefault(description string) {
	PrintTitle(os.Args[0], description)
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

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

func enoughCount(params []string, count int) bool {
	return len(params) > count
}
