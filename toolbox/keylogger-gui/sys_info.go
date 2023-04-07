package main

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
	"github.com/shirou/gopsutil/mem"
	"log"
	"time"
)

const (
	refreshSec = 5
)

type MonitorItem struct {
	initX, initY, x, y, red, green, blue float64
	deltaFunc                            func() float64
	widget                               *gtk.DrawingArea
}

func (t *MonitorItem) delta() {
	t.x = t.deltaFunc()
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

func memoryInfo() float64 {
	memInfo, _ := mem.VirtualMemory()
	return (100 - memInfo.UsedPercent) * width / 100
}

func swapMemoryInfo() float64 {
	memInfo, _ := mem.SwapMemory()
	return -memInfo.UsedPercent * width / 100
}

func refreshDrawArea(items []*MonitorItem) {
	for range time.NewTicker(time.Second * refreshSec).C {
		for _, i := range items {
			i.delta()
		}
		win.QueueDraw()
	}
}
