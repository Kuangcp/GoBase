package conf

import (
	"github.com/kuangcp/gobase/pkg/ctk"
	"github.com/kuangcp/logger"
	"os"
)

var (
	MainDir = "/.config/app-conf/keylogger"
	LogDir  = MainDir + "/log"
	LogPath string
	DbPath  string
)

func ConfigLogger() {
	//logger.SetLogPathTrim("/keylogger/")

	home, err := ctk.Home()
	ctk.CheckIfError(err)
	MainDir = home + MainDir

	err = os.MkdirAll(MainDir, 0755)
	ctk.CheckIfError(err)
	LogDir = home + LogDir

	DbPath = MainDir + "/db"

	err = os.MkdirAll(LogDir, 0755)
	ctk.CheckIfError(err)

	LogPath = LogDir + "/main.log"
	_ = logger.SetLoggerConfig(&logger.LogConfig{
		TimeFormat: ctk.YYYY_MM_DD_HH_MM_SS_MS,
		Console: &logger.ConsoleLogger{
			Level:    logger.DebugDesc,
			Colorful: true,
		},
		File: &logger.FileLogger{
			Filename:   LogPath,
			Level:      logger.DebugDesc,
			Colorful:   true,
			Append:     true,
			PermitMask: "0660",
		},
	})
}
