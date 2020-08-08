package main

import (
	"flag"
	"github.com/kuangcp/gobase/mybook/app/common/web"
	"github.com/kuangcp/gobase/mybook/app/service"
)

var (
	updateDb      bool
	printCategory bool
	printAccount  bool
	webServer     bool
	debugStatic   bool
)

func init() {
	flag.BoolVar(&updateDb, "u", true, "create or update database table")
	flag.BoolVar(&printCategory, "pc", false, "print all category")
	flag.BoolVar(&printAccount, "pa", false, "print all account")
	flag.BoolVar(&webServer, "s", false, "start web server")
	flag.BoolVar(&debugStatic, "d", false, "debug for static file")
}

func main() {
	flag.Parse()

	if updateDb {
		service.AutoMigrateAll()
	}

	if printCategory {
		service.PrintCategory()
	}

	if printAccount {
		service.PrintAccount()
	}

	if webServer {
		web.Server(debugStatic)
	}
}
