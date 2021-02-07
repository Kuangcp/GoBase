package app

import (
	"github.com/kuangcp/logger"
	"github.com/skratchdot/open-golang/open"

	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/hosts-group/app/icon"
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

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("Hosts Group")
		systray.SetTooltip("Hosts Group")

		addGroupItem()
		addPageItem()
	}()
}

func addPageItem() {
	mUrl := systray.AddMenuItem("Open Page", "page")
	for {
		select {
		case <-mUrl.ClickedCh:
			err := open.Run("http://localhost:8066")
			if err != nil {
				logger.Fatal(err.Error())
			}
		}
	}
}

func addGroupItem() {
	groupMenu := systray.AddMenuItem("Groups", "SubMenu Test (top)")
	list := getFileList()
	for _, vo := range list {
		vo := vo
		go func() {
			checkbox := groupMenu.AddSubMenuItemCheckbox(vo.Name, "Check Me", vo.Use)
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
}

func OnExit() {
	logger.Info("exit")
}
