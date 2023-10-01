package core

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/ctool/stream"
	"github.com/kuangcp/logger"
	"github.com/tidwall/pretty"
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

	ProxySwitch interface {
		HasUse() bool
		SwitchUse()
		GetName() string
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
		Id          string        `json:"id"`
		Redis       *RedisConf    `json:"redis"`
		Groups      []*ProxyGroup `json:"groups"`
		ProxySelf   *ProxySelf    `json:"proxy"`  // 抓包
		ProxyDirect *ProxySelf    `json:"direct"` // 不抓包，不存储
	}
)

var (
	// 文件相对路径
	configFilePath = "/.dev-proxy/dev-proxy.json"
	logFilePath    = "/.dev-proxy/dev-proxy.log"
	pacFilePath    = "/.dev-proxy/dev-proxy.pac"
	dbDirPath      = "/.dev-proxy/leveldb-request-log"
)

var (
	ProxyConfVar ProxyConf
	// url map: src -> dst
	replaceMap   = make(map[string]string)
	trackList    []string // 代理类型的地址（抓包）
	directList   []string // 直连类型的地址
	lock         = &sync.RWMutex{}
	ConfigReload = make(chan bool, 1)
	guiMode      = false
	// DirectType https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
	DirectType = []string{"html", "javascript", "css", "image/", "pdf", "msword", "octet-stream", "audio", "video"}
)

func (p *ProxyGroup) GetName() string {
	return p.Name
}

func (g *ProxyGroup) HasUse() bool {
	return g.ProxyType == Open
}

func (g *ProxyGroup) SwitchUse() {
	if g.HasUse() {
		g.ProxyType = Close
	} else {
		g.ProxyType = Open
	}
}

func (p *ProxySelf) GetName() string {
	return p.Name
}

func (g *ProxySelf) HasUse() bool {
	return g.ProxyType == Open
}

func (g *ProxySelf) SwitchUse() {
	if g.HasUse() {
		g.ProxyType = Close
	} else {
		g.ProxyType = Open
	}
}

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

func MatchNeedStorage(proxyReq http.Request) bool {
	fullUrl := proxyReq.URL.Scheme + "://" + proxyReq.URL.Host + proxyReq.URL.Path
	return stream.Just(directList...).NoneMatch(func(item any) bool {
		return matchConf(item.(string), fullUrl)
	})
}

func FindReplaceByRegexp(proxyReq http.Request) (*url.URL, string) {
	lock.RLock()
	defer lock.RUnlock()

	fullUrl := proxyReq.URL.Scheme + "://" + proxyReq.URL.Host + proxyReq.URL.Path
	matchKey := ""
	matchResult := ""

	for _, conf := range directList {
		if matchConf(conf, fullUrl) {
			return nil, Direct
		}
	}

	for k, v := range replaceMap {
		tryResult := tryToReplacePath(k, v, fullUrl)
		if tryResult == "" {
			continue
		}

		// 多个匹配时，按最长正则匹配
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

		return parse, Replace
	}

	for _, conf := range trackList {
		if matchConf(conf, fullUrl) {
			parse, err := url.Parse(fullUrl)
			if err != nil {
				logger.Error(err)
			}

			return parse, Proxy
		}
	}

	return nil, Direct
}

func InitConfig() {
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
	if Debug {
		logLevel = logger.DebugDesc
	}
	logger.SetLoggerConfig(&logger.LogConfig{
		TimeFormat: "01-02 15:04:05.000",
		File: &logger.FileLogger{
			Filename:   home + logFilePath,
			Level:      logger.DebugDesc,
			Append:     true,
			PermitMask: "0660",
			MaxDays:    -1,
		},
		Console: &logger.ConsoleLogger{
			Level:    logLevel,
			Colorful: true,
		}})

	if JsonPath != "" {
		configFilePath = JsonPath
	} else {
		configFilePath = home + configFilePath
	}
	if PacPath != "" {
		pacFilePath = PacPath
	} else {
		pacFilePath = home + pacFilePath
	}

	dbDirPath = home + dbDirPath
	exist := ctool.IsFileExist(configFilePath)
	if !exist {
		initMainProxyJson()
	}

	cleanAndRegisterFromFile(configFilePath)

	var hostId string
	if ProxyConfVar.Id == "" {
		hostId = hostname + ":tmp-" + uuid.NewString()[:6]
	} else {
		hostId = hostname + ":" + ProxyConfVar.Id
	}

	RequestList = fmt.Sprintf(listFmt, Prefix, hostId)
	RequestUrlList = fmt.Sprintf(urlListFmt, Prefix, hostId)
	if ReloadConf {
		go listenConfig()
	}
}

