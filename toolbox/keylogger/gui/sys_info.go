package main

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"math"
	"time"
)

const (
	refreshSec = 1
	lineWeight = 1.1
)

type MonitorItem struct {
	height             int
	initX, initY, x, y float64
	red, green, blue   float64
	deltaFunc          func(*MonitorItem)
}

func drawCpuItem(grid *gtk.Grid, item *MonitorItem) {
	grid.Attach(item.buildItem(), 0, 0, 1, height)

	background, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("", err)
	}
	background.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGBA(50, 50, 50, 0.6)
		cr.Rectangle(0, 0, 1, float64(item.height))
		cr.Fill()
	})
	grid.Attach(background, 0, 0, 1, item.height)
}

func drawMemItem(grid *gtk.Grid, left *MonitorItem) {
	grid.Attach(left.buildItem(), 0, 0, width*2, height)

	drawLineBackground(grid)
}

func drawMemAndSwapItem(grid *gtk.Grid, left, right *MonitorItem) {
	drawMidSeparation(grid)

	grid.Attach(left.buildItem(), 0, 0, width, height)
	grid.Attach(right.buildItem(), 0, 0, width, height)

	drawLineBackground(grid)
}

func drawLineBackground(grid *gtk.Grid) {
	background, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("", err)
	}
	background.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGBA(50, 50, 50, 0.6)
		cr.Rectangle(0, 0, lineWeight+width, lineWeight)
		cr.Fill()
	})
	grid.Attach(background, 0, 0, width, height)
}

func drawMidSeparation(grid *gtk.Grid) {
	// 黑色隔断
	mid, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("", err)
	}
	mid.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(0, 0, 0)
		//cr.SetSourceRGB(54, 106, 206)
		cr.Rectangle(width/2-0.5, 0, lineWeight+0.5, lineWeight)
		cr.Fill()
	})
	grid.Attach(mid, 0, 0, width, height)
}

// https://docs.gtk.org/gtk3/class.DrawingArea.html
func (t *MonitorItem) buildItem() *gtk.DrawingArea {
	da, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("", err)
	}

	da.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(t.red, t.green, t.blue)
		cr.Rectangle(t.initX, t.initY, lineWeight+t.x, lineWeight+t.y)
		cr.Fill()
	})
	return da
}

func memoryInfoOne(t *MonitorItem) {
	memInfo, _ := mem.VirtualMemory()
	t.x = (100 - memInfo.UsedPercent) * width / 100
}

func memoryInfo(t *MonitorItem) {
	memInfo, _ := mem.VirtualMemory()
	t.x = (100 - memInfo.UsedPercent) * width / 200
}

func cpuInfo(t *MonitorItem) {
	go func() {
		// 阻塞计算采样区间内的平均值
		memInfo, _ := cpu.Percent(time.Millisecond*100, false)
		t.y = math.Round((memInfo[0]) * float64(t.height) / 100)
		//fmt.Println(memInfo[0], t.height, t.y)
	}()
}

func swapMemoryInfo(t *MonitorItem) {
	memInfo, _ := mem.SwapMemory()
	t.initX = width/2 + (100-memInfo.UsedPercent)*width/200
	t.initX = float64(int64(t.initX))
	t.x = width - t.initX
	// fmt.Println(t.initX, t.x)
}

func refreshDrawArea(items []*MonitorItem) {
	refreshOnce(items)
	glib.TimeoutAdd(uint((time.Second * refreshSec).Milliseconds()), func() bool {
		refreshOnce(items)
		return true
	})

	//for range time.NewTicker(time.Second * refreshSec).C {
	//	refreshOnce(items)
	//}
}

func refreshOnce(items []*MonitorItem) {
	for _, i := range items {
		i.deltaFunc(i)
	}
	win.QueueDraw()
}
