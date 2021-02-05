package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/logger"
)

var (
	mainDir  = "/.hosts-group"
	groupDir = "/group"
	bakFile  = "/hosts.origin.bak"
)

func init() {
	initPrepare()
}

func initPrepare() {
	if "windows" == runtime.GOOS {
		logger.Fatal("not support")
	}

	home, err := cuibase.Home()
	cuibase.CheckIfError(err)

	mainDir = home + mainDir
	groupDir = mainDir + groupDir
	bakFile = mainDir + bakFile

	mkDir(groupDir)

	exists, err := isPathExists(bakFile)
	cuibase.CheckIfError(err)
	if !exists {
		CopyFile("/etc/hosts", bakFile)
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
