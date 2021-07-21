package app

import (
	"bufio"
	"fmt"
	"github.com/kuangcp/logger"
	"io"
	"os"
	"runtime"

	"github.com/kuangcp/gobase/pkg/cuibase"
)

const (
	groupDirStr     = "group/"
	bakFileStr      = "origin.hosts.bak"
	winHostFileStr  = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	unixHostFileStr = "/etc/hosts"
)

var (
	mainDir     = "/.hosts-group/" // 用户目录下
	groupDir    string
	bakFile     string
	curHostFile string
)
var (
	Debug            bool
	DebugStatic      bool
	Version          bool
	LogPath          string
	FinalHostFile    string // 入参指定hosts文件
	MainPath         string
	GenerateAfterCmd string
)

var Info = cuibase.HelpInfo{
	Description:   "Hosts Group, switch host tool",
	Version:       "1.4.0",
	SingleFlagLen: -2,
	DoubleFlagLen: 0,
	ValueLen:      -5,
	Flags: []cuibase.ParamVO{
		{Short: "-h", Comment: "help info"},
		{Short: "-f", Comment: "set main hosts file"},
		{Short: "-m", Comment: "main config file"},
		{Short: "-cmd", Comment: "the cmd run after generate hosts file"},
		{Short: "-d", Comment: "debug mode, use test dir and host-file"},
		{Short: "-D", Comment: "debug static mode, use static dir not embed packaged"},
		{Short: "-v", Comment: "version"},
	},
	Options: []cuibase.ParamVO{
		{Short: "-l", Value: "path", Comment: "log path"},
	},
}

// 初始化日志配置，目标hosts文件，应用配置目录
func InitConfigAndEnv() {
	initLogConfig()

	if MainPath != "" {
		mainDir = MainPath
	}

	home, err := cuibase.Home()
	cuibase.CheckIfError(err)

	mainDir = home + mainDir
	groupDir = mainDir + groupDirStr
	bakFile = mainDir + bakFileStr
	mkDir(groupDir)

	fillFinalHostsFile()
	logger.Info("current hosts file:", curHostFile)

	backupOriginFile()
}

func fillFinalHostsFile() {
	if FinalHostFile != "" {
		curHostFile = FinalHostFile
		return
	}

	if Debug {
		logger.Info("using debug mode")
		curHostFile = mainDir + "hosts"
		return
	}

	if runtime.GOOS == "windows" {
		curHostFile = winHostFileStr
	} else {
		curHostFile = unixHostFileStr
	}
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

// 备份原始 hosts 文件
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
