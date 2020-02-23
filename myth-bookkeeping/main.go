package main

import (
	"github.com/kuangcp/gobase/cuibase"
	"github.com/kuangcp/gobase/myth-bookkeeping/service"
)

func help(params []string) {
	info := cuibase.HelpInfo{
		Description: "Myth Bookkeeping",
		VerbLen:     -4,
		ParamLen:    -35,
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
				Param:   "AccountId CategoryId Type Amount ",
				Comment: "create record ",
			}, {
				Verb:    "-re",
				Param:   "AccountId CategoryId Amount ",
				Comment: "create expense record ",
			}, {
				Verb:    "-rt",
				Param:   "AccountId CategoryId Amount ",
				Comment: "create transfer record ",
			},
		}}
	cuibase.Help(info)
}

func updateDatabaseStructure(params []string) {
	// 建立数据库结构
	service.AutoMigrateAll()
}

func main() {
	cuibase.RunAction(map[string]func(params []string){
		"-h":  help,
		"-u":  updateDatabaseStructure,
		"-r":  service.CreateRecordByParams,
		"-re": service.CreateExpenseRecordByParams,
		"-rt": service.CreateTransRecordByParams,
	}, help)

}
