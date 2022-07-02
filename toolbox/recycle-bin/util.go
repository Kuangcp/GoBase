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
	listTrash    bool
	log          bool
	showConfig   bool
	initConfig   bool
	pipeline     bool
	listOrder    int
	buildVersion string
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

	flag.IntVar(&listOrder, "o", 0, "")

	flag.StringVar(&restore, "R", "", "")
	flag.StringVar(&retentionStr, "r", retentionStr, "")
	flag.StringVar(&periodStr, "p", periodStr, "")
	flag.StringVar(&suffix, "s", "", "")
}

func initConfigValue() {
	logger.SetLogPathTrim("recycle-bin/")

	home, err := ctk.Home()
	ctk.CheckIfError(err)

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
	ctk.CheckIfError(err)
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

var info = ctk.HelpInfo{
	Description:   "Recycle bin",
	Version:       "1.0.7",
	BuildVersion:  buildVersion,
	SingleFlagLen: -3,
	ValueLen:      -10,
	Flags: []ctk.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "Help info"},
		{Short: "-P", BoolVar: &pipeline, Comment: "Pipeline"},
		{Short: "-D", BoolVar: &debug, Comment: "Release mode"},
		{Short: "-X", BoolVar: &exit, Comment: "Exit daemon"},
		{Short: "-C", BoolVar: &check, Comment: "Start check"},
		{Short: "-d", BoolVar: &daemon, Comment: "Start check by daemon"},
		{Short: "-l", BoolVar: &listTrash, Comment: "List trash"},
		{Short: "-g", BoolVar: &log, Comment: "Show log file path"},
		{Short: "-c", BoolVar: &showConfig, Comment: "Show config file"},
		{Short: "-i", BoolVar: &initConfig, Comment: "Init dir and config"},
	},
	Options: []ctk.ParamVO{
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
	ctk.CheckIfError(err)
	err = os.MkdirAll(configDir, 0755)
	ctk.CheckIfError(err)
	err = os.MkdirAll(logDir, 0755)
	ctk.CheckIfError(err)
	exist, err := isPathExists(configFile)
	ctk.CheckIfError(err)
	if !exist {
		result, err := json.Marshal(Setting{Retention: retentionStr, CheckPeriod: periodStr})
		ctk.CheckIfError(err)
		err = ioutil.WriteFile(configFile, result, 0644)
		ctk.CheckIfError(err)
	}
}

func isPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	//fmt.Println("stat:  ",path , err)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
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
