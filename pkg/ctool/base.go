package ctool

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

// CheckIfError assert err is nil
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Home
// eg: /home/test
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

func CopyStream(src io.ReadCloser) ([]byte, io.ReadCloser) {
	bodyBt, err := io.ReadAll(src)
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	return bodyBt, io.NopCloser(bytes.NewBuffer(bodyBt))
}

func ConvertList[T any, R any](src []T, mapFun func(T) R, filterFun func(T) bool) []R {
	var result []R
	for _, d := range src {
		if filterFun == nil || filterFun(d) {
			result = append(result, mapFun(d))
		}
	}
	return result
}
