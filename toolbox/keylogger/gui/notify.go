package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"os"
	"os/exec"
	"time"
)

// notify something

type (
	NotifyMsg struct {
		Duration   string `json:"duration"`
		Pic        string `json:"pic"`
		Title      string `json:"title"`
		Content    string `json:"content"`
		ExpireTime int    `json:"expireTime"`
	}
)

func notifyAny() {
	home, err := ctool.Home()
	if err != nil {
		logger.Error(err)
		return
	}

	jsonFile := home + "/.config/app-conf/notify/notify.json"
	var msgs []NotifyMsg
	file, err := os.ReadFile(jsonFile)
	if err != nil {
		logger.Error(err)
		return
	}
	err = json.Unmarshal(file, &msgs)
	if err != nil {
		logger.Error(err)
		return
	}
	for _, msg := range msgs {
		logger.Info("add notify:", msg.Title)
		duration, err := time.ParseDuration(msg.Duration)
		if err != nil {
			logger.Error(err)
			continue
		}
		for t := range time.NewTicker(duration).C {
			logger.Info(t.String(), msg.Title)
			execCommand(fmt.Sprintf("notify-send -i %s %s %s -t %v", msg.Pic, msg.Title, msg.Content, msg.ExpireTime))
		}
	}
}

func execCommand(command string) (string, bool) {
	cmd := exec.Command("/usr/bin/bash", "-c", command)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.Error(command, err)
		return "", false
	}

	result := out.String()
	return result, true
}
