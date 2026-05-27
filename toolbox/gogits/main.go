package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
)

var (
	buildVersion string
	help         bool
	openInChrome bool
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
		{Short: "-o", BoolVar: &openInChrome, Comment: "open report in Chrome after generation"},
	},
	Options: []ctool.ParamVO{
		{Short: "-p", StringVar: &repoPath, String: ".", Value: "path",
			Comment: "git repository path"},
		{Short: "-f", StringVar: &outputPath, String: "", Value: "file",
			Comment: "output HTML report file"},
	},
}

func main() {
	info.Parse()
	if outputPath == "" {
		if _, err := os.Stat(filepath.Join(repoPath, ".git")); err == nil {
			outputPath = filepath.Join(repoPath, ".git", "gogits-report.html")
		} else {
			outputPath = "gogits-report.html"
		}
	}
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

	if openInChrome {
		absPath, err := filepath.Abs(outputPath)
		if err != nil {
			logger.Fatal(err)
		}
		chromeURL := "file://" + absPath
		cmd := exec.Command("google-chrome-stable", "--window-size=1300,800", "--window-position=500,200", "--app="+chromeURL)
		cmd.Stdout = nil
		cmd.Stderr = nil
		err = cmd.Start()
		if err != nil {
			logger.Fatal(err)
		}
	}

	os.Exit(0)
}
