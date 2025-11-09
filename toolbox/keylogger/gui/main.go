package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gotk3/gotk3/cairo"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/toolbox/keylogger/app/store"
	"github.com/shirou/gopsutil/mem"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kuangcp/logger"
)

var (
	DashboardMsMode bool // 以 ms 格式刷新时间
	SwapMemory      bool
	DashboardMs     = DefaultRefreshMs
)

const (
	width            = 80
	height           = 10
	DefaultRefreshMs = 57
	appId            = "com.github.kuangcp.keylogger"
)

var (
	app           *gtk.Application
	win           *gtk.Window
	kpmLabel      *gtk.Label
	refreshPeriod = time.Millisecond * 400
)

var (
	help bool

	// redis
	host string
	port string
	pwd  string
	db   int

	option redis.Options
)

var (
	buildVersion string
)

var info = ctool.HelpInfo{
	Description:   "Record key input, show rank",
	Version:       "1.2.0",
	BuildVersion:  buildVersion,
	SingleFlagLen: -5,
	DoubleFlagLen: 0,
	ValueLen:      -6,
	Flags: []ctool.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help info"},
		{Short: "-m", BoolVar: &DashboardMsMode, Comment: "show time with ms"},
		{Short: "-s", BoolVar: &SwapMemory, Comment: "show swap memory rate"},
	},
	Options: []ctool.ParamVO{
		{Short: "-host", Value: "host", Comment: "redis host"},
		{Short: "-port", Value: "port", Comment: "redis port"},
		{Short: "-pwd", Value: "pwd", Comment: "redis password"},
		{Short: "-db", Value: "db", Comment: "redis db"},
		{Short: "-ms", Value: "ms", Comment: "gui refresh ms"},
	},
}

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "")
	flag.StringVar(&port, "port", "6667", "")
	flag.StringVar(&pwd, "pwd", "", "")
	flag.IntVar(&db, "db", 5, "")
	flag.IntVar(&DashboardMs, "ms", DashboardMs, "")
}

func main() {
	info.Parse()

	if help {
		info.PrintHelp()
		return
	}

	//main2()
	//go notifyAny()

	// killer
	var list = []string{"firefox"}
	go runKiller(list)
	//go oomKiller()

	option = redis.Options{Addr: host + ":" + port, Password: pwd, DB: db}
	store.InitConnection(option, false)

	if DashboardMsMode || (!DashboardMsMode && DashboardMs != DefaultRefreshMs) {
		refreshPeriod = time.Millisecond * time.Duration(DashboardMs)
	}
	logger.Info("refresh:", refreshPeriod)

	gtk.Init(nil)

	app, _ = gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	app.Connect("activate", createWindow)
	app.Run(nil)

	//createWindow()
	//gtk.Main()
}

func oomKiller() {
	type Score struct {
		pid   string
		score int
	}
	for range time.NewTicker(time.Second * 10).C {
		memInfo, _ := mem.VirtualMemory()
		//fmt.Println(memInfo.UsedPercent)
		if memInfo.UsedPercent < 95 {
			continue
		}
		var list []Score
		out, _ := execCommand("ls /proc/*/oom_score")
		pids := strings.Split(out, "\n")
		for _, pid := range pids {
			if strings.Contains(pid, "self") || pid == "" {
				continue
			}
			oomScore, success := execCommand("cat " + pid)
			if !success {
				continue
			}

			id := strings.Split(pid, "/")[2]
			//fmt.Printf("%v  %v [%v] \n", id, pid, strings.TrimSpace(oomScore))
			score, _ := strconv.Atoi(strings.TrimSpace(oomScore))
			if score == 0 {
				continue
			}
			list = append(list, Score{pid: id, score: score})
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].score >= list[j].score
		})
		if len(list) == 0 {
			continue
		}
		logger.Error("all: ", len(list), " kill PID:", list[0].pid, list[0].score)
		execCommand("kill -9" + list[0].pid)
	}

}

func runKiller(list []string) {
	for range time.NewTicker(time.Second * 5).C {
		memInfo, _ := mem.VirtualMemory()
		//fmt.Println(memInfo.UsedPercent)
		if memInfo.UsedPercent < 92 {
			continue
		}
		for _, app := range list {
			command, _ := execCommand("ps ux | awk '/" + app + "/ && !/awk/ {print $2}'")
			if command == "" {
				continue
			}
			pids := strings.Split(command, "\n")
			for _, id := range pids {
				if id == "" {
					continue
				}
				logger.Warn("kill " + id)
				execCommand("kill " + id)
			}
		}
	}
}

