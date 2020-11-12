package main

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/logger"
)

var info = cuibase.HelpInfo{
	Description:   "Recycle bin",
	Version:       "1.0.2",
	SingleFlagLen: -3,
	ValueLen:      -10,
	Flags: []cuibase.ParamVO{
		{
			Short:   "-h",
			Value:   "",
			Comment: "Help info and init",
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
		},{
			Short:   "-q",
			Value:   "",
			Comment: "Remove pid file",
		},
	},
	Options: []cuibase.ParamVO{
		{
			Short:   "",
			Value:   "file",
			Comment: "Delete file",
		}, {
			Short:   "-s",
			Value:   "suffix",
			Comment: "Delete *.suffix",
		}, {
			Short:   "-l",
			Value:   "duration",
			Comment: "File live time, default 1m. (unit: ms/s/m/h) ",
		}, {
			Short:   "-c",
			Value:   "duration",
			Comment: "Check trash period, default 1m. (unit: ms/s/m/h)",
		},
	}}

func InitDir() {
	err := os.MkdirAll(trashDir, 0755)
	cuibase.CheckIfError(err)
	err = os.MkdirAll(configDir, 0755)
	cuibase.CheckIfError(err)
	err = os.MkdirAll(logDir, 0755)
	cuibase.CheckIfError(err)
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
