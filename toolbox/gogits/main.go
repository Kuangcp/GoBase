package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
)

var (
	buildVersion string
	help         bool
	repoPath     string
	outputPath   string
)

var info = ctool.HelpInfo{
	Description:   "Analyze git repository commits and generate HTML report with ECharts",
	BuildVersion:  buildVersion,
	Version:       "1.0.0",
	SingleFlagLen: -2,
	ValueLen:      -6,
	Flags: []ctool.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help"},
	},
	Options: []ctool.ParamVO{
		{Short: "-p", StringVar: &repoPath, String: ".", Value: "path",
			Comment: "git repository path"},
		{Short: "-o", StringVar: &outputPath, String: "report.html", Value: "file",
			Comment: "output HTML report file"},
	},
}

func main() {
	info.Parse()
	if help {
		info.PrintHelp()
		return
	}

	start := time.Now()
	result, err := Analyze(repoPath)
	if err != nil {
		logger.Fatal(err)
	}
	result.GenerationDuration = time.Since(start)

	err = GenerateReport(result, outputPath)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println("Report generated:", outputPath)
	os.Exit(0)
}
