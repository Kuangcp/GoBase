package cuibase

import (
	"math/rand"
	"testing"
	"time"
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

func TestProcessBar(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i <= 100; i++ {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(300)))
		DrawProgressBar("task name", float32(i)/100.0, 40)
	}
}
