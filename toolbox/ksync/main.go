package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/kuangcp/logger"
)

var (
	sideList  []string // 对端列表 格式 host:port
	port      int
	version   bool
	initSide  string
	localAddr string
	syncDir   string
)

var (
	lastFile = make(map[string]struct{})
)

func init() {
	flag.IntVar(&port, "p", 8000, "port")
	flag.BoolVar(&version, "v", false, "version")
	flag.StringVar(&initSide, "s", "", "init server side. ag: 192.168.0.1:8000")
	flag.StringVar(&localAddr, "l", "", "local side host. ag: 192.168.0.2")
	flag.StringVar(&syncDir, "d", "./", "sync dir. end with / or \\\\")
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("1.0.0")
		return
	}

	initSideBind()

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {
			syncFile()
		}
	}()

	webServer()
}

func readNeedSyncFile() []string {
	var result []string
	dir, err := ioutil.ReadDir(syncDir)
	if err != nil {
		fmt.Println(err)
		return result
	}

	init := len(lastFile) == 0
	for _, info := range dir {
		if info.IsDir() {
			continue
		}
		_, ok := lastFile[info.Name()]
		lastFile[info.Name()] = struct{}{}

		if !ok || init {
			logger.Info("need sync", info.Name(), info.ModTime())
			result = append(result, info.Name())
		}
	}
	return result
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func webServer() {
	http.HandleFunc("/exist", func(writer http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		unescape, err := url.QueryUnescape(name)
		if err != nil {
			logger.Error(err)
			return
		}
		exist := isFileExist(syncDir + "" + unescape)
		if exist {
			writer.Write([]byte("EXIST"))
		} else {
			writer.Write([]byte("NONE"))
		}
	})

	// 接收文件
	http.HandleFunc("/sync", func(writer http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		unescape, err := url.QueryUnescape(name)
		if err != nil {
			logger.Error(err)
			return
		}

		open, err := os.Create(syncDir + unescape)
		if err != nil {
			logger.Error(err)
			return
		}

		var buf = make([]byte, 4096)
		for {
			read, err := request.Body.Read(buf)
			if read != 0 {
				open.Write(buf[:read])
			}
			if err != nil {
				break
			}
		}

		open.Close()
	})
	// 注册
	http.HandleFunc("/register", func(writer http.ResponseWriter, request *http.Request) {
		client := request.Header.Get("self")
		logger.Info("register new", client)
		sideList = append(sideList, client)
		writer.Write([]byte("OK"))
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func initSideBind() {
	if initSide == "" || localAddr == "" {
		logger.Warn("client config missing")
		return
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", "http://"+initSide+"/register", nil)
	req.Header.Set("self", fmt.Sprintf("%v:%v", localAddr, port))
	if err != nil {
		fmt.Println(err)
		return
	}
	rsp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
	sideList = append(sideList, initSide)
}

func syncFile() {
	if len(sideList) == 0 {
		return
	}

	//logger.Info("check sync %v", sideList)
	fileList := readNeedSyncFile()
	for _, path := range fileList {
		for _, r := range sideList {
			postFile(r, path)
		}
	}
}

func postFile(server string, path string) {
	name := url.QueryEscape(path)
	existURL := "http://" + server + "/exist?name=" + name

	existRsp, err := http.Get(existURL)
	if existRsp == nil || err != nil {
		return
	}
	rspStr, err := ioutil.ReadAll(existRsp.Body)
	if string(rspStr) == "EXIST" {
		logger.Info("%s exist", name)
		return
	}

	open, err := os.Open(path)
	if err != nil {
		logger.Error(err)
		return
	}

	defer open.Close()

	syncURL := "http://" + server + "/sync?name=" + name
	post, err := http.Post(syncURL, "", open)
	if err != nil {
		return
	}

	defer post.Body.Close()
	logger.Info("send to ", server, name)
}
