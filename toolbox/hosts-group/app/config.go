package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/kuangcp/logger"

	"github.com/kuangcp/gobase/pkg/cuibase"
)

const (
	groupDirStr     = "group/"
	bakFileStr      = "origin.hosts.bak"
	winHostFileStr  = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	unixHostFileStr = "/etc/hosts"
)

var (
	mainDir     = "/.hosts-group/"
	groupDir    string
	bakFile     string
	curHostFile string
)
var (
	Debug   bool
	Win     bool
	Version bool
	LogPath string
)

var Info = cuibase.HelpInfo{
	Description:   "Hosts switch tool",
	Version:       "1.3.6",
	SingleFlagLen: -2,
	DoubleFlagLen: 0,
	ValueLen:      -5,
	Flags: []cuibase.ParamVO{
		{Short: "-h", Comment: "help info"},
		{Short: "-d", Comment: "debug"},
		{Short: "-v", Comment: "version"},
	},
	Options: []cuibase.ParamVO{
		{Short: "-l", Value: "path", Comment: "log path"},
	},
}

func InitConfigAndEnv() {
	initLogConfig()

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

	if Debug {
		logger.Info("using debug mode")
		curHostFile = mainDir + "hosts"
	}
	logger.Info("current hosts file:", curHostFile)

	mkDir(groupDir)

	backupOriginFile()
}

func initLogConfig() {
	logger.SetLogPathTrim("/hosts-group/")

	if LogPath == "" {
		return
	}
	exists, err := isPathExists(LogPath)
	if !exists || err != nil {
		logger.Fatal("log path invalid")
	}

	err = logger.SetLoggerConfig(&logger.LogConfig{
		Console: &logger.ConsoleLogger{
			Level:    logger.DebugDesc,
			Colorful: true,
		},
		File: &logger.FileLogger{
			Filename:   LogPath,
			Level:      logger.DebugDesc,
			Colorful:   true,
			Append:     true,
			PermitMask: "0660",
		},
	})
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func backupOriginFile() {
	exists, err := isPathExists(bakFile)
	cuibase.CheckIfError(err)
	if !exists {
		CopyFile(curHostFile, bakFile)
		CopyFile(curHostFile, groupDir+"origin-backup"+use)
		err := generateHost()
		if err != nil {
			logger.Fatal(err.Error())
		}
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
