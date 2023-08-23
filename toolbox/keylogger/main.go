package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/kuangcp/gobase/toolbox/keylogger/app/conf"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/kuangcp/gobase/pkg/ctk"
	"github.com/kuangcp/gobase/toolbox/keylogger/app"
	"github.com/kuangcp/gobase/toolbox/keylogger/app/store"
	"github.com/kuangcp/gobase/toolbox/keylogger/app/web"

	"github.com/go-redis/redis"
)

//go:embed static
var fs embed.FS

var user = ctk.Red.Print("root")
var redisStr = ctk.Cyan.Print("redis")
var info = ctk.HelpInfo{
	Description:   "Record key input, show rank",
	Version:       "1.1.0",
	BuildVersion:  buildVersion,
	SingleFlagLen: -5,
	DoubleFlagLen: 0,
	ValueLen:      -14,
	Flags: []ctk.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help info"},
		{Short: "-l", BoolVar: &listKeyboardDevice, Comment: user + " list keyboard device"},
		{Short: "-L", BoolVar: &listAllDevice, Comment: user + " list all device"},
		{Short: "-p", BoolVar: &printKeyMap, Comment: user + " print key map"},
		{Short: "-c", BoolVar: &cacheKeyMap, Comment: user + " cache key map"},
		{Short: "-s", BoolVar: &listenDevice, Comment: user + " listen keyboard with last device or specific device"},
		{Short: "-i", BoolVar: &interactiveListen, Comment: user + " listen keyboard with interactive select device\n"},
		{Short: "-dt", BoolVar: &printDay, Comment: "print daily total by before x day ago and duration"},
		{Short: "-dr", BoolVar: &printDayRank, Comment: "print daily rank by before x day ago and duration"},
		{Short: "-tr", BoolVar: &printTotalRank, Comment: "print total rank by before x day ago and duration\n"},
		{Short: "-S", BoolVar: &webServer, Comment: "web server"},
		{Short: "-d", BoolVar: &debug, Comment: "debug: logic and static file(must run on root dir)"},
		{Short: "-O", BoolVar: &notOpenPage, Comment: "not auto open web page by browser"},
		{Short: "-g", BoolVar: &showLog, Comment: "show log"},
	},
	Options: []ctk.ParamVO{
		{Short: "-t", Value: "x,duration", Comment: "before x day ago and duration. Provide to " + ctk.Green.Print("-T") + " " + ctk.Green.Print("-R") + " for use"},
		{Short: "-e", Value: "device", Comment: "operation target device. Provide to " + ctk.Green.Print("-p -ca -s") + " for use"},
		{Short: "-P", Value: "port", Comment: "web Server port. default 9902"},
		{Short: "-host", Value: "host", Comment: redisStr + " host"},
		{Short: "-port", Value: "port", Comment: redisStr + " port"},
		{Short: "-pwd", Value: "pwd", Comment: redisStr + " password"},
		{Short: "-db", Value: "db", Comment: redisStr + " db"},
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

	debug       bool
	notOpenPage bool
	option      redis.Options
)
var (
	buildVersion string
)

func init() {
	conf.ConfigLogger()

	flag.BoolVar(&help, "help", false, "")
	flag.StringVar(&timePair, "t", "1", "")
	flag.StringVar(&targetDevice, "e", "", "")

	flag.StringVar(&host, "host", "127.0.0.1", "")
	flag.StringVar(&port, "port", "6667", "")
	flag.StringVar(&pwd, "pwd", "", "")
	flag.IntVar(&db, "db", 5, "")

	flag.StringVar(&webPort, "P", "9902", "")
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
		fmt.Println(conf.LogPath)
		return
	}

	invokeThenExit(help, info.PrintHelp, nil)
	invokeThenExit(listKeyboardDevice, app.ListAllKeyBoardDevice, nil)
	invokeThenExit(listAllDevice, app.ListAllDevice, nil)

	// 以下逻辑都依赖Redis
	app.SetFormatTargetDevice(targetDevice)
	app.SetTimePair(timePair)
	store.InitConnection(option, true)
	defer store.CloseConnection()

	store.InitDb()
	go web.ScheduleSyncAllDetails()

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

	if webServer {
		web.Server(fs, debug, notOpenPage, webPort)
		return
	}

	invoke(printKeyMap, app.PrintKeyMap)
	invoke(printDay, app.PrintDay)
	invoke(printDayRank, app.PrintDayRank)
	invoke(printTotalRank, app.PrintTotalRank)
}
