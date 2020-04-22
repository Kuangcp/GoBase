package main

import "github.com/kuangcp/gobase/cuibase"

func help(_ []string) {
	info := cuibase.HelpInfo{
		Description: "Start simple http server on current path",
		VerbLen:     -5,
		ParamLen:    -5,
		Params: []cuibase.ParamInfo{
			{
				Verb:    "-h",
				Param:   "",
				Comment: "help",
			},
			{
				Verb:    "",
				Param:   "file",
				Comment: "refresh catalog",
			},
			{
				Verb:    "-a",
				Param:   "file",
				Comment: "append catalog",
			},
			{
				Verb:    "-at",
				Param:   "file",
				Comment: "append title and catalog",
			},
			{
				Verb:    "-mm",
				Param:   "file",
				Comment: "show mind map",
			},
		}}
	cuibase.Help(info)
}

func main() {
	cuibase.RunAction(map[string]func(params []string){
		"-h": help,
	}, help)
}
