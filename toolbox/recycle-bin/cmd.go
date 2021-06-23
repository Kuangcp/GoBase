package main

import (
	"bytes"
	"github.com/kuangcp/logger"
	"os"
	"os/exec"
)

// 静默执行 不关心返回值
func execCmdWithQuite(cmd *exec.Cmd) {
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.Error(cmd, err)
		os.Exit(1)
	}
}

func execCommand(command string) string {
	cmd := exec.Command("/usr/bin/bash", "-c", command)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	result := out.String()
	return result
}
