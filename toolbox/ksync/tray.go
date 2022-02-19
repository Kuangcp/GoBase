package main

import (
	"time"

	"net/http"
	"net/url"
	"os"

	"github.com/kuangcp/logger"
)

func OnReady() {

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			select {
			case <-ticker.C:
				syncFile()
			}
		}
	}()
}

func syncFile() {
	logger.Info("check sync %v", sideList)
	fileList := readNeedSyncFile()
	for _, path := range fileList {
		for _, r := range sideList {
			postFile(r, path)
		}
	}
}

func postFile(server string, path string) {
	logger.Info(path, "post to", server)
	client := http.Client{}
	open, err := os.Open(path)
	if err != nil {
		logger.Error(err)
		return
	}

	defer open.Close()
	name := url.QueryEscape(path)
	syncURL := "http://" + server + "/sync?name=" + name
	//fmt.Println(syncURL)
	post, err := client.Post(syncURL, "", open)
	if err != nil {
		return
	}

	defer post.Body.Close()
	logger.Info(post, err)
}

func OnExit() {
	logger.Info("exit")
}
