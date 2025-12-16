package main

import (
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"
)

type ProxyConf struct {
	Name    string   `json:"name"`
	Enable  int      `json:"enable"`
	Routers []string `json:"routers"`
}

const (
	configPath = "/.dev-proxy/dev-proxy.json"
	logPath    = "/.dev-proxy/dev-proxy.log"
)

var (
	proxyValMap = make(map[string]string)
	lock        = &sync.RWMutex{}
	dbPath      = "/.dev-proxy/leveldb-request-log"
)

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

func findReplaceByRegexp(proxyReq http.Request) *url.URL {
	lock.RLock()
	defer lock.RUnlock()

	for k, v := range proxyValMap {
		fullUrl := proxyReq.URL.Scheme + "://" + proxyReq.URL.Host + proxyReq.URL.Path
		tryResult := tryToReplacePath(k, v, fullUrl)
		if tryResult == "" {
			continue
		}

		parse, err := url.Parse(tryResult)
		if err != nil {
			logger.Error(err)
		}

		return parse
	}

	return nil
}

func initConfig() {
	home, err := ctool.Home()
	if err != nil {
		fmt.Println(err)
		return
	}

	logger.SetLoggerConfig(&logger.LogConfig{
		TimeFormat: "01-02 15:04:05.000",
		File: &logger.FileLogger{
			Filename:   home + logPath,
			Level:      logger.DebugDesc,
			Append:     true,
			PermitMask: "0660",
			MaxDays:    -1,
		}})

	configFile := home + configPath
	dbPath = home + dbPath
	cleanAndRegisterFromFile(configFile)
	if reloadConf {
		go listenConfig(configFile)
	}
}

func cleanAndRegisterFromFile(configFile string) {
	file, err := os.ReadFile(configFile)
	if err != nil {
		logger.Error(err)
		return
	}

	lock.Lock()
	defer lock.Unlock()

	var confList []ProxyConf
	err = json.Unmarshal(file, &confList)
	if err != nil {
		logger.Error(err)
		return
	}

	proxyValMap = make(map[string]string)
	for _, conf := range confList {
		if conf.Enable == 0 {
			continue
		}
		logger.Info("Register group:", conf.Name)
		pair := len(conf.Routers) / 2
		for i := 0; i < pair; i++ {
			match := conf.Routers[i*2]
			replace := conf.Routers[i*2+1]
			proxyValMap[match] = replace
			logger.Debug("Register", match, ctool.Yellow.Print("►"), ctool.Cyan.Print(replace))
		}
	}
}

func listenConfig(configFile string) {
	var lastModTime = time.Now()
	for range time.NewTicker(time.Second * 2).C {
		stat, err := os.Stat(configFile)
		if err != nil {
			logger.Error(err)
			continue
		}

		curModTime := stat.ModTime()
		if curModTime.After(lastModTime) {
			//logger.Info(stat.ModTime())
			lastModTime = curModTime
			cleanAndRegisterFromFile(configFile)
		}
	}
}