func MarkGuiMode() {
	guiMode = true
}

func initMainProxyJson() {
	logger.Warn("init new config")
	conf := ProxyConf{
		Id: uuid.NewString()[:5],
		Redis: &RedisConf{
			Addr:     "127.0.0.1:6379",
			DB:       0,
			PoolSize: 3,
		},
		Groups: []*ProxyGroup{
			{Name: "temp", ProxyType: Close, Routers: []ProxyRouter{
				{Src: "http://127.0.0.1/(.*)", Dst: "http://127.0.0.1/$1", ProxyType: Open},
			}},
		},
		ProxyDirect: &ProxySelf{Name: "direct", ProxyType: Open, Paths: []string{"http://172.22.133.255:8989/(.*)"}},
		ProxySelf:   &ProxySelf{Name: "proxy", ProxyType: Open, Paths: []string{"http://172.22.133.255:8990/(.*)"}},
	}
	storeByMemory(conf)
}

func storeByMemory(conf ProxyConf) {
	bts, err := json.Marshal(conf)
	if err != nil {
		logger.Error(err)
		return
	}
	var Options = &pretty.Options{Width: 80, Prefix: "", Indent: "  ", SortKeys: false}
	fmtBts := pretty.PrettyOptions(bts, Options)

	os.WriteFile(configFilePath, fmtBts, 0644)
}

func cleanAndRegisterFromFile(configFile string) {
	file, err := os.ReadFile(configFile)
	if err != nil {
		logger.Error(err)
		return
	}

	lock.Lock()
	defer lock.Unlock()

	err = json.Unmarshal(file, &ProxyConfVar)
	if err != nil {
		logger.Error(err)
		return
	}

	ReloadConfByCacheObj()
	if guiMode {
		ConfigReload <- true
	}
}

func ReloadConfByCacheObj() {
	logger.Info("Start reload proxy rule by file")

	// 替换
	replaceMap = make(map[string]string)
	for _, conf := range ProxyConfVar.Groups {
		if conf.ProxyType == Close {
			continue
		}
		logger.Info("Register group:", conf.Name)
		for _, router := range conf.Routers {
			if router.ProxyType == Close {
				continue
			}
			replaceMap[router.Src] = router.Dst
			logger.Debug("Register", ctool.White.Print(router.Src),
				ctool.Yellow.Print("►"), ctool.Cyan.Print(router.Dst))
		}
	}

	// 抓包，存储
	trackList = []string{}
	parsePath(ProxyConfVar.ProxySelf, "track", func(s string) {
		trackList = append(trackList, s)
	})

	// 不抓包，不存储
	directList = []string{}
	parsePath(ProxyConfVar.ProxyDirect, "direct", func(s string) {
		directList = append(directList, s)
	})
	logger.Info("Finish reload proxy rule by file")
}

func parsePath(proxy *ProxySelf, name string, listAppendFunc func(string)) {
	if proxy == nil || proxy.ProxyType != Open {
		return
	}
	logger.Debug("Register %v: %v \n %v", name, proxy.Name, strings.Join(proxy.Paths, " \n "))

	for _, path := range proxy.Paths {
		if path == "" {
			continue
		}
		listAppendFunc(path)
	}
}

func listenConfig() {
	var lastModTime = time.Now()
	for range time.NewTicker(time.Second * 2).C {
		stat, err := os.Stat(configFilePath)
		if err != nil {
			logger.Error(err)
			continue
		}

		curModTime := stat.ModTime()
		if curModTime.After(lastModTime) {
			//logger.Info(stat.ModTime())
			execCommand("notify-send -i folder-new Dev-Proxy 'start reload config file'")
			lastModTime = curModTime
			cleanAndRegisterFromFile(configFilePath)
		}
	}
}
