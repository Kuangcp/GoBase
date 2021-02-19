package app

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/kuangcp/gobase/pkg/cuibase"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	width              = 100
	height             = 20
	refreshLabelPeriod = time.Second
	appId              = "com.github.kuangcp.keylogger"
)

var (
	app      *gtk.Application
	win      *gtk.Window
	kpmLabel *gtk.Label
)

func ShowPopWindow() {
	gtk.Init(nil)
	app, _ = gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	_, err := app.Connect("activate", createWindow)
	cuibase.CheckIfError(err)
	app.Run(nil)
}

func createWindow() {
	win, _ = gtk.WindowNew(gtk.WINDOW_POPUP)
	win.Add(buildWidget())
	win.SetDefaultSize(width, height)
	win.SetPosition(gtk.WIN_POS_MOUSE)
	_, err := win.Connect("destroy", gtk.MainQuit)
	cuibase.CheckIfError(err)
	bindMouseActionForWindow()

	// init label
	now := time.Now()
	refreshLabel(now)

	app.AddWindow(win)
	win.ShowAll()

	// 启动后的计算并刷新缓存
	go func() {
		ticker := time.NewTicker(refreshLabelPeriod)
		for now := range ticker.C {
			refreshLabel(now)
		}
	}()
}

func buildWidget() *gtk.Widget {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	kpmLabel, err = gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	grid.Add(kpmLabel)
	kpmLabel.SetHExpand(true)
	kpmLabel.SetVExpand(true)

	return &grid.Container.Widget
}

func bindMouseActionForWindow() {
	// 鼠标按下事件
	var x, y int
	win.SetEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK))

	//鼠标按下事件处理
	_, _ = win.Connect("button-press-event", func(widget *gtk.Window, ctx *gdk.Event) {
		//获取鼠键按下属性结构体变量，系统内部的变量，不是用户传参变量
		event := *(*gdk.EventButton)(unsafe.Pointer(&ctx))
		if event.Button() == 1 { //左键
			x, y = int(event.X()), int(event.Y()) //保存点击的起点坐标
		} else if event.Button() == 3 { //右键
			//右键，关闭窗口
			app.Quit()
		}
	})

	//鼠标移动事件处理
	_, _ = win.Connect("motion-notify-event", func(widget *gtk.Window, ctx *gdk.Event) {
		//获取鼠标移动属性结构体变量，系统内部的变量，不是用户传参变量
		event := *(*gdk.EventButton)(unsafe.Pointer(&ctx))
		win.Move(int(event.XRoot())-x, int(event.YRoot())-y)
	})
}

// 从缓存中更新窗口内面板
func refreshLabel(now time.Time) {
	conn := GetConnection()
	today := now.Format(DateFormat)

	tempValue, err := conn.Get(GetTodayTempKPMKeyByString(today)).Result()
	if err != nil {
		tempValue = "0"
	}
	maxValue, err := conn.Get(GetTodayMaxKPMKeyByString(today)).Result()
	if err != nil {
		maxValue = "0"
	}
	total := conn.ZScore(TotalCount, today).Val()

	// https://blog.csdn.net/bitscro/article/details/3874616
	str := fmt.Sprintf(" 🕒 %s\n%s %s %s",
		fmt.Sprintf("<span foreground='#F2F3F5' font_desc='10'>%s</span>", now.Format(TimeFormat)),
		fmt.Sprintf("<span foreground='#5AFF00' font_desc='14'>%s</span>", tempValue),
		fmt.Sprintf("<span foreground='#F2F3F5' font_desc='12'>%d</span>", int(total)),
		fmt.Sprintf("<span foreground='yellow' font_desc='9'>%s</span>", maxValue),
	)
	_, err = glib.IdleAdd(kpmLabel.SetMarkup, str)
	if err != nil {
		log.Fatal("IdleAdd() failed:", err)
	}
}
