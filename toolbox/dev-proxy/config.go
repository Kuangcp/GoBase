package main

import (
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type ProxyConf struct {
	Name    string   `json:"name"`
	Use     bool     `json:"use"`
	Routers []string `json:"routers"`
}

// origin url -> target url
var proxyMap = make(map[*url.URL]*url.URL)
var lock = &sync.RWMutex{}

//	/api/a -> /a
//
// registerReplace(map[string]string{"http://host1:port1/api": "http://host2:port2"})
//
//	/api/a -> /api2/a
//
// registerReplace(map[string]string{"http://host1:port1/api": "http://host2:port2/api2"})
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

// TODO 当前按Host维度替换，需要实现按路径维度替换
func findTargetReplace(r *http.Request) (*url.URL, *url.URL) {
	lock.RLock()
	defer lock.RUnlock()
	for k, v := range proxyMap {
		if k.Host == r.Host {
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

	configFile := home + "/.dev-proxy.json"
	cleanAndRegister(configFile)
	if checkConf {
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
			if !conf.Use {
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
