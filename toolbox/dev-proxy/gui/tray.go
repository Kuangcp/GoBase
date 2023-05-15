package main

import (
	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/logger"
	"reflect"
	"sync"
)

var (
	fileMap sync.Map
)

func OnExit() {
	logger.Info("exit")
	//storeByMemory()
}

func OnReady() {
	systray.SetTemplateIcon(Data, Data)
	systray.SetTitle("Dev Proxy")
	systray.SetTooltip("Dev Proxy")

	exitItem := systray.AddMenuItem("Exit  (v"+core.HelpInfo.Version+")", "Exit the whole app")
	go func() {
		<-exitItem.ClickedCh
		logger.Info("Requesting quit")
		systray.Quit()
		logger.Info("Finished quitting")
	}()

	systray.AddSeparator()

	for _, vo := range core.ProxyConfVar.Groups {
		addFileItem(vo)
	}
	addFileItem(core.ProxyConfVar.ProxySelf)
	addFileItem(core.ProxyConfVar.ProxyBlock)
}

func addFileItem(vo core.ProxySwitch) {
	if vo == nil || (reflect.ValueOf(vo).Kind() == reflect.Ptr && reflect.ValueOf(vo).IsNil()) {
		return
	}
	checkbox := systray.AddMenuItemCheckbox(vo.GetName(), "Check Me", vo.HasUse())
	fileMap.Store(vo.GetName(), checkbox)
	go func() {
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
