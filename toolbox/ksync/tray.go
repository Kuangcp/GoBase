package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/kuangcp/logger"
	"ksync/icon"
	"net/http"
	"net/url"
	"os"
	"runtime"
)

func OnReady() {
	if "windows" == runtime.GOOS {
		systray.SetTemplateIcon(icon.Data, icon.Data)
	} else {
		systray.SetTemplateIcon([]byte(iconImg), []byte(iconImg))
	}
	systray.SetTitle("K-Sync")
	systray.SetTooltip("K-Sync")

	syncItem := systray.AddMenuItem("Sync", "sync file")
	go func() {
		for {
			select {
			case <-syncItem.ClickedCh:
				logger.Info("sync %v", sideList)
				fileList := syncFile()
				for _, path := range fileList {
					for _, r := range sideList {
						postFile(r, path)
					}
				}
			}
		}
	}()

	exitItem := systray.AddMenuItem(fmt.Sprintf("Exit (%v)", port), "Exit the whole app")
	go func() {
		<-exitItem.ClickedCh
		logger.Info("Requesting quit")
		systray.Quit()
		logger.Info("Finished quitting")
	}()

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
