package main

import (
	"github.com/getlantern/systray"
	"github.com/kuangcp/logger"
	"ksync/icon"
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

	exitItem := systray.AddMenuItem("Exit", "Exit the whole app")
	go func() {
		<-exitItem.ClickedCh
		logger.Info("Requesting quit")
		systray.Quit()
		logger.Info("Finished quitting")
	}()

	systray.AddSeparator()
	syncItem := systray.AddMenuItem("Sync", "sync file")
	go func() {
		for {
			select {
			case <-syncItem.ClickedCh:
				logger.Info("sync")
			}
		}
	}()
}

func OnExit() {
	logger.Info("exit")
}
