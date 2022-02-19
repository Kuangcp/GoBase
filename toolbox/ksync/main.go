package main

import (
	_ "embed"
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

//go:embed sync.png
var iconImg string

var (
	sideList  []string // 对端列表 格式 host:port
	port      int
	version   bool
	initSide  string
	localAddr string
)

var (
	lastFile = make(map[string]struct{})
)

func init() {
	flag.IntVar(&port, "p", 8000, "port")
	flag.BoolVar(&version, "v", false, "version")
	flag.StringVar(&initSide, "s", "", "init server side. ag: 192.168.0.1:8000")
	flag.StringVar(&localAddr, "l", "", "local side host. ag: 192.168.0.2")
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("1.0.0")
		return
	}

	initSideBind()

	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {

	}
	go func() {
		for {
			select {
			case <-ticker.C:
				syncFile()
			}
		}
	}()

	webServer()
}

func readNeedSyncFile() []string {
	var result []string
	dir, err := ioutil.ReadDir("./")
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
		if init {
			lastFile[info.Name()] = struct{}{}
		}

		if !ok || init {
			logger.Info("need sync", info.Name(), info.ModTime())
			result = append(result, info.Name())
		}

	}
	return result
}

func webServer() {
	http.HandleFunc("/sync", func(writer http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		unescape, err := url.QueryUnescape(name)
		if err != nil {
			logger.Error(err)
			return
		}

		open, err := os.Create(unescape)
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
	http.HandleFunc("/register", func(writer http.ResponseWriter, request *http.Request) {
		client := request.Header.Get("self")
		logger.Info("add new", client)
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
		fmt.Println("config error")
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
