package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/cuibase"
	"github.com/kuangcp/gobase/keyboard-man/app"
)

var info = cuibase.HelpInfo{
	Description: "Record key input, show rank",
	Version:     "1.0.2",
	VerbLen:     -5,
	ParamLen:    -14,
	Params: []cuibase.ParamInfo{
		{
			Verb:    "-h",
			Param:   "",
			Comment: "Help info",
		}, {
			Verb:    "-l",
			Param:   "",
			Comment: cuibase.Red.Print("root") + " List keyboard device",
		}, {
			Verb:    "-ld",
			Param:   "",
			Comment: cuibase.Red.Print("root") + " List all device",
		}, {
			Verb:    "-p",
			Param:   "",
			Comment: cuibase.Red.Print("root") + " Print key map",
		}, {
			Verb:    "-ca",
			Param:   "",
			Comment: cuibase.Red.Print("root") + " Cache key map",
		}, {
			Verb:    "-s",
			Param:   "",
			Comment: cuibase.Red.Print("root") + " Listen keyboard with last device or specific device",
		}, {
			Verb:    "-d",
			Param:   "",
			Comment: "Print daily total by before x day ago and duration",
		}, {
			Verb:    "-dr",
			Param:   "",
			Comment: "Print daily rank by before x day ago and duration",
		}, {
			Verb:    "-t",
			Param:   "<x>,<duration>",
			Comment: "Before x day ago and duration. For -d and -dr",
		}, {
			Verb:    "-e",
			Param:   "device",
			Comment: "Device. For -p -ca -s",
		}, {
			Verb:    "-host",
			Param:   "",
			Comment: "Redis host",
		}, {
			Verb:    "-port",
			Param:   "",
			Comment: "Redis port",
		}, {
			Verb:    "-pwd",
			Param:   "",
			Comment: "Redis pwd",
		}, {
			Verb:    "-db",
			Param:   "",
			Comment: "Redis db",
		}, {
			Verb:    "-ws",
			Param:   "",
			Comment: "Start Web Server",
		}, {
			Verb:    "-wp",
			Param:   "port",
			Comment: "Web Server port. default 9902",
		},
	}}

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
	flag.BoolVar(&help, "h", false, "")
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
			http.ListenAndServe("0.0.0.0:"+debugPort, nil)
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
