package main

import (
	"flag"
	"github.com/kuangcp/gobase/mybook/app/common/web"
	"github.com/kuangcp/gobase/mybook/app/service"
)

func main() {
	update := flag.Bool("u", true, "generate or update database structure")
	printCategory := flag.Bool("pc", false, "print all category")
	printAccount := flag.Bool("pa", false, "print all account")
	webServer := flag.Bool("s", false, "start web server")
	debug := flag.Bool("d", false, "debug with static file")

	flag.Parse()

	if *update {
		service.AutoMigrateAll()
	}
	if *printCategory {
		service.PrintCategory()
	}
	if *printAccount {
		service.PrintAccount()
	}
	if *webServer {
		web.Server(*debug)
	}
}
