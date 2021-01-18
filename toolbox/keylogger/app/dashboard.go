package app

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/go-redis/redis"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	width         = 140
	height        = 24
	refreshPeriod = time.Millisecond * 1688
)

func ShowWindow() {
	gtk.Init(nil)

	app, _ := gtk.ApplicationNew("com.github.kuangcp.keylogger", glib.APPLICATION_FLAGS_NONE)
	_, err := app.Connect("activate", func() {
		createWindow(app)
	})
	if err != nil {
		log.Fatal(err)
	}

	app.Run(nil)
}

func createWindow(app *gtk.Application) {
	win, _ := gtk.WindowNew(gtk.WINDOW_POPUP)
	app.AddWindow(win)
	area, _ := gtk.DrawingAreaNew()

	win.Add(area)
	win.SetTitle("Dashboard")
	win.SetDefaultSize(width, height)
	win.Connect("destroy", gtk.MainQuit)
	win.SetPosition(gtk.WIN_POS_MOUSE)
	win.ShowAll()

	var x, y int
	// 鼠标按下事件
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
		//fmt.Println(x, int(event.XRoot())-x, y,int(event.YRoot())-y)
		win.Move(int(event.XRoot())-x, int(event.YRoot())-y)
	})

	// Event handlers
	area.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		total, bpm, todayMax := buildTotalAndBPM()

		cr.MoveTo(5, 18)

		cr.SetSourceRGB(0, 150, 0)
		cr.SetFontSize(18)
		cr.ShowText(fmt.Sprintf("%3d ", bpm))

		cr.SetSourceRGB(255, 255, 255)
		cr.SetFontSize(19)
		cr.ShowText(fmt.Sprintf("%-5d ", total))

		cr.SetSourceRGB(150, 150, 0)
		cr.SetFontSize(18)
		cr.ShowText(fmt.Sprintf("%-3d ", todayMax))

		cr.Fill()
	})

	refresh(win)
}
func buildTotalAndBPM() (int, int, int) {
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
	if todayMax < bpm {
		conn.Set(bpmKey, bpm, 0)
		todayMax = bpm
	}

	return int(total), bpm, todayMax
}

func calculateBPM(conn *redis.Client, total float64, now time.Time) int {
	lastKey := OddKey
	curKey := EvenKey
	if now.Minute()%2 == 1 {
		lastKey = EvenKey
		curKey = OddKey
	}

	conn.Set(curKey, total, 0)
	lastTotal, err := conn.Get(lastKey).Float64()
	if err == nil {
		delta := total - lastTotal
		if delta <= 0 {
			return 0
		}
		return int(delta * 60 / float64(now.Second()))
	} else {
		fmt.Println(err)
	}
	return 0
}

func refresh(win *gtk.Window) {
	go func() {
		for true {
			time.Sleep(refreshPeriod)
			win.QueueDraw()
		}
	}()
}
