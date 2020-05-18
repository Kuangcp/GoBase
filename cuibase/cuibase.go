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

// ParamInfo one line struct
type (
	ParamInfo struct {
		Verb    string
		Param   string
		Comment string
		Handler func(params []string)
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

// RunActionFromInfo defaultAction 为空时默认PrintHelp
func RunActionFromInfo(info HelpInfo, defaultAction func(params []string)) {
	if len(info.Params) == 0 {
		return
	}
	params := os.Args
	if len(params) < 2 {
		if defaultAction != nil {
			defaultAction(params)
		} else {
			info.PrintHelp()
		}
		return
	}

	verb := params[1]
	for _, param := range info.Params {
		if verb != param.Verb {
			continue
		}

		if param.Handler != nil {
			param.Handler(params)
			return
		} else {
			info.PrintHelp()
		}
	}
}

func RunAction(actions map[string]func(params []string), defaultAction func(params []string)) {
	runAction(os.Args, actions, defaultAction)
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
	if action != nil {
		action(os.Args)
		return
	}
	if defaultAction != nil {
		defaultAction(os.Args)
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
