package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/kuangcp/gobase/pkg/ctool"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	_ "net/http/pprof"

	"github.com/kuangcp/logger"
)

var (
	port       int
	version    bool
	serverAddr string
	localHost  string
	localAddr  string
	syncDir    string
	checkSec   int
)

var (
	sideList   = ctool.NewSet[string]() // 对端列表 格式 host:port
	configName = ".ksync.config.json"
)

type (
	ArgVO struct {
		ServerAddr string `json:"server_addr"`
		LocalHost  string `json:"local_host"`
		LocalPort  int    `json:"local_port"`
		SyncDir    string `json:"sync_dir"`
	}
)

func (a ArgVO) String() string {
	bt, err := json.Marshal(a)
	if err != nil {
		return "Error: " + err.Error()
	}
	return string(bt)
}
func init() {
	flag.IntVar(&port, "p", 8000, "port")
	flag.IntVar(&checkSec, "c", 5, "check duration second")
	flag.BoolVar(&version, "v", false, "version")
	flag.StringVar(&serverAddr, "s", "", "init server host and port. ag: 192.168.0.4:8000")
	flag.StringVar(&localHost, "l", "", "local side host. ag: 192.168.0.5")
	flag.StringVar(&syncDir, "d", "", "sync dir.")
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("1.0.3")
		return
	}

	initFromConfig()
	normalizeParam()

	watchFileThenSync()

	logger.Info("start success. server:", serverAddr, "local:", localHost)
	go displayRemoteSide()
	go pprofDebug()
	webServer()
}

func displayRemoteSide() {
	for range time.NewTicker(time.Minute).C {
		logger.Info("Remote:", sideList.Join(","))
	}
}

func pprofDebug() {
	debugPort := "33054"
	logger.Info("http://127.0.0.1:" + debugPort + "/debug/pprof/")
	_ = http.ListenAndServe("0.0.0.0:"+debugPort, nil)
}

// 当命令行未指定才从配置文件加载
func initFromConfig() {
	if serverAddr != "" && localHost != "" {
		return
	}

	home, err := ctool.Home()
	if err == nil {
		configName = home + "/" + configName
	}
	file, err := os.ReadFile(configName)
	if err != nil {
		logger.Error(err)
		return
	}

	var arg ArgVO
	err = json.Unmarshal(file, &arg)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("load config", arg)
	if arg.LocalPort != 0 {
		port = arg.LocalPort
	}
	if serverAddr == "" {
		serverAddr = arg.ServerAddr
	}
	if syncDir == "" && arg.SyncDir != "" {
		syncDir = arg.SyncDir
	}
	if localHost == "" {
		localHost = arg.LocalHost
	}
}

func normalizeParam() {
	if syncDir == "" {
		syncDir = "./"
	}
}

func absPath(filename string) string {
	if "windows" == runtime.GOOS {
		return syncDir + "\\" + filename
	} else {
		return syncDir + "/" + filename
	}
}
func trimAbs(path string) string {
	if "windows" == runtime.GOOS {
		return strings.TrimPrefix(path, syncDir+"\\")
	} else {
		return strings.TrimPrefix(path, syncDir+"/")
	}
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// 接收模式：注册自身到远程
func registerOnRemote() {
	if serverAddr == "" {
		return
	}
	// 必须要声明自身可被访问的host
	if localHost == "" {
		logger.Fatal("local host is missing")
	}

	req, _ := http.NewRequest("GET", "http://"+serverAddr+"/register", nil)
	req.Header.Set("self", localAddr)

	for range time.NewTicker(time.Second * 5).C {
		rsp, err := http.DefaultClient.Do(req)
		if rsp != nil && rsp.Body != nil {
			rsp.Body.Close()
		}
		if err != nil {
			logger.Error(err)
			return
		}
		sideList.Add(serverAddr)
	}
}

// 检查本地可访问性
func checkLocalHostReachable(localAddr string) {
	resp, err := http.DefaultClient.Get("http://" + localAddr + "/ping")
	if err != nil {
		logger.Fatal("local host unreachable", localAddr, err)
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Fatal("local host unreachable", localAddr, err)
		}
		respStr := string(all)
		if respStr != pong {
			logger.Fatal("un except local host", localAddr, respStr)
		}
	}
}

func watchFileThenSync() {
	if serverAddr != "" {
		localAddr = fmt.Sprintf("%v:%v", localHost, port)
		go func() {
			time.Sleep(time.Second * 5)
			checkLocalHostReachable(localAddr)
		}()
	}

	go registerOnRemote()
	go watchSyncFile()
}

// 发送模式：检查文件变化 发送到对端
func watchSyncFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//logger.Info("event:", event)
				if event.Has(fsnotify.Create) {
					sideList.Loop(func(addr string) {
						postFile(addr, trimAbs(event.Name))
					})
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error("error:", err)
			}
		}
	}()

	logger.Info("Start watch ", syncDir)
	err = watcher.Add(syncDir)
	if err != nil {
		log.Fatal(err)
	}

	// 为什么需要 block，当前web环境下主协程是不会退出的
	<-make(chan struct{})
}

func postFile(server string, file string) {
	name := url.QueryEscape(file)
	existURL := "http://" + server + "/exist?name=" + name

	existRsp, err := http.Get(existURL)
	if err != nil || existRsp == nil || existRsp.Body == nil {
		return
	}
	rspStr, err := ioutil.ReadAll(existRsp.Body)
	if err != nil || string(rspStr) == "EXIST" {
		logger.Info("%s exist", name)
		return
	}
	defer existRsp.Body.Close()

	open, err := os.Open(absPath(file))
	if err != nil {
		logger.Error(err)
		return
	}
	defer open.Close()

	postRsp, err := http.Post("http://"+server+"/upload?name="+name, "", open)
	if err != nil || postRsp == nil || postRsp.Body == nil {
		logger.Error(err)
		return
	}
	defer postRsp.Body.Close()
	logger.Info("send to", server, name)
}
