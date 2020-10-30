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

// AssertParamCount os.Args 参数构成: 0 go源文件; 1 参数1; 2 参数2; count 必填参数个数
func AssertParamCount(count int, msg string) {
	param := os.Args
	flag := enoughCount(param, count)
	if !flag {
		log.Printf("param count less than %v \n", count)
		log.Fatal(msg)
	}
}

// BuildFormat
func BuildFormat(info HelpInfo) string {
	single := strconv.Itoa(info.SingleFlagLen)
	double := strconv.Itoa(info.DoubleFlagLen)
	value := strconv.Itoa(info.ValueLen)
	return "    %v %" + single + "v%" + double + "v %v %" + value + "v %v %v\n"
}

// PrintParams
func PrintParams(format string, flagColor Color, params []ParamVO) {
	for _, vo := range params {
		if vo.Long == "" {
			fmt.Printf(format, flagColor, vo.Short, "", Yellow, vo.Value, End, vo.Comment)
		} else {
			fmt.Printf(format, flagColor, vo.Short, ", "+vo.Long, Yellow, vo.Value, End, vo.Comment)
		}
	}
}

// PrintTitle
func PrintTitle(command string, helpInfo HelpInfo) {
	flagStr := ""
	for _, flagVO := range helpInfo.Flags {
		flagStr += flagVO.Short
	}
	flagStr = strings.Replace(flagStr, "-", "", -1)

	optionStr := ""
	for _, option := range helpInfo.Options {
		optionStr += fmt.Sprintf("[%s %s] ", option.Short, option.Value)
	}
	fmt.Printf("%s\n\n  %v %v %v\n\n",
		LightCyan.Print("Usage:"),
		command,
		Yellow.PrintNoEnd("[-"+flagStr+"]"),
		Purple.Print(optionStr))

	fmt.Printf("%s\n\n  %v\n", LightCyan.Print("Description:"), helpInfo.Description)
}

// CheckIfError assert err is nil
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// ReadFileLines
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

// Home 
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

func OpenBrowser(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	return exec.Command(cmd, append(args, url)...).Start()
}

// PrintWithColorful
func PrintWithColorful() {
	for i := 0; i < 255; i++ {
		fmt.Printf("\x1b[48;5;%dm%3d\u001B[0m", i, i)
		if i == 15 || (i > 15 && ((i-15)%6 == 0)) {
			println()
		}
	}
}

func printTitleDefault(helpInfo HelpInfo) {
	PrintTitle(os.Args[0], helpInfo)
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
