package cuibase

import (
	"testing"
)

func TestHelpInfo_PrintHelp(t *testing.T) {
	info := HelpInfo{
		Description:   "Translation between Chinese and English By Baidu API",
		Version:       "1.0.1",
		SingleFlagLen: -2,
		DoubleFlagLen: -8,
		ValueLen:      -9,
		Flags: []ParamVO{
			{Short: "-s", Value: "<value>", Comment: "use desc"},
			{Long: "--s", Value: "<value>", Comment: "use desc"},
		},
		Options: []ParamVO{{Short: "-s", Long: "--s", Value: "<value>", Comment: "use desc"}},
	}
	info.PrintHelp()
}

func TestColor(t *testing.T) {
	PrintWithColorful()
    println()
	print(Red.Print("Red"))
}
