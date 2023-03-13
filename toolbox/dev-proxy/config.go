package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type (
	ProxyRouter struct {
		Src       string `json:"src"`
		Dst       string `json:"dst"`
		ProxyType int    `json:"proxy_type"`
	}

	ProxyGroup struct {
		Name      string        `json:"name"`
		ProxyType int           `json:"proxy_type"`
		Routers   []ProxyRouter `json:"routers"`
	}

	ProxySelf struct {
		Name      string   `json:"name"`
		ProxyType int      `json:"proxy_type"`
		Paths     []string `json:"paths"`
	}

	RedisConf struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		DB       int    `json:"db"`
		PoolSize int    `json:"pool_size"`
	}

	ProxyConf struct {
		Groups    []ProxyGroup `json:"groups"`
		ProxySelf *ProxySelf   `json:"proxy"`
		Redis     *RedisConf   `json:"redis"`
		Id        string       `json:"id"`
	}
)

const (
	configPath = "/.dev-proxy/dev-proxy.json"
	logPath    = "/.dev-proxy/dev-proxy.log"
	Open       = 1
	Close      = 0
	Proxy      = 2
)

var (
	proxyConf     ProxyConf
	proxyValMap   = make(map[string]string)
	proxySelfList []string
	lock          = &sync.RWMutex{}
	dbPath        = "/.dev-proxy/leveldb-request-log"
)

// 处理源路径到目标路径的转换
// originConf 正则匹配规则
// targetConf 正则替换规则
// fullUrl    实际请求的完整路径

// 例如：
//
//	原始 http://192.168.9.12:30011/api/(.*)       目标 http://127.0.0.1:8081/api/v2/$1
//	请求 http://192.168.9.12:30011/api/pageList   结果 http://127.0.0.1:8081/api/v2/pageList
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

func matchConf(originConf, fullUrl string) bool {
	compile, err := regexp.Compile(originConf)
	if err != nil {
		logger.Error(err)
		return false
	}

	return compile.Match([]byte(fullUrl))
}

func findReplaceByRegexp(proxyReq http.Request) (*url.URL, int) {
	lock.RLock()
	defer lock.RUnlock()

	fullUrl := proxyReq.URL.Scheme + "://" + proxyReq.URL.Host + proxyReq.URL.Path
	matchKey := ""
	matchResult := ""
	for k, v := range proxyValMap {
		tryResult := tryToReplacePath(k, v, fullUrl)
		if tryResult == "" {
			continue
		}

		//logger.Debug("match", k, tryResult)
		if len(matchKey) < len(k) {
			matchResult = tryResult
			matchKey = k
		}
	}
	if len(matchResult) > 0 {
		//logger.Debug("Final", matchKey, matchResult)
		parse, err := url.Parse(matchResult)
		if err != nil {
			logger.Error(err)
		}

		return parse, Open
	}

	for _, conf := range proxySelfList {
		if matchConf(conf, fullUrl) {
			parse, err := url.Parse(fullUrl)
			if err != nil {
				logger.Error(err)
			}

			return parse, Proxy
		}
	}

	return nil, Close
}

func initConfig() {
	home, err := ctool.Home()
	if err != nil {
		fmt.Println(err)
		return
	}
	hostname, err := os.Hostname()
	if err != nil {
		logger.Error("Random hostname. err:", err)
		hostname = uuid.NewString()
	}

	logLevel := logger.InformationalDesc
	if debug {
		logLevel = logger.DebugDesc
	}
	logger.SetLoggerConfig(&logger.LogConfig{
		TimeFormat: "01-02 15:04:05.000",
		File: &logger.FileLogger{
			Filename:   home + logPath,
			Level:      logLevel,
			Append:     true,
			PermitMask: "0660",
			MaxDays:    -1,
		}})

	configFile := home + configPath
	dbPath = home + dbPath
	cleanAndRegisterFromFile(configFile)

	if proxyConf.Id == "" {
		listVar += ":" + hostname + ":tmp-" + uuid.NewString()[:6]
	} else {
		listVar += ":" + hostname + ":" + proxyConf.Id
	}

	RequestList = listVar
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

	err = json.Unmarshal(file, &proxyConf)
	if err != nil {
		logger.Error(err)
		return
	}

	proxyValMap = make(map[string]string)
	for _, conf := range proxyConf.Groups {
		if conf.ProxyType == Close {
			continue
		}
		logger.Info("Register group:", conf.Name)
		for _, router := range conf.Routers {
			if router.ProxyType == Close {
				continue
			}
			proxyValMap[router.Src] = router.Dst
			logger.Debug("Register", router.Src, ctool.Yellow.Print("►"), ctool.Cyan.Print(router.Dst))
		}
	}

	// 代理自身
	if proxyConf.ProxySelf != nil {
		logger.Info("Register proxy group:", proxyConf.ProxySelf.Name)
		logger.Debug("Register %v", strings.Join(proxyConf.ProxySelf.Paths, " , "))

		for _, path := range proxyConf.ProxySelf.Paths {
			if path == "" {
				continue
			}
			proxySelfList = append(proxySelfList, path)
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
			execCommand("notify-send -i folder-new Dev-Proxy 'start reload config file'")
			lastModTime = curModTime
			cleanAndRegisterFromFile(configFile)
		}
	}
}
