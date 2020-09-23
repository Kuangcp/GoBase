package cuibase

import "fmt"

type (
	ParamInfo struct {
		Verb    string
		Param   string
		Comment string
		Handler func(params []string)
	}
	ParamVO struct {
		Short   string
		Long    string
		Value   string
		Comment string
	}

	HelpInfo struct {
		Version       string
		Description   string
		SingleFlagLen int
		DoubleFlagLen int
		ValueLen      int
		Flags         []ParamVO
		Options       []ParamVO
		Args          []ParamVO
	}
)

// PrintHelp info msg
func (helpInfo HelpInfo) PrintHelp() {
	// Usage Description
	printTitleDefault(helpInfo.Description)

	format := BuildFormat(helpInfo)

	if len(helpInfo.Flags) > 0 {
		fmt.Printf("\n%v\n", LightCyan.Print("Flags:"))
		PrintParams(format, Green, helpInfo.Flags)
	}
	if len(helpInfo.Options) > 0 {
		fmt.Printf("\n%v\n", LightCyan.Print("Options:"))
		PrintParams(format, Purple, helpInfo.Options)
	}

	if len(helpInfo.Args) > 0 {
		fmt.Printf("\n%v\n", LightCyan.Print("Args:"))
		PrintParams(format, Green, helpInfo.Args)
	}

	if helpInfo.Version != "" {
		fmt.Printf("\n%s  %v\n\n", LightCyan.Print("Version:"), helpInfo.Version)
	}
}
