package main

import (
	"embed"
	"flag"
	"fmt"
	"keylogger/app"
	"keylogger/app/store"
	"keylogger/app/web"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/kuangcp/logger"

	"github.com/webview/webview"

	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/pkg/cuibase"
)

//go:embed static
var fs embed.FS

var user = cuibase.Red.Print("root")
var info = cuibase.HelpInfo{
	Description:   "Record key input, show rank",
	Version:       "1.0.9",
	BuildVersion:  buildVersion,
	SingleFlagLen: -5,
	DoubleFlagLen: 0,
	ValueLen:      -14,
	Flags: []cuibase.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help info"},
		{Short: "-l", BoolVar: &listKeyboardDevice, Comment: user + " list keyboard device"},
		{Short: "-L", BoolVar: &listAllDevice, Comment: user + " list all device"},
		{Short: "-p", BoolVar: &printKeyMap, Comment: user + " print key map"},
		{Short: "-c", BoolVar: &cacheKeyMap, Comment: user + " cache key map"},
		{Short: "-s", BoolVar: &listenDevice, Comment: user + " listen keyboard with last device or specific device"},
		{Short: "-i", BoolVar: &interactiveListen, Comment: user + " listen keyboard with interactive select device"},
		{Short: "-T", BoolVar: &printDay, Comment: "print daily total by before x day ago and duration"},
		{Short: "-R", BoolVar: &printDayRank, Comment: "print daily rank by before x day ago and duration"},
		{Short: "-r", BoolVar: &printTotalRank, Comment: "print total rank by before x day ago and duration"},
		{Short: "-S", BoolVar: &webServer, Comment: "web server"},
		{Short: "-b", BoolVar: &dashboard, Comment: "open small window to show total and KPM(Keystrokes Per Minute)"},
		{Short: "-d", BoolVar: &debug, Comment: "debug: logic and static file"},
		{Short: "-O", BoolVar: &notOpenPage, Comment: "not open url by browser"},
		{Short: "-g", BoolVar: &showLog, Comment: "show log"},
	},
	Options: []cuibase.ParamVO{
		{Short: "-t", Value: "x,duration", Comment: "before x day ago and duration. For -T and -R"},
		{Short: "-e", Value: "device", Comment: "operation target device. For -p -ca -s"},
		{Short: "-P", Value: "port", Comment: "web Server port. default 9902"},
		{Short: "-host", Value: "host", Comment: "redis host"},
		{Short: "-port", Value: "port", Comment: "redis port"},
		{Short: "-pwd", Value: "pwd", Comment: "redis password"},
		{Short: "-db", Value: "db", Comment: "redis db"},
	},
}

var (
	help               bool
	printKeyMap        bool
	cacheKeyMap        bool
	listKeyboardDevice bool
	listAllDevice      bool
	listenDevice       bool
	interactiveListen  bool
	printDay           bool
	printDayRank       bool
	printTotalRank     bool
	dashboard          bool
	showLog            bool

	targetDevice string
	timePair     string

	// redis
	host string
	port string
	pwd  string
	db   int

	// web
	webPort   string
	webServer bool
	webView   bool

	debug       bool
	notOpenPage bool
	option      redis.Options
	logPath     string
)
var (
	buildVersion string
)
var (
	mainDir = "/.config/app-conf/keylogger"
)

func init() {
	configLogger()

	flag.BoolVar(&webView, "v", false, "start webview")
	flag.BoolVar(&help, "help", false, "")
	flag.StringVar(&timePair, "t", "1", "")
	flag.StringVar(&targetDevice, "e", "", "")

	flag.StringVar(&host, "host", "127.0.0.1", "")
	flag.StringVar(&port, "port", "6667", "")
	flag.StringVar(&pwd, "pwd", "", "")
	flag.IntVar(&db, "db", 5, "")

	flag.StringVar(&webPort, "P", "9902", "")
}

func configLogger() {
	logger.SetLogPathTrim("/keylogger/")

	home, err := cuibase.Home()
	cuibase.CheckIfError(err)
	mainDir = home + mainDir

	err = os.MkdirAll(mainDir, 0755)
	cuibase.CheckIfError(err)
	logDir := mainDir + "/log"

	err = os.MkdirAll(logDir, 0755)
	cuibase.CheckIfError(err)

	logPath = logDir + "/main.log"
	_ = logger.SetLoggerConfig(&logger.LogConfig{
		TimeFormat: cuibase.YYYY_MM_DD_HH_MM_SS_MS,
		Console: &logger.ConsoleLogger{
			Level:    logger.DebugDesc,
			Colorful: true,
		},
		File: &logger.FileLogger{
			Filename:   logPath,
			Level:      logger.DebugDesc,
			Colorful:   true,
			Append:     true,
			PermitMask: "0660",
		},
	})
}

func pprofDebug() {
	if !debug {
		return
	}

	debugPort := "8891"
	go func() {
		fmt.Println("http://127.0.0.1:" + debugPort + "/debug/pprof/")
		_ = http.ListenAndServe("0.0.0.0:"+debugPort, nil)
	}()
}

func invokeThenExit(condition bool, action func(), clean func()) {
	if !condition {
		return
	}
	action()
	if clean != nil {
		clean()
	}

	os.Exit(0)
}

func invoke(condition bool, action func()) {
	if condition {
		action()
	}
}

func main() {
	info.Parse()

	pprofDebug()
	option = redis.Options{Addr: host + ":" + port, Password: pwd, DB: db}
	if showLog {
		fmt.Println(logPath)
		return
	}

	invokeThenExit(help, info.PrintHelp, nil)
	invokeThenExit(listKeyboardDevice, app.ListAllKeyBoardDevice, nil)
	invokeThenExit(listAllDevice, app.ListAllDevice, nil)

	// 以下逻辑都依赖Redis
	app.SetFormatTargetDevice(targetDevice)
	app.SetTimePair(timePair)
	store.InitConnection(option)
	defer store.CloseConnection()

	invokeThenExit(dashboard, app.ShowPopWindow, store.CloseConnection)
	invokeThenExit(listenDevice, app.ListenDevice, store.CloseConnection)
	invokeThenExit(cacheKeyMap, app.CacheKeyMap, store.CloseConnection)

	if interactiveListen {
		device, err := app.SelectDevice()
		if err != nil {
			return
		}
		app.SetFormatTargetDevice(device)
		app.ListenDevice()
		return
	}

	if webServer && !webView {
		web.Server(fs, debug, notOpenPage, webPort)
		return
	}
	if webServer && webView {
		go web.Server(fs, debug, notOpenPage, webPort)
		mainWin()
		return
	}

	invoke(printKeyMap, app.PrintKeyMap)
	invoke(printDay, app.PrintDay)
	invoke(printDayRank, app.PrintDayRank)
	invoke(printTotalRank, app.PrintTotalRank)
}

func mainWin() {
	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("Keylogger webview")
	w.SetSize(1800, 960, webview.HintNone)
	w.Navigate("http://localhost:" + webPort)
	w.Run()
}