func createWindow() {
	grid := createMainGrid()

	win, _ = gtk.WindowNew(gtk.WINDOW_POPUP)
	win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)  // 设置为停靠窗口类型，保持在顶层
	win.SetDefaultSize(width, height)
	win.SetPosition(gtk.WIN_POS_MOUSE)
	win.Add(grid)
	win.Connect("destroy", gtk.MainQuit)
	
	// 设置窗口支持透明
	screen, _ := gdk.ScreenGetDefault()
	visual, _ := screen.GetRGBAVisual()
	if visual != nil {
		win.SetVisual(visual)
	}
	win.SetAppPaintable(true)
	
	// 设置窗口背景色
	win.Connect("draw", func(widget *gtk.Window, cr *cairo.Context) bool {
		width := float64(widget.GetAllocatedWidth())
		height := float64(widget.GetAllocatedHeight())
		
		// 设置圆角矩形路径
		radius := float64(8)  // 减小圆角半径
		degrees := math.Pi / 180.0
		
		cr.NewPath()
		cr.Arc(width-radius, radius, radius, -90*degrees, 0*degrees)
		cr.Arc(width-radius, height-radius, radius, 0*degrees, 90*degrees)
		cr.Arc(radius, height-radius, radius, 90*degrees, 180*degrees)
		cr.Arc(radius, radius, radius, 180*degrees, 270*degrees)
		cr.ClosePath()
		
		// 设置背景颜色（RGB: 43,45,51, Alpha: 0.8）
		cr.SetSourceRGBA(40.0/255.0, 40.0/255.0, 40.0/255.0, 1)
		cr.Fill()
		
		return false
	})
	
	// Load CSS stylesheet
	mRefProvider, _ := gtk.CssProviderNew()
	mRefProvider.LoadFromPath("style.css")
	
	// Apply CSS to window and its children
	styleContext, _ := win.GetStyleContext()
	styleContext.AddProvider(mRefProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	styleContext.AddClass("rounded")
	
	// Apply CSS to grid
	gridStyleContext, _ := grid.GetStyleContext()
	gridStyleContext.AddProvider(mRefProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	gridStyleContext.AddClass("rounded")
	
	// Apply CSS to label
	labelStyleContext, _ := kpmLabel.GetStyleContext()
	labelStyleContext.AddProvider(mRefProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	labelStyleContext.AddClass("rounded")
	
	bindMouseActionForWindow()
	
	win.SetOpacity(0.96)
	app.AddWindow(win)
	win.ShowAll()

	// TODO 系统信息2D绘制

	// 启动后的计算并刷新缓存
	go timeoutRefresh(refreshPeriod)
}

// https://developer.gnome.org/gtk4/unstable/GtkLabel.html
func createMainGrid() *gtk.Widget {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	grid.SetHAlign(gtk.ALIGN_FILL)
	grid.SetVAlign(gtk.ALIGN_FILL)

	kpmLabel, err = gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	kpmLabel.SetMarkup(latestLabelStr(time.Now()))
	kpmLabel.SetHExpand(true)
	kpmLabel.SetVExpand(true)
	grid.Attach(kpmLabel, 0, 0, width, height)

	// var items []*MonitorItem
	// cpuItem := &MonitorItem{initX: 0, initY: 0, red: 0, height: 3.8 * height, green: 100, blue: 50, deltaFunc: cpuInfo}
	// drawCpuItem(grid, cpuItem)
	// items = append(items, cpuItem)

	// if SwapMemory {
	// 	left := &MonitorItem{initX: 0, initY: 0, red: 0, green: 100, blue: 50, deltaFunc: memoryInfo}
	// 	right := &MonitorItem{initX: width / 2, initY: 0, red: 125, green: 0, blue: 0, deltaFunc: swapMemoryInfo}
	// 	drawMemAndSwapItem(grid, left, right)
	// 	items = append(items, left, right)
	// } else {
	// 	left := &MonitorItem{initX: 0, initY: 0, red: 0, green: 100, blue: 50, deltaFunc: memoryInfoOne}
	// 	drawMemItem(grid, left)
	// 	items = append(items, left)
	// }
	// go refreshDrawArea(items)

	return &grid.Container.Widget
}

func bindMouseActionForWindow() {
	// 鼠标最后点击坐标
	var x, y int
	win.SetEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK))

	//鼠标按下事件处理
	win.Connect("button-press-event", func(widget *gtk.Window, ctx *gdk.Event) {
		//获取鼠键按下属性结构体变量，系统内部的变量，不是用户传参变量
		event := *(*gdk.EventButton)(unsafe.Pointer(&ctx))
		if event.Button() == 1 { //左键
			x, y = int(event.X()), int(event.Y()) //保存点击的起点坐标
		} else if event.Button() == 3 { //右键
			//右键，关闭窗口
			app.Quit()
		}
	})

	// 鼠标移动事件处理
	// 注意： 不出现偏差的前提是 应用 内不出现其他点击事件和交互
	win.Connect("motion-notify-event", func(widget *gtk.Window, ctx *gdk.Event) {
		//获取鼠标移动属性结构体变量，系统内部的变量，不是用户传参变量
		event := *(*gdk.EventButton)(unsafe.Pointer(&ctx))
		win.Move(int(event.XRoot())-x, int(event.YRoot())-y)
	})
}

func latestLabelStr(now time.Time) string {
	tempValue := store.TempKPMVal(now)
	maxValue := store.MaxKPMVal(now)
	total := store.TotalCountVal(now)

	var timeFmt = ""
	if DashboardMsMode {
		timeFmt = store.MsTimeFormat
	} else {
		timeFmt = store.TimeFormat
	}

	// style https://blog.csdn.net/bitscro/article/details/3874616
	// 时间居中调整： 时间的fmt以及字体大小
	return "<span>" +
		// now time
		"<span foreground='#00FFF6' font_desc='11'>" + fmt.Sprintf("%9s", now.Format(timeFmt)) + "</span>\n" +
		// today kpm
		"<span foreground='#5AFF00' font_desc='9'>" + fmt.Sprintf("%3s", tempValue) + "</span> " +
		// today max kpm
		"<span foreground='gray' font_desc='6'>" + fmt.Sprintf("%3s", maxValue) + "</span> " +
		// today total
		"<span foreground='white' font_desc='6'>" + fmt.Sprintf("%-6d", total) + "</span>" +
		"</span>"
}

// 从缓存中更新窗口内面板
func timeoutRefresh(period time.Duration) {
	// 返回 true 才能一直执行
	glib.TimeoutAdd(uint(period.Milliseconds()), func() bool {
		str := latestLabelStr(time.Now())
		kpmLabel.SetMarkup(str)
		return true
	})
}
