package main

import (
	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/logger"
	"sync"
)

type (
	FileItemVO struct {
		Name    string `json:"name"`
		Use     bool   `json:"use"`
		Content string `json:"content,omitempty"`
	}
)

var (
	fileMap sync.Map
)

func OnExit() {
	logger.Info("exit")
	//storeByMemory()
}

func OnReady() {
	//systray.SetTemplateIcon(Data, Data)
	systray.SetTitle("Hosts Group")
	systray.SetTooltip("Hosts Group")

	//addPageLinkItem()

	//versionItem := systray.AddMenuItem("v"+Info.Version, Info.Version)
	//versionItem.Disable()
	exitItem := systray.AddMenuItem("Exit", "Exit the whole app")
	go func() {
		<-exitItem.ClickedCh
		logger.Info("Requesting quit")
		systray.Quit()
		logger.Info("Finished quitting")
	}()

	systray.AddSeparator()

	var latch sync.WaitGroup
	for _, vo := range core.ProxyConfVar.Groups {
		latch.Add(1)
		addFileItem(vo, &latch)
		latch.Wait()
	}
}

func addFileItem(vo *core.ProxyGroup, s *sync.WaitGroup) {
	go func() {
		checkbox := systray.AddMenuItemCheckbox(vo.Name, "Check Me", vo.HasUse())
		fileMap.Store(vo.Name, checkbox)
		if s != nil {
			s.Done()
		}
		for {
			select {
			case <-checkbox.ClickedCh:
				useState := vo.HasUse()

				vo.SwitchUse()
				core.ReloadConfByCacheObj()

				// Windows need this line, linux not need
				updateCheckBox(useState, checkbox)
			}
		}
	}()
}

func updateCheckBox(useState bool, checkbox *systray.MenuItem) {
	if useState {
		checkbox.Uncheck()
	} else {
		checkbox.Check()
	}
}
