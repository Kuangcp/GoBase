package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"

	"github.com/wonderivan/logger"

	"github.com/go-redis/redis"
	"keylogger/app"
	"github.com/kuangcp/gobase/pkg/cuibase"
)

var user = cuibase.Red.Print("root")
var info = cuibase.HelpInfo{
	Description:   "Record key input, show rank",
	Version:       "1.0.5",
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
		{Short: "-T", Comment: "print daily total by before x day ago and duration"},
		{Short: "-R", Comment: "print daily rank by before x day ago and duration"},
		{Short: "-r", Comment: "print total rank by before x day ago and duration"},
		{Short: "-S", Comment: "web server"},
		{Short: "-d", Comment: "debug"},
	},
	Options: []cuibase.ParamVO{
		{Short: "-t", Value: "x,duration", Comment: "before x day ago and duration. For -T and -R"},
		{Short: "-e", Value: "device", Comment: "operation target device, work for -p -ca -s"},
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
	day                bool
	dayRank            bool
	totalRank          bool

	targetDevice string
	timePair     string

	// redis
	host string
	port string
	pwd  string
	db   int

	webPort   string
	webServer bool

	debug bool
)

func init() {
	logger.SetLogPathTrim("/keylogger/")

	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&printKeyMap, "p", false, "")
	flag.StringVar(&targetDevice, "e", "", "")
	flag.BoolVar(&cacheKeyMap, "c", false, "")
	flag.BoolVar(&listKeyboardDevice, "l", false, "")
	flag.BoolVar(&listAllDevice, "L", false, "")
	flag.BoolVar(&listenDevice, "s", false, "")
	flag.BoolVar(&day, "T", false, "")
	flag.BoolVar(&dayRank, "R", false, "")
	flag.BoolVar(&totalRank, "r", false, "")

	flag.StringVar(&timePair, "t", "1", "")

	flag.StringVar(&host, "host", "127.0.0.1", "")
	flag.StringVar(&port, "port", "6667", "")
	flag.StringVar(&pwd, "pwd", "", "")
	flag.IntVar(&db, "db", 5, "")

	flag.StringVar(&webPort, "P", "9902", "")
	flag.BoolVar(&webServer, "S", false, "")
	flag.BoolVar(&debug, "d", false, "")

	flag.Parse()
}

func main() {
	options := redis.Options{
		PoolSize: 20,
		Addr:     host + ":" + port,
		Password: pwd,
		DB:       db,
	}
	app.InitConnection(options)
	defer app.CloseConnection()

	debugPort := "8891"
	if debug {
		go func() {
			_ = http.ListenAndServe("0.0.0.0:"+debugPort, nil)
		}()
	}

	targetDevice = app.FormatEvent(targetDevice)

	if help {
		info.PrintHelp()
		return
	} else if webServer {
		app.Server(debug, webPort)
		return
	} else if listKeyboardDevice {
		app.ListAllKeyBoardDevice()
		return
	} else if listAllDevice {
		app.ListAllDevice()
		return
	} else if cacheKeyMap {
		app.CacheKeyMap(targetDevice)
		return
	} else if listenDevice {
		app.ListenDevice(targetDevice)
		return
	}

	// simple query info

	if printKeyMap {
		app.PrintKeyMap(targetDevice)
	}

	if day {
		app.PrintDay(timePair)
	}

	if dayRank {
		app.PrintDayRank(timePair)
	}

	if totalRank {
		app.PrintTotalRank(timePair)
	}
}
