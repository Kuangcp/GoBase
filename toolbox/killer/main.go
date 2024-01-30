package main

import (
	"fmt"
	"github.com/kuangcp/logger"
	"os/exec"
	"strings"
)

func main() {
	pid, err := GetPid("firefox")
	if err != nil {
		logger.Error(err)
	}

	fmt.Println(pid)

	pids := strings.Split(pid, "\n")
	for _, id := range pids {
		runInLinux("kill " + id)
	}
}

// GetPid 根据进程名称获取进程ID
func GetPid(serverName string) (string, error) {
	a := `ps ux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	pid, err := runInLinux(a)
	return pid, err
}

func runInLinux(cmd string) (string, error) {
	//fmt.Println("Running Linux cmd:" + cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}
