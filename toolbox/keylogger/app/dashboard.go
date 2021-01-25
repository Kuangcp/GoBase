package app

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/kuangcp/gobase/pkg/cuibase"

	"github.com/go-redis/redis"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kuangcp/logger"
)

const (
	width              = 100
	height             = 20
	refreshPeriod      = time.Millisecond * 512
	recordBPMThreshold = 59 // 当前秒数 > xx(s) 才存储 bpm
	appId              = "com.github.kuangcp.keylogger"
)

var (
	app        *gtk.Application
	win        *gtk.Window
	bpmLabel   *gtk.Label
	firstStart = true
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

	go func() {
		ticker := time.NewTicker(refreshPeriod)
		for range ticker.C {
			now := time.Now()
			refreshLabel(now)
		}
	}()

	app.AddWindow(win)
	win.ShowAll()
}

func buildWidget() *gtk.Widget {
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
	bpmLabel.SetHExpand(true)
	bpmLabel.SetVExpand(true)

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

func refreshLabel(now time.Time) {
	total, bpm, todayMax := buildShowData(now)

	// https://blog.csdn.net/bitscro/article/details/3874616
	str := fmt.Sprintf(" 🕒 %s\n%s %s %s",
		fmt.Sprintf("<span foreground='#F2F3F5' font_desc='10'>%s</span>", now.Format(TimeFormat)),
		fmt.Sprintf("<span foreground='#5AFF00' font_desc='14'>%d</span>", bpm),
		fmt.Sprintf("<span foreground='#F2F3F5' font_desc='12'>%d</span>", total),
		fmt.Sprintf("<span foreground='yellow' font_desc='9'>%d</span>", todayMax),
	)
	_, err := glib.IdleAdd(bpmLabel.SetMarkup, str)
	if err != nil {
		log.Fatal("IdleAdd() failed:", err)
	}
}

func buildShowData(now time.Time) (int, int, int) {
	conn := GetConnection()
	today := now.Format(DateFormat)
	total := conn.ZScore(TotalCount, today).Val()

	bpm := calculateBPM(conn, total, now)

	maxBPMKey := GetTodayMaxBPMKey(now)
	todayMax, err := conn.Get(maxBPMKey).Int()
	if err != nil {
		todayMax = 0
	}

	if now.Second() < recordBPMThreshold {
		return int(total), bpm, todayMax
	}

	if todayMax < bpm {
		conn.Set(maxBPMKey, bpm, 0)
		todayMax = bpm
		logger.Info("Today max BPM up to", bpm)
	}
	return int(total), bpm, todayMax
}

func calculateBPM(conn *redis.Client, total float64, now time.Time) int {
	if firstStart {
		firstStart = false
		return coverOldValue(conn, total)
	}

	second := now.Second()
	if second <= 1 {
		return 0
	}

	lastKey, curKey := selectActualKey(now)
	conn.Set(curKey, total, 0)
	lastTotal, err := conn.Get(lastKey).Float64()

	if err != nil {
		fmt.Println(err)
		return 0
	}

	// everyDay first min
	if lastTotal > total {
		return coverOldValue(conn, total)
	}

	delta := total - lastTotal
	result := int(delta * 60 / float64(now.Second()))
	//fmt.Println(delta, "* 60 / ", second, "=", result)
	return result
}

func coverOldValue(conn *redis.Client, total float64) int {
	conn.Set(OddKey, total, 0)
	conn.Set(EvenKey, total, 0)
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
