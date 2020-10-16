package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"

	"github.com/wonderivan/logger"

	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/cuibase"
	"github.com/kuangcp/gobase/keylogger/app"
)

var user = cuibase.Red.Print("root")
var info = cuibase.HelpInfo{
	Description:   "Record key input, show rank",
	Version:       "1.0.4",
	SingleFlagLen: -5,
	DoubleFlagLen: -10,
	ValueLen:      -14,
	Flags: []cuibase.ParamVO{
		{Short: "-h", Long: "--help", Comment: "Help info"},
		{Short: "-l", Comment: user + " List keyboard device"},
		{Short: "-ld", Comment: user + " List all device"},
		{Short: "-p", Comment: user + " Print key map"},
		{Short: "-ca", Comment: user + " Cache key map"},
		{Short: "-s", Comment: user + " Listen keyboard with last device or specific device"},
		{Short: "-d", Comment: "Print daily total by before x day ago and duration"},
		{Short: "-dr", Comment: "Print daily rank by before x day ago and duration"},
		{Short: "-ws", Comment: "Web server"},
	},
	Options: []cuibase.ParamVO{
		{Short: "-t", Long: "--time", Value: "<x>,<duration>", Comment: "Before x day ago and duration. For -d and -dr"},
		{Short: "-e", Long: "--device", Value: "<device>", Comment: "Operation target device, work for -p -ca -s"},
		{Short: "-wp", Value: "<port>", Comment: "Web Server port. default 9902"},
		{Short: "-host", Value: "<host>", Comment: "Redis host"},
		{Short: "-port", Value: "<port>", Comment: "Redis port"},
		{Short: "-pwd", Value: "<pwd>", Comment: "Redis password"},
		{Short: "-db", Value: "<db>", Comment: "Redis db"},
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
	flag.StringVar(&targetDevice, "e", "", "specific device")
	flag.BoolVar(&cacheKeyMap, "ca", false, "")
	flag.BoolVar(&listKeyboardDevice, "l", false, "")
	flag.BoolVar(&listAllDevice, "la", false, "")
	flag.BoolVar(&listenDevice, "s", false, "")
	flag.BoolVar(&day, "d", false, "")
	flag.BoolVar(&dayRank, "dr", false, "")
	flag.StringVar(&timePair, "t", "1", "")

	flag.StringVar(&host, "host", "127.0.0.1", "")
	flag.StringVar(&port, "port", "6667", "")
	flag.StringVar(&pwd, "pwd", "", "")
	flag.IntVar(&db, "db", 5, "")

	flag.StringVar(&webPort, "wp", "9902", "")
	flag.BoolVar(&webServer, "ws", false, "")
	flag.BoolVar(&debug, "debug", false, "")
}

func main() {
	flag.Parse()

	options := redis.Options{
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

	if help {
		info.PrintHelp()
		return
	} else if webServer {
		app.Server(debug, webPort)
		return
	}

	targetDevice = app.FormatEvent(targetDevice)

	if listKeyboardDevice {
		app.ListAllKeyBoardDevice()
		return
	}

	if listAllDevice {
		app.ListAllDevice()
		return
	}

	if cacheKeyMap {
		app.CacheKeyMap(targetDevice)
		return
	}

	if listenDevice {
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
}
