package main

import (
	"github.com/kuangcp/gobase/cuibase"
	"github.com/kuangcp/gobase/mybook/app/common/web"
	"github.com/kuangcp/gobase/mybook/app/service"
)

var info = cuibase.HelpInfo{
	Description: "Simple Bookkeeping",
	VerbLen:     -3,
	ParamLen:    -58,
	Version:     "1.0.1",
	Params: []cuibase.ParamInfo{
		{
			Verb:    "-h",
			Param:   "",
			Comment: "help",
		}, {
			Verb:    "-u",
			Param:   "",
			Comment: "update database structure",
			Handler: func(_ []string) {
				service.AutoMigrateAll()
			},
		}, {
			Verb:    "-r",
			Param:   "Type AccountId CategoryId Amount Date [Comment]",
			Comment: "create record ",
			Handler: service.CreateRecordByParams,
		}, {
			Verb:    "-re",
			Param:   "AccountId CategoryId Amount Date [Comment]",
			Comment: "create expense record ",
			Handler: service.CreateExpenseRecordByParams,
		}, {
			Verb:    "-ri",
			Param:   "AccountId CategoryId Amount Date [Comment]",
			Comment: "create income record ",
			Handler: service.CreateIncomeRecordByParams,
		}, {
			Verb:    "-rt",
			Param:   "OutAccountId CategoryId Amount Date InAccountId [Comment]",
			Comment: "create transfer record ",
			Handler: service.CreateTransRecordByParams,
		}, {
			Verb:    "-pc",
			Param:   "",
			Comment: "print all category",
			Handler: service.PrintCategory,
		}, {
			Verb:    "-pa",
			Param:   "",
			Comment: "print all account",
			Handler: service.PrintAccount,
		}, {
			Verb:    "-s",
			Param:   "",
			Comment: "start web service",
			Handler: web.Server,
		},
		{
			Verb:    "-sd",
			Param:   "",
			Comment: "start debug web service",
			Handler: web.DebugServer,
		},
	}}

func main() {
	cuibase.RunActionFromInfo(info, nil)
}
