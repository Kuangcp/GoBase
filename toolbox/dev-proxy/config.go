package main

import (
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"os"
	"regexp"
	"sync"
	"time"
)

type ProxyConf struct {
	Name    string   `json:"name"`
	Enable  bool     `json:"enable"`
	Routers []string `json:"routers"`
}

var proxyValMap = make(map[string]string)
var lock = &sync.RWMutex{}

// 处理源路径到目标路径的转换
// originConf 正则匹配规则
// targetConf 正则替换规则
// fullUrl    实际请求的完整路径

// 例如：
//
//	原始 http://192.168.9.12:30011/api/(.*)
//	目标 http://127.0.0.1:8081/api/v2/$1
//	请求 http://192.168.9.12:30011/api/pageList
//	结果 http://127.0.0.1:8081/api/v2/pageList
func tryToReplacePath(originConf, targetConf, fullUrl string) string {
	compile, err := regexp.Compile(originConf)
	if err != nil {
		logger.Error(err)
		return ""
	}

	replaceResult := compile.ReplaceAllString(fullUrl, targetConf)
	if replaceResult == fullUrl {
		return ""
	}

	return replaceResult
}

func initConfig() {
	home, err := ctool.Home()
	if err != nil {
		fmt.Println(err)
		return
	}

	logger.SetLoggerConfig(&logger.LogConfig{
		TimeFormat: "01-02 15:04:05",
		File: &logger.FileLogger{
			Filename:   home + "/.dev-proxy.log",
			Level:      logger.DebugDesc,
			Append:     true,
			PermitMask: "0660",
			MaxDays:    -1,
		}})

	configFile := home + "/.dev-proxy.json"
	cleanAndRegister(configFile)
	if reloadConf {
		go listenConfig(configFile)
	}
}

func cleanAndRegister(configFile string) {
	file, err := os.ReadFile(configFile)
	if err == nil {
		var confList []ProxyConf
		err := json.Unmarshal(file, &confList)
		if err != nil {
			logger.Error(err)
			return
		}
		proxy := make(map[string]string)
		for _, conf := range confList {
			if !conf.Enable {
				continue
			}
			logger.Info("Register group:", conf.Name)
			pair := len(conf.Routers) / 2
			for i := 0; i < pair; i++ {
				proxy[conf.Routers[i*2]] = conf.Routers[i*2+1]
			}
		}

		proxyValMap = proxy
	}
}

func listenConfig(configFile string) {
	var lastModTime = time.Now()
	for range time.NewTicker(time.Second * 5).C {
		stat, err := os.Stat(configFile)
		if err != nil {
			logger.Error(err)
		}
		curModTime := stat.ModTime()
		if curModTime.After(lastModTime) {
			//logger.Info(stat.ModTime())
			lastModTime = curModTime
			cleanAndRegister(configFile)
		}
	}
}
