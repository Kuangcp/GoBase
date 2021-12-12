package main

import (
	"embed"
	"flag"
	"fmt"
	"mybook/app/account"
	"mybook/app/category"
	"mybook/app/common"
	"mybook/app/common/config"
	"mybook/app/common/web"
)

var (
	updateDb      bool
	printCategory bool
	printAccount  bool
	webServer     bool
)

//go:embed mybook-static/dist
var staticFS embed.FS

func init() {
	flag.BoolVar(&updateDb, "u", false, "create or update database table")
	flag.BoolVar(&printCategory, "pc", false, "print all category")
	flag.BoolVar(&printAccount, "pa", false, "print all account")
	flag.BoolVar(&webServer, "s", false, "start web server")

	flag.BoolVar(&config.AppConf.DebugStatic, "d", false, "debug for static file")
	flag.BoolVar(&config.AppConf.Release, "r", false, "Release logic")
	flag.IntVar(&config.AppConf.Port, "p", config.DefaultPort, "web server port")
}

func main() {
	flag.Parse()
	fmt.Println("config: ", config.AppConf)
	config.InitAppConfig()

	if updateDb {
		common.AutoMigrateAll()
	}

	if printCategory {
		category.PrintCategory()
	}

	if printAccount {
		account.PrintAccount()
	}

	if webServer {
		web.Server(staticFS)
	}
}
