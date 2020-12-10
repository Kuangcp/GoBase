package main

import (
	"flag"
	"mybook/app/common"
	"mybook/app/common/web"
	"mybook/app/service"
)

var (
	updateDb      bool
	printCategory bool
	printAccount  bool
	webServer     bool
	debugStatic   bool
	port          int
)

func init() {
	flag.BoolVar(&updateDb, "u", true, "create or update database table")
	flag.BoolVar(&printCategory, "pc", false, "print all category")
	flag.BoolVar(&printAccount, "pa", false, "print all account")
	flag.BoolVar(&webServer, "s", false, "start web server")
	flag.BoolVar(&debugStatic, "d", false, "debug for static file")
	flag.IntVar(&port, "p", 0, "web server port")

	flag.Parse()
}

func main() {
	if updateDb {
		common.AutoMigrateAll()
	}

	if printCategory {
		service.PrintCategory()
	}

	if printAccount {
		service.PrintAccount()
	}

	if webServer {
		web.Server(debugStatic, port)
	}
}
