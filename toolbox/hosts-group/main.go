package main

import (
	"flag"
	"fmt"

	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/hosts-group/app"
	"github.com/kuangcp/gobase/toolbox/hosts-group/app/icon"
	"github.com/kuangcp/logger"
	"github.com/skratchdot/open-golang/open"
)

func init() {
	flag.BoolVar(&app.Debug, "d", false, "")
}

func main() {
	flag.Parse()

	app.InitPrepare()

	go func() {
		app.WebServer("8066")
	}()

	onExit := func() {
		logger.Info("exit")
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Hosts Group")
	systray.SetTooltip("Tips")
	mQuitOrig := systray.AddMenuItem("Exit", "Exit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("Hosts Group")
		systray.SetTooltip("Hosts Group")

		//subMenuTop := systray.AddMenuItem("Groups", "SubMenu Test (top)")
		//mChecked := subMenuTop.AddSubMenuItemCheckbox("Unchecked", "Check Me", true)

		mUrl := systray.AddMenuItem("Open UI", "my home")

		for {
			select {
			//case <-mChecked.ClickedCh:
			//	if mChecked.Checked() {
			//		mChecked.Uncheck()
			//		mChecked.SetTitle("Unchecked")
			//	} else {
			//		mChecked.Check()
			//		mChecked.SetTitle("Checked")
			//	}
			case <-mUrl.ClickedCh:
				open.Run("http://localhost:8066")
			}
		}
	}()
}
