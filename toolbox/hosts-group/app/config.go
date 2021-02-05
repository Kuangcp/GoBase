package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/kuangcp/gobase/pkg/cuibase"
)

const (
	groupDirStr     = "group/"
	bakFileStr      = "hosts.origin.bak"
	winHostFileStr  = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	unixHostFileStr = "/etc/hosts"
)

var (
	mainDir     = "/.hosts-group/"
	groupDir    string
	bakFile     string
	curHostFile string
)

func init() {
	initPrepare()
}

func initPrepare() {
	if "windows" == runtime.GOOS {
		curHostFile = winHostFileStr
	} else {
		curHostFile = unixHostFileStr
	}

	home, err := cuibase.Home()
	cuibase.CheckIfError(err)

	mainDir = home + mainDir
	groupDir = mainDir + groupDirStr
	bakFile = mainDir + bakFileStr

	mkDir(groupDir)

	exists, err := isPathExists(bakFile)
	cuibase.CheckIfError(err)
	if !exists {
		CopyFile(curHostFile, bakFile)
	}
}

func mkDir(path string) {
	pathExists, err := isPathExists(path)
	cuibase.CheckIfError(err)
	if !pathExists {
		err := os.MkdirAll(path, 0755)
		cuibase.CheckIfError(err)
	}
}

func isPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CopyFile(srcFileName string, dstFileName string) (written int64, err error) {
	srcFile, err := os.Open(srcFileName)

	if err != nil {
		fmt.Printf("open file err = %v\n", err)
		return
	}

	defer srcFile.Close()

	//通过srcFile，获取到Reader
	reader := bufio.NewReader(srcFile)

	//打开dstFileName
	dstFile, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open file err = %v\n", err)
		return
	}

	writer := bufio.NewWriter(dstFile)
	defer func() {
		writer.Flush() //把缓冲区的内容写入到文件
		dstFile.Close()

	}()

	return io.Copy(writer, reader)
}
