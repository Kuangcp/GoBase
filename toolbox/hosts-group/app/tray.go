package app

import (
	"github.com/pkg/errors"
	"os/exec"
	"os/user"
	"strconv"
	"sync"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/kuangcp/gobase/toolbox/hosts-group/app/icon"
	"github.com/kuangcp/logger"
	"github.com/skratchdot/open-golang/open"
)

var (
	fileMap sync.Map
)

func OnReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Hosts Group")
	systray.SetTooltip("Hosts Group")

	//addPageLinkItem()

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
	//winItem := systray.AddMenuItem("Webview", "Webview")
	pageURL := systray.AddMenuItem("Hosts Group", "page")
	feedbackURL := systray.AddMenuItem("Feedback", "Feedback")
	go func() {
		for {
			select {
			//case <-winItem.ClickedCh:
			//	go OpenWebView("http://localhost:"+PortStr)
			case <-feedbackURL.ClickedCh:

				err := open.Run("https://github.com/Kuangcp/GoBase/issues")
				if err != nil {
					logger.Fatal(err.Error())
				}
			case <-pageURL.ClickedCh:
				openWithNormalUser("http://localhost:" + PortStr)
				//open.Run("http://localhost:" + PortStr)
			}
		}
	}()
}

func openWithNormalUser(input string) {
	cmd := exec.Command("xdg-open", input)
	//cmd := exec.Command("sleep", "10m")
	err := setUserAttr(cmd, "zk")
	if err != nil {
		logger.Error(err)
		return
	}

	err = cmd.Run()
	if err != nil {
		logger.Error(err)
	}
}

// 修改启动进程所属用户
func setUserAttr(cmd *exec.Cmd, name string) error {
	// 检测用户是否存在
	sysUser, err := user.Lookup(name)
	logger.Info(sysUser.Uid, sysUser.Gid)
	if err != nil {
		return errors.Wrapf(err, "invalid user %s", name)
	}
	// set process attr
	// 获取用户 id
	uid, err := strconv.ParseUint(sysUser.Uid, 10, 32)
	if err != nil {
		return err
	}
	// 获取用户组 id
	gid, err := strconv.ParseUint(sysUser.Gid, 10, 32)
	if err != nil {
		return err
	}
	attr := cmd.SysProcAttr

	logger.Info("attr: ", attr)
	if attr == nil {
		attr = &syscall.SysProcAttr{}
	}
	//设置进程执行用户
	attr.Credential = &syscall.Credential{
		Uid:         uint32(uid),
		Gid:         uint32(gid),
		NoSetGroups: false,
	}

	cmd.SysProcAttr = attr
	return nil
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
