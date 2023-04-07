package main

import (
	"fmt"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
	"github.com/shirou/gopsutil/mem"
	"log"
	"time"
)

const (
	refreshSec = 1
)

type MonitorItem struct {
	initX, initY, x, y float64
	red, green, blue   float64
	deltaFunc          func(*MonitorItem)
	widget             *gtk.DrawingArea
}

// https://docs.gtk.org/gtk3/class.DrawingArea.html
func (t *MonitorItem) buildItem() *gtk.DrawingArea {
	da, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("", err)
	}

	unitSize := 1.0
	da.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(t.red, t.green, t.blue)
		cr.Rectangle(t.initX, t.initY, unitSize+t.x, unitSize+t.y)
		cr.Fill()
	})
	t.widget = da
	return t.widget
}

func memoryInfo(t *MonitorItem) {
	memInfo, _ := mem.VirtualMemory()
	t.x = (100 - memInfo.UsedPercent) * width / 200
}

func swapMemoryInfo(t *MonitorItem) {
	memInfo, _ := mem.SwapMemory()
	t.initX = width/2 + (100-memInfo.UsedPercent)*width/200
	t.initX = float64(int64(t.initX))
	t.x = width - t.initX
	fmt.Println(t.initX, t.x)
}

func refreshDrawArea(items []*MonitorItem) {
	for range time.NewTicker(time.Second * refreshSec).C {
		for _, i := range items {
			i.deltaFunc(i)
		}
		win.QueueDraw()
	}
}
