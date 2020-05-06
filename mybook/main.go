package main

import (
	"github.com/kuangcp/gobase/cuibase"
	"github.com/kuangcp/gobase/mybook/app/cui"
	"github.com/kuangcp/gobase/mybook/app/service"
	"github.com/kuangcp/gobase/mybook/app/web"
)

func help(_ []string) {
	info := cuibase.HelpInfo{
		Description: "Myth Bookkeeping",
		VerbLen:     -4,
		ParamLen:    -60,
		Params: []cuibase.ParamInfo{
			{
				Verb:    "-h",
				Param:   "",
				Comment: "help",
			}, {
				Verb:    "-u",
				Param:   "",
				Comment: "update database structure",
			}, {
				Verb:    "-r",
				Param:   "Type AccountId CategoryId Amount Date [Comment]",
				Comment: "create record ",
			}, {
				Verb:    "-re",
				Param:   "AccountId CategoryId Amount Date [Comment]",
				Comment: "create expense record ",
			}, {
				Verb:    "-ri",
				Param:   "AccountId CategoryId Amount Date [Comment]",
				Comment: "create income record ",
			}, {
				Verb:    "-rt",
				Param:   "OutAccountId CategoryId Amount Date InAccountId [Comment]",
				Comment: "create transfer record ",
			}, {
				Verb:    "-pc",
				Param:   "",
				Comment: "print all category",
			},
		}}
	info.PrintHelp()
}

func updateDatabaseStructure(_ []string) {
	service.AutoMigrateAll()
}

func main() {
	cuibase.RunAction(map[string]func(params []string){
		"-h":  help,
		"-u":  updateDatabaseStructure,
		"-r":  service.CreateRecordByParams,
		"-re": service.CreateExpenseRecordByParams,
		"-ri": service.CreateIncomeRecordByParams,
		"-rt": service.CreateTransRecordByParams,
		"-pc": service.PrintCategory,
		"-pa": service.PrintAccount,
		"-s":  web.Server,
		"-c":  cui.CUIMain,
	}, help)
}
