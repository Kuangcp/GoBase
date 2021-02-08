package app

import (
	"sync"

	"github.com/kuangcp/logger"
	"github.com/skratchdot/open-golang/open"

	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/hosts-group/app/icon"
)

var (
	//groupMenu *systray.MenuItem
	fileMap sync.Map
)

func OnReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Hosts Group")
	systray.SetTooltip("Tips")

	exitItem := systray.AddMenuItem("Exit", "Exit the whole app")
	go func() {
		<-exitItem.ClickedCh
		logger.Info("Requesting quit")
		systray.Quit()
		logger.Info("Finished quitting")
	}()

	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Hosts Group")
	systray.SetTooltip("Hosts Group")

	systray.AddMenuItem("v"+Info.Version, Info.Version)
	addPageLinkItem()
	systray.AddSeparator()

	list := getFileList()
	for _, vo := range list {
		addFileItem(vo)
	}
}

func addPageLinkItem() {
	mUrl := systray.AddMenuItem("Open Page", "page")
	go func() {
		for {
			select {
			case <-mUrl.ClickedCh:
				err := open.Run("http://localhost:8066")
				if err != nil {
					logger.Fatal(err.Error())
				}
			}
		}
	}()
}

func addFileItem(vo FileItemVO) {
	go func() {
		checkbox := systray.AddMenuItemCheckbox(vo.Name, "Check Me", vo.Use)
		fileMap.Store(vo.Name, checkbox)
		for {
			select {
			case <-checkbox.ClickedCh:
				state, err := fileState(vo.Name)
				if err != nil {
					logger.Warn("switch failed", err)
					break
				}
				if state {
					checkbox.Uncheck()
				} else {
					checkbox.Check()
				}
				success, err := switchFileState(vo.Name)
				if !success {
					logger.Warn("switch failed", err)
				}
			}
		}
	}()
}

func updateFileItemState(vo FileItemVO) {
	value, ok := fileMap.Load(vo.Name)
	if ok {
		if vo.Use {
			value.(*systray.MenuItem).Check()
		} else {
			value.(*systray.MenuItem).Uncheck()
		}
	}
}

func OnExit() {
	logger.Info("exit")
}
