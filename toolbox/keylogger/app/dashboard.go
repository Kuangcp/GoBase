package app

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/cuibase"
	"log"
	"time"
	"unsafe"

	"github.com/go-redis/redis"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/wonderivan/logger"
)

const (
	width           = 128
	height          = 24
	refreshPeriod   = time.Millisecond * 520
	recordThreshold = 58 // xx s后才存储 bpm
)

var (
	bpmLabel   *gtk.Label
	firstStart = true
)

func ShowWindow() {
	gtk.Init(nil)
	app, _ := gtk.ApplicationNew("com.github.kuangcp.keylogger", glib.APPLICATION_FLAGS_NONE)
	_, err := app.Connect("activate", func() {
		createWindow(app)
	})
	cuibase.CheckIfError(err)
	app.Run(nil)
}

func createWindow(app *gtk.Application) {
	win, _ := gtk.WindowNew(gtk.WINDOW_POPUP)
	win.Add(windowWidget())
	app.AddWindow(win)

	win.SetTitle("Dashboard")
	win.SetDefaultSize(width, height)
	_, err := win.Connect("destroy", gtk.MainQuit)
	cuibase.CheckIfError(err)
	win.SetPosition(gtk.WIN_POS_MOUSE)
	win.ShowAll()

	// 鼠标按下事件
	var x, y int
	win.SetEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK))

	//鼠标按下事件处理
	_, _ = win.Connect("button-press-event", func(widget *gtk.Window, ctx *gdk.Event) {
		//获取鼠键按下属性结构体变量，系统内部的变量，不是用户传参变量
		event := *(*gdk.EventButton)(unsafe.Pointer(&ctx))
		//x, y = int(event.X()), int(event.Y())

		if event.Button() == 1 { //左键
			x, y = int(event.X()), int(event.Y()) //保存点击的起点坐标
		} else if event.Button() == 3 { //右键
			//右键，关闭窗口
			//gtk.MainQuit()
			app.Quit()
		}
	})

	//鼠标移动事件处理
	_, _ = win.Connect("motion-notify-event", func(widget *gtk.Window, ctx *gdk.Event) {
		//获取鼠标移动属性结构体变量，系统内部的变量，不是用户传参变量
		event := *(*gdk.EventButton)(unsafe.Pointer(&ctx))
		win.Move(int(event.XRoot())-x, int(event.YRoot())-y)
	})

	go func() {
		for {
			time.Sleep(refreshPeriod)
			total, bpm, todayMax := buildShowData()

			str := fmt.Sprintf("<span foreground='green' font_desc='16'>%d</span> <span font_desc='15'>%d</span> <span foreground='yellow' font_desc='10'>%d</span>",
				bpm, total, todayMax)
			_, err := glib.IdleAdd(bpmLabel.SetMarkup, str)
			if err != nil {
				log.Fatal("IdleAdd() failed:", err)
			}
		}
	}()
}

func windowWidget() *gtk.Widget {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	bpmLabel, err = gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	grid.Add(bpmLabel)
	bpmLabel.SetHExpand(false)
	bpmLabel.SetVExpand(true)

	return &grid.Container.Widget
}

func buildShowData() (int, int, int) {
	conn := GetConnection()
	now := time.Now()
	today := now.Format(DateFormat)
	score := conn.ZScore(TotalCount, today)
	total := score.Val()

	bpmKey := GetMaxBPMKey(now)
	todayMax, err := conn.Get(bpmKey).Int()
	if err != nil {
		todayMax = 0
	}

	bpm := calculateBPM(conn, total, now)
	if todayMax < bpm && now.Second() > recordThreshold {
		conn.Set(bpmKey, bpm, 0)
		todayMax = bpm
		logger.Info("bump to", bpm)
	}

	return int(total), bpm, todayMax
}

func calculateBPM(conn *redis.Client, total float64, now time.Time) int {
	if firstStart {
		conn.Set(OddKey, total, 0)
		conn.Set(EvenKey, total, 0)
		firstStart = false
		return 0
	}

	second := now.Second()
	if second <= 1 {
		return 0
	}

	lastKey, curKey := selectActualKey(now)

	conn.Set(curKey, total, 0)
	lastTotal, err := conn.Get(lastKey).Float64()
	if err == nil {
		delta := total - lastTotal
		result := int(delta * 60 / float64(second))
		//fmt.Println(delta, "* 60 / ", second, "=", result)
		return result
	} else {
		fmt.Println(err)
	}
	return 0
}

func selectActualKey(now time.Time) (string, string) {
	lastKey := OddKey
	curKey := EvenKey
	if now.Minute()%2 == 1 {
		lastKey = EvenKey
		curKey = OddKey
	}
	return lastKey, curKey
}
