package main

import (
	"flag"
	"fmt"
	"keylogger/app"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/kuangcp/logger"

	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/pkg/cuibase"
)

var user = cuibase.Red.Print("root")
var info = cuibase.HelpInfo{
	Description:   "Record key input, show rank",
	Version:       "1.0.7",
	SingleFlagLen: -5,
	DoubleFlagLen: 0,
	ValueLen:      -14,
	Flags: []cuibase.ParamVO{
		{Short: "-h", Comment: "help info"},
		{Short: "-l", Comment: user + " list keyboard device"},
		{Short: "-L", Comment: user + " list all device"},
		{Short: "-p", Comment: user + " print key map"},
		{Short: "-c", Comment: user + " cache key map"},
		{Short: "-s", Comment: user + " listen keyboard with last device or specific device"},
		{Short: "-i", Comment: user + " listen keyboard with interactive select device"},
		{Short: "-T", Comment: "print daily total by before x day ago and duration"},
		{Short: "-R", Comment: "print daily rank by before x day ago and duration"},
		{Short: "-r", Comment: "print total rank by before x day ago and duration"},
		{Short: "-S", Comment: "web server"},
		{Short: "-b", Comment: "open small window to show total and KPM(Keystrokes Per Minute)"},
		{Short: "-d", Comment: "debug"},
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

	debug   bool
	option  redis.Options
	logPath string
)

var (
	mainDir = "/.config/app-conf/keylogger"
)

func init() {
	configLogger()

	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&printKeyMap, "p", false, "")
	flag.StringVar(&targetDevice, "e", "", "")
	flag.BoolVar(&cacheKeyMap, "c", false, "")
	flag.BoolVar(&listKeyboardDevice, "l", false, "")
	flag.BoolVar(&listAllDevice, "L", false, "")
	flag.BoolVar(&listenDevice, "s", false, "")
	flag.BoolVar(&interactiveListen, "i", false, "")
	flag.BoolVar(&printDay, "T", false, "")
	flag.BoolVar(&printDayRank, "R", false, "")
	flag.BoolVar(&printTotalRank, "r", false, "")
	flag.BoolVar(&dashboard, "b", false, "")
	flag.BoolVar(&showLog, "g", false, "")

	flag.StringVar(&timePair, "t", "1", "")

	flag.StringVar(&host, "host", "127.0.0.1", "")
	flag.StringVar(&port, "port", "6667", "")
	flag.StringVar(&pwd, "pwd", "", "")
	flag.IntVar(&db, "db", 5, "")

	flag.StringVar(&webPort, "P", "9902", "")
	flag.BoolVar(&webServer, "S", false, "")
	flag.BoolVar(&debug, "d", false, "")

	flag.Usage = info.PrintHelp
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
		TimeFormat: "2006-01-02 15:04:05.000",
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
	flag.Parse()

	option = redis.Options{
		PoolSize: 5,
		Addr:     host + ":" + port,
		Password: pwd,
		DB:       db,
	}

	pprofDebug()

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
	app.InitConnection(option)
	defer app.CloseConnection()

	invokeThenExit(dashboard, app.ShowPopWindow, app.CloseConnection)
	invokeThenExit(listenDevice, app.ListenDevice, app.CloseConnection)
	invokeThenExit(cacheKeyMap, app.CacheKeyMap, app.CloseConnection)

	if interactiveListen {
		device, err := app.SelectDevice()
		if err != nil {
			return
		}
		app.SetFormatTargetDevice(device)
		app.ListenDevice()
		return
	}

	if webServer {
		app.Server(debug, webPort)
		return
	}

	invoke(printKeyMap, app.PrintKeyMap)
	invoke(printDay, app.PrintDay)
	invoke(printDayRank, app.PrintDayRank)
	invoke(printTotalRank, app.PrintTotalRank)
}
