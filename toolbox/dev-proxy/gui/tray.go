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
	core.Go(func() {
		<-exitItem.ClickedCh
		logger.Info("Requesting quit")
		systray.Quit()
		logger.Info("Finished quitting")
	})

	systray.AddSeparator()

	for _, vo := range core.ProxyConfVar.Groups {
		addFileItem(vo)
	}
	systray.AddSeparator()
	addFileItem(core.ProxyConfVar.ProxySelf)
	addFileItem(core.ProxyConfVar.ProxyDirect)

	core.Go(refreshUIByConfigReload)
}

func refreshUIByConfigReload() {
	for {
		select {
		case <-core.ConfigReload:
			logger.Info("config reload, start refresh ui")
			fileMap.Range(func(key, value any) bool {
				updateCheckBox(key.(core.ProxySwitch).HasUse(), value.(*systray.MenuItem))
				return true
			})
		}
	}
}

func addFileItem(vo core.ProxySwitch) {
	if vo == nil || (reflect.ValueOf(vo).Kind() == reflect.Ptr && reflect.ValueOf(vo).IsNil()) {
		return
	}
	checkbox := systray.AddMenuItemCheckbox(vo.GetName(), "Check Me", vo.HasUse())
	fileMap.Store(vo, checkbox)
	core.Go(func() {
		for {
			select {
			case <-checkbox.ClickedCh:
				useState := vo.HasUse()

				vo.SwitchUse()
				core.ReloadConfByCacheObj()

				// Windows need this line, linux not need
				// 手动更新内存标记为正确的值 选中or不选中
				updateCheckBox(!useState, checkbox)
			}
		}
	})
}

func updateCheckBox(use bool, checkbox *systray.MenuItem) {
	if use {
		checkbox.Check()
	} else {
		checkbox.Uncheck()
	}
}
