package main

import (
	"github.com/kuangcp/gobase/cuibase"
	"github.com/wonderivan/logger"
)

func HelpInfo(_ []string) {
	info := cuibase.HelpInfo{
		Description: "Git repository manager",
		VerbLen:     -3,
		ParamLen:    -5,
		Params: []cuibase.ParamInfo{
			{
				Verb:    "-h",
				Param:   "",
				Comment: "Help info",
			},
		}}
	cuibase.Help(info)
}

func main() {
	logger.SetLogPathTrim("/toolbox/")
	cuibase.RunAction(map[string]func(params []string){
		"-h": HelpInfo,
	}, HelpInfo)
}
