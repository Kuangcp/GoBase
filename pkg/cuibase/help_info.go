package cuibase

import (
	"flag"
	"fmt"
)

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

		DefaultVar interface{}

		BoolType    bool
		BoolVar     *bool
		Float64Type bool
		Float64Var  *float64
		IntType     bool
		IntVar      *int
		Int64Type   bool
		Int64Var    *int64
		StringType  bool
		StringVar   *string
		UintType    bool
		UintVar     *uint
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
	printTitleDefault(helpInfo)

	format := BuildFormat(helpInfo)

	if len(helpInfo.Flags) > 0 {
		fmt.Printf("\n%v\n", Yellow.Print("Flags:"))
		PrintParams(format, Green, helpInfo.Flags)
	}
	if len(helpInfo.Options) > 0 {
		fmt.Printf("\n%v\n", Purple.Print("Options:"))
		PrintParams(format, Green, helpInfo.Options)
	}

	if len(helpInfo.Args) > 0 {
		fmt.Printf("\n%v\n", LightCyan.Print("Args:"))
		PrintParams(format, Green, helpInfo.Args)
	}

	if helpInfo.Version != "" {
		fmt.Printf("\n%s  %v\n\n", LightCyan.Print("Version:"), helpInfo.Version)
	}
}

func (helpInfo HelpInfo) Parse() {
	flag.Usage = helpInfo.PrintHelp

	for _, flagVO := range helpInfo.Flags {
		flag.BoolVar(flagVO.BoolVar, flagVO.Short[1:], false, "")
	}
	flag.Parse()
}
