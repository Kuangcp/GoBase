package app

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

const (
	width         = 130
	height        = 24
	refreshPeriod = time.Second * 2
)

func ShowWindow() {
	gtk.Init(nil)

	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	area, _ := gtk.DrawingAreaNew()

	win.Add(area)
	win.SetTitle("Dashboard")
	win.SetDefaultSize(width, height)
	win.Connect("destroy", gtk.MainQuit)
	win.SetPosition(gtk.WIN_POS_MOUSE)
	win.ShowAll()

	// Event handlers
	area.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(0, 150, 0)
		cr.MoveTo(2, 20)
		cr.SetFontSize(20)
		fillText(cr)
		cr.Fill()
	})

	refresh(win)

	gtk.Main()
}

func fillText(cr *cairo.Context) {
	conn := GetConnection()
	now := time.Now()
	today := now.Format(DateFormat)
	score := conn.ZScore(TotalCount, today)
	total := score.Val()

	bpm := calculateBPM(conn, total, now)

	text := fmt.Sprintf("%5d | %3d", int(total), int(bpm))
	cr.ShowText(text)
}

func calculateBPM(conn *redis.Client, total float64, now time.Time) float64 {
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
		return delta * 60 / float64(now.Second())
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
