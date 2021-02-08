package app

import (
	"os"
	"os/exec"
	"sync"

	"github.com/kuangcp/logger"
	"github.com/skratchdot/open-golang/open"

	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/hosts-group/app/icon"
)

var (
	fileMap sync.Map
)

func OnReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Hosts Group")
	systray.SetTooltip("Hosts Group")

	addPageLinkItem()

	versionItem := systray.AddMenuItem("v"+Info.Version, Info.Version)
	versionItem.Disable()
	exitItem := systray.AddMenuItem("Exit", "Exit the whole app")
	go func() {
		<-exitItem.ClickedCh
		logger.Info("Requesting quit")
		systray.Quit()
		logger.Info("Finished quitting")
	}()

	systray.AddSeparator()

	list := getFileList()
	var latch sync.WaitGroup
	for _, vo := range list {
		latch.Add(1)
		addFileItem(vo, &latch)
		latch.Wait()
	}
}

func addPageLinkItem() {
	winItem := systray.AddMenuItem("Web", "Web")
	pageURL := systray.AddMenuItem("Hosts Group", "page")
	feedbackURL := systray.AddMenuItem("Feedback", "Feedback")
	go func() {
		for {
			select {
			case <-winItem.ClickedCh:
				command := exec.Command(os.Args[0], "-w")
				err := command.Start()
				if err != nil {
					logger.Fatal(err.Error())
				}
			case <-feedbackURL.ClickedCh:
				err := open.Run("https://github.com/Kuangcp/GoBase/issues")
				if err != nil {
					logger.Fatal(err.Error())
				}
			case <-pageURL.ClickedCh:
				err := open.Run("http://localhost:8066")
				if err != nil {
					logger.Fatal(err.Error())
				}
			}
		}
	}()
}

func addFileItem(vo FileItemVO, s *sync.WaitGroup) {
	go func() {
		checkbox := systray.AddMenuItemCheckbox(vo.Name, "Check Me", vo.Use)
		fileMap.Store(vo.Name, checkbox)
		if s != nil {
			s.Done()
		}
		for {
			select {
			case <-checkbox.ClickedCh:
				useState, err := fileUseState(vo.Name)
				if err != nil {
					logger.Warn("switch failed", err)
					break
				}
				success, err := switchFileState(vo.Name)
				if !success {
					logger.Warn("switch failed", err)
					systray.AddMenuItem("Error: "+err.Error(), "")
					// rollback check action
					updateCheckBox(!useState, checkbox)
					break
				}
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
