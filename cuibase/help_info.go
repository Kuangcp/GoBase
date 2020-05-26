package cuibase

import "fmt"

type (
	// HelpInfo data
	HelpInfo struct {
		Version     string
		Description string
		VerbLen     int
		ParamLen    int
		Params      []ParamInfo
	}
)

// PrintHelp info msg
func (helpInfo HelpInfo) PrintHelp() {
	printTitleDefault(helpInfo.Description)
	format := BuildFormat(helpInfo.VerbLen, helpInfo.ParamLen)
	PrintParams(format, helpInfo.Params)
	if helpInfo.Version != "" {
		fmt.Printf("\n%sVersion:%s  %v\n\n", LightGreen, End, helpInfo.Version)
	}
}
