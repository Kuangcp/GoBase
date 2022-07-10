package ctool

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		BuildVersion  string
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
		fmt.Printf("\n%s  %v", LightCyan.Print("Version:"), helpInfo.Version)
	}
	if helpInfo.BuildVersion != "" {
		fmt.Printf("\n%s  %v", LightCyan.Print("Build:"), helpInfo.BuildVersion)
	}
	fmt.Println()
}

func (helpInfo HelpInfo) Parse() {
	flag.Usage = helpInfo.PrintHelp

	for _, flagVO := range helpInfo.Flags {
		flag.BoolVar(flagVO.BoolVar, flagVO.Short[1:], false, "")
	}
	flag.Parse()
}

// BuildFormat
func BuildFormat(info HelpInfo) string {
	single := strconv.Itoa(info.SingleFlagLen)
	double := strconv.Itoa(info.DoubleFlagLen)
	value := strconv.Itoa(info.ValueLen)
	return "    %v %" + single + "v%" + double + "v %v %" + value + "v %v %v\n"
}

// PrintParams
func PrintParams(format string, flagColor Color, params []ParamVO) {
	for _, vo := range params {
		if vo.Long == "" {
			fmt.Printf(format, flagColor, vo.Short, "", Yellow, vo.Value, End, vo.Comment)
		} else {
			fmt.Printf(format, flagColor, vo.Short, ", "+vo.Long, Yellow, vo.Value, End, vo.Comment)
		}
	}
}

// PrintTitle
func PrintTitle(command string, helpInfo HelpInfo) {
	flagStr := ""
	for _, flagVO := range helpInfo.Flags {
		flagStr += flagVO.Short
	}
	flagStr = strings.Replace(flagStr, "-", "", -1)

	optionStr := ""
	for _, option := range helpInfo.Options {
		optionStr += fmt.Sprintf("[%s %s] ", option.Short, option.Value)
	}
	fmt.Printf("%s\n\n  %v %v %v\n\n",
		LightCyan.Print("Usage:"),
		command,
		Yellow.PrintNoEnd("[-"+flagStr+"]"),
		Purple.Print(optionStr))

	fmt.Printf("%s\n\n  %v\n", LightCyan.Print("Description:"), helpInfo.Description)
}

func printTitleDefault(helpInfo HelpInfo) {
	PrintTitle(os.Args[0], helpInfo)
}
