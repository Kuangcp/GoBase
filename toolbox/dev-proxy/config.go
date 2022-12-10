package main

import (
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type ProxyConf struct {
	Name    string   `json:"name"`
	Enable  bool     `json:"enable"`
	Routers []string `json:"routers"`
}

// origin url -> target url
var proxyMap = make(map[*url.URL]*url.URL)
var lock = &sync.RWMutex{}

func registerReplace(proxy map[string]string) {
	lock.Lock()
	proxyMap = make(map[*url.URL]*url.URL)
	for k, v := range proxy {
		kUrl, err := url.Parse(k)
		if err != nil {
			continue
		}
		vUrl, err := url.Parse(v)
		if err != nil {
			continue
		}
		proxyMap[kUrl] = vUrl
		logger.Info("register:", k, "=>", v)
	}
	lock.Unlock()
}

// TODO 待实现源路径支持正则替换
func needReplace(originUrl *url.URL, request *http.Request) bool {
	sameHost := originUrl.Host == request.Host
	samePath := originUrl.Path == "" || strings.HasPrefix(request.URL.Path, originUrl.Path)
	return sameHost && samePath
}

func findTargetReplace(r *http.Request) (*url.URL, *url.URL) {
	lock.RLock()
	defer lock.RUnlock()
	for k, v := range proxyMap {
		if needReplace(k, r) {
			return k, v
		}
	}
	return nil, nil
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

		registerReplace(proxy)
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
