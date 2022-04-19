package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	_ "net/http/pprof"

	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/logger"
)

var (
	port       int
	version    bool
	serverAddr string
	localAddr  string
	localHost  string
	syncDir    string
	checkSec   int
)

var (
	lastFile   = cuibase.NewSet()
	sideList   = cuibase.NewSet() // 对端列表 格式 host:port
	configName = "ksync.config.json"
)

type (
	ArgVO struct {
		ServerAddr string `json:"server_addr"`
		LocalHost  string `json:"local_host"`
	}
)

func init() {
	flag.IntVar(&port, "p", 8000, "port")
	flag.IntVar(&checkSec, "c", 2, "check duration second")
	flag.BoolVar(&version, "v", false, "version")
	flag.StringVar(&serverAddr, "s", "", "init server host&port. ag: 192.168.0.1:8000")
	flag.StringVar(&localHost, "l", "", "local side host. ag: 192.168.0.2")
	flag.StringVar(&syncDir, "d", "./", "sync dir.")
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("1.0.2")
		return
	}

	initFromConfig()
	//pprofDebug()
	logger.Info("start success. server:", serverAddr, "local:", localHost)
	normalizeParam()

	go syncTimerTask()
	webServer()
}

func pprofDebug() {
	debugPort := "8891"
	go func() {
		fmt.Println("http://127.0.0.1:" + debugPort + "/debug/pprof/")
		_ = http.ListenAndServe("0.0.0.0:"+debugPort, nil)
	}()
}

// 当命令行未指定才从配置文件加载
func initFromConfig() {
	if serverAddr != "" && localHost != "" {
		return
	}

	file, err := ioutil.ReadFile(configName)
	if err != nil {
		return
	}

	var arg ArgVO
	err = json.Unmarshal(file, &arg)
	if err != nil {
		return
	}

	if serverAddr == "" {
		serverAddr = arg.ServerAddr
	}
	if localHost == "" {
		localHost = arg.LocalHost
	}
}

func normalizeParam() {
	localAddr = localHost + ":" + fmt.Sprint(port)
	if "windows" == runtime.GOOS {
		//if !strings.HasSuffix(syncDir, "\\") {
		//	syncDir += "\\"
		//}
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")

		if home == "" {

			home = os.Getenv("USERPROFILE")

		}
		fmt.Println(home)
	} else {
		if !strings.HasSuffix(syncDir, "/") {
			syncDir += "/"
		}
	}
}

func readNeedSyncFile() []string {
	var result []string
	dir, err := ioutil.ReadDir(syncDir)
	if err != nil {
		fmt.Println(err)
		return result
	}

	firstInit := lastFile.IsEmpty()
	for _, info := range dir {
		if info.IsDir() {
			continue
		}
		fileName := info.Name()
		if strings.HasPrefix(fileName, ".") {
			continue
		}
		if fileName == configName {
			continue
		}
		contains := lastFile.Contains(fileName)
		lastFile.Add(fileName)

		if !contains || firstInit {
			logger.Info("need sync", fileName, info.ModTime())
			result = append(result, fileName)
		}
	}
	return result
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// 注册到服务端
func registerOnServer() {
	// 如果需要注册到服务端，必须要声明自己可被访问的host
	if serverAddr != "" && localHost == "" {
		logger.Warn("local host is missing")
		os.Exit(1)
	}

	if serverAddr == "" {
		return
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", "http://"+serverAddr+"/register", nil)
	req.Header.Set("self", fmt.Sprintf("%v:%v", localHost, port))
	if err != nil {
		fmt.Println(err)
		return
	}

	rsp, err := client.Do(req)
	if rsp != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		fmt.Println(rsp)
		fmt.Println(err)
		return
	}
	sideList.Add(serverAddr)
}

func syncTimerTask() {
	ticker := time.NewTicker(time.Second * time.Duration(checkSec))
	for range ticker.C {
		registerOnServer()
		syncFile()
	}
}

func syncFile() {
	if sideList.Len() == 0 {
		return
	}

	//logger.Info("check sync %v", sideList)
	fileList := readNeedSyncFile()
	for _, path := range fileList {
		sideList.Loop(func(i interface{}) {
			postFile(i.(string), path)
		})
	}
}

func postFile(server string, path string) {
	name := url.QueryEscape(path)
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

	open, err := os.Open(path)
	if err != nil {
		logger.Error(err)
		return
	}
	defer open.Close()

	postRsp, err := http.Post("http://"+server+"/upload?name="+name, "", open)
	if err != nil || postRsp == nil || postRsp.Body == nil {
		return
	}
	defer postRsp.Body.Close()
	logger.Info("send to ", server, name)
}
