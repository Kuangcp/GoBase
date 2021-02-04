package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/logger"
)

var (
	help         bool
	suffix       string
	check        bool
	daemon       bool
	debug        bool
	exit         bool
	illegalQuit  bool
	listTrash    bool
	log          bool
	showConfig   bool
	initConfig   bool
	listOrder    int
	restore      string
	retentionStr = "168h" // time.ParseDuration()
	periodStr    = "1h"
)

type (
	Setting struct {
		Retention   string `json:"retention"`
		CheckPeriod string `json:"checkPeriod"`
	}
)

func init() {
	initConfigValue()

	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "H", false, "")
	flag.BoolVar(&debug, "D", false, "")
	flag.BoolVar(&check, "C", false, "")
	flag.BoolVar(&daemon, "d", false, "")
	flag.BoolVar(&exit, "X", false, "")
	flag.BoolVar(&illegalQuit, "q", false, "")
	flag.BoolVar(&listTrash, "l", false, "")
	flag.BoolVar(&log, "g", false, "")
	flag.BoolVar(&showConfig, "c", false, "")
	flag.BoolVar(&initConfig, "i", false, "")
	flag.IntVar(&listOrder, "o", 0, "")

	flag.StringVar(&restore, "R", "", "")
	flag.StringVar(&retentionStr, "r", retentionStr, "")
	flag.StringVar(&periodStr, "p", periodStr, "")
	flag.StringVar(&suffix, "s", "", "")

	flag.Usage = info.PrintHelp
}

func initConfigValue() {
	logger.SetLogPathTrim("recycle-bin")

	home, err := cuibase.Home()
	cuibase.CheckIfError(err)

	mainDir = home + mainDir
	configDir = mainDir + "/conf"
	pidFile = configDir + "/pid"
	configFile = configDir + "/main.json"

	logDir = mainDir + "/log"
	logFile = logDir + "/main.log"

	trashDir = mainDir + "/trash"

	_ = logger.SetLoggerConfig(&logger.LogConfig{
		Console: &logger.ConsoleLogger{
			Level:    logger.DebugDesc,
			Colorful: true,
		},
		File: &logger.FileLogger{
			Filename:   logFile,
			Level:      logger.DebugDesc,
			Colorful:   true,
			Append:     true,
			PermitMask: "0660",
		},
	})

	_, err = loadConfig()
	cuibase.CheckIfError(err)
}

func loadConfig() (*Setting, error) {
	exists, err := isPathExists(configFile)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	setting := Setting{}
	err = json.Unmarshal(file, &setting)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if setting.Retention != "" {
		retentionStr = setting.Retention
	}
	if setting.CheckPeriod != "" {
		periodStr = setting.CheckPeriod
	}
	return &setting, nil
}

var info = cuibase.HelpInfo{
	Description:   "Recycle bin",
	Version:       "1.0.6",
	SingleFlagLen: -3,
	ValueLen:      -10,
	Flags: []cuibase.ParamVO{
		{
			Short:   "-h",
			Value:   "",
			Comment: "Help info",
		}, {
			Short:   "-D",
			Value:   "",
			Comment: "Debug mode",
		}, {
			Short:   "-X",
			Value:   "",
			Comment: "Exit daemon",
		}, {
			Short:   "-C",
			Value:   "",
			Comment: "Start check",
		}, {
			Short:   "-d",
			Value:   "",
			Comment: "Start check by daemon",
		}, {
			Short:   "-q",
			Value:   "",
			Comment: "Remove pid file",
		}, {
			Short:   "-l",
			Value:   "",
			Comment: "List trash",
		}, {
			Short:   "-g",
			Value:   "",
			Comment: "Show log file path",
		}, {
			Short:   "-c",
			Value:   "",
			Comment: "Show config file",
		}, {
			Short:   "-i",
			Value:   "",
			Comment: "Init dir and config",
		},
	},
	Options: []cuibase.ParamVO{
		{
			Short:   "",
			Value:   "file",
			Comment: "Delete file",
		}, {
			Short:   "-R",
			Value:   "file",
			Comment: "Restore file",
		}, {
			Short:   "-o",
			Value:   "order",
			Comment: "Order for list(1/2 asc/desc)",
		}, {
			Short:   "-s",
			Value:   "suffix",
			Comment: "Delete *.suffix",
		}, {
			Short:   "-r",
			Value:   "duration",
			Comment: "File retention time, default " + retentionStr + ". (unit: ms/s/m/h) ",
		}, {
			Short:   "-p",
			Value:   "duration",
			Comment: "Check trash period, default " + periodStr + ". (unit: ms/s/m/h)",
		},
	}}

func InitConfig() {
	fmt.Println("init")
	err := os.MkdirAll(trashDir, 0755)
	cuibase.CheckIfError(err)
	err = os.MkdirAll(configDir, 0755)
	cuibase.CheckIfError(err)
	err = os.MkdirAll(logDir, 0755)
	cuibase.CheckIfError(err)
	exist, err := isPathExists(configFile)
	cuibase.CheckIfError(err)
	if !exist {
		result, err := json.Marshal(Setting{Retention: retentionStr, CheckPeriod: periodStr})
		cuibase.CheckIfError(err)
		err = ioutil.WriteFile(configFile, result, 0644)
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

func newSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}

func startProc(args []string, logFile string) (*exec.Cmd, error) {
	cmd := &exec.Cmd{
		Path:        args[0],
		Args:        args,
		SysProcAttr: newSysProcAttr(),
	}

	if logFile != "" {
		stdout, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			logger.Error("打开日志文件错误: ", os.Getpid(), err)
			return nil, err
		}
		cmd.Stderr = stdout
		//cmd.Stdout = stdout
	}

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func invokeWithBool(flag bool, action func()) {
	if flag {
		action()
		os.Exit(0)
	}
}

func invokeWithStr(param string, action func(string)) {
	if param != "" {
		action(param)
		os.Exit(0)
	}
}
