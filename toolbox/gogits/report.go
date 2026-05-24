package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"
)

type templateData struct {
	RepoPath         string
	RepoName         string
	Branch           string
	Version          string
	GeneratedAt      string
	GenDuration      string
	TotalCommits     int
	TotalAdded       int
	TotalDeleted     int
	TotalLoc         int
	TotalFiles       int
	ReportStart      string
	ReportEnd        string
	AgeDays          int
	TotalActiveDays  int
	ActivePct        string
	AvgPerActive     string
	AvgPerDay        string
	AuthorCount      int
	AvgPerAuthor     string
	Authors          []AuthorStat
	CommitChartOpt      template.JS
	LineChartOpt        template.JS
	AuthorLineChartOpt  template.JS
	HourWeekOpt         template.JS
	MonthOfYearOpt      template.JS
	YearMonthOpt        template.JS
	YearMonthLabels     []string
	YearMonthData       []int
	YearMonthAddedData  []int
	YearMonthDeletedData []int
	CumCommitOpt        template.JS
	CumAddedOpt         template.JS
	MonthAuthorStats    []PeriodAuthorStat
	YearAuthorStats     []PeriodAuthorStat
	FilesChartOpt       template.JS
	LocChartOpt         template.JS
}

func GenerateReport(result *AnalysisResult, outputPath string) error {
	commitOpt, err := buildCommitChartOption(result)
	if err != nil {
		return err
	}
	lineOpt, err := buildLineChartOption(result)
	if err != nil {
		return err
	}
	authorLineOpt, err := buildAuthorLineChartOption(result)
	if err != nil {
		return err
	}
	hourWeekOpt, err := buildHourWeekChartOption(result)
	if err != nil {
		return err
	}
	monthOfYearOpt, err := buildMonthOfYearChartOption(result)
	if err != nil {
		return err
	}
	yearMonthOpt, err := buildYearMonthChartOption(result)
	if err != nil {
		return err
	}
	cumCommitOpt, err := buildCumChartOption(result, true)
	if err != nil {
		return err
	}
	cumAddedOpt, err := buildCumChartOption(result, false)
	if err != nil {
		return err
	}
	fileOpt, err := buildFileChartOption(result)
	if err != nil {
		return err
	}
	locOpt, err := buildLocChartOption(result)
	if err != nil {
		return err
	}

	start := result.ReportStart
	end := result.ReportEnd
	ageDays := 0
	if !start.IsZero() && !end.IsZero() {
		ageDays = int(end.Sub(start).Hours() / 24)
		if ageDays < 1 {
			ageDays = 1
		}
	}

	activePct := 0.0
	if ageDays > 0 {
		activePct = float64(result.TotalActiveDays) / float64(ageDays) * 100
	}

	avgPerActive := 0.0
	if result.TotalActiveDays > 0 {
		avgPerActive = float64(result.TotalCommits) / float64(result.TotalActiveDays)
	}

	avgPerDay := 0.0
	if ageDays > 0 {
		avgPerDay = float64(result.TotalCommits) / float64(ageDays)
	}

	avgPerAuthor := 0.0
	authorCount := len(result.Authors)
	if authorCount > 0 {
		avgPerAuthor = float64(result.TotalCommits) / float64(authorCount)
	}

	dur := result.GenerationDuration
	durStr := fmt.Sprintf("%.1fs", dur.Seconds())
	if dur.Seconds() < 1 {
		durStr = fmt.Sprintf("%dms", dur.Milliseconds())
	}

	data := templateData{
		RepoPath:           result.RepoPath,
		RepoName:           result.RepoName,
		Branch:             result.Branch,
		Version:            "1.0.0",
		GeneratedAt:        time.Now().Format("2006-01-02 15:04:05"),
		GenDuration:        durStr,
		TotalCommits:       result.TotalCommits,
		TotalAdded:         result.TotalAdded,
		TotalDeleted:       result.TotalDeleted,
		TotalLoc:           result.TotalLinesOfCode,
		TotalFiles:         result.TotalFiles,
		ReportStart:        result.ReportStart.Format("2006-01-02 15:04:05"),
		ReportEnd:          result.ReportEnd.Format("2006-01-02 15:04:05"),
		AgeDays:            ageDays,
		TotalActiveDays:    result.TotalActiveDays,
		ActivePct:          fmt.Sprintf("%.2f", activePct),
		AvgPerActive:       fmt.Sprintf("%.1f", avgPerActive),
		AvgPerDay:          fmt.Sprintf("%.1f", avgPerDay),
		AuthorCount:        authorCount,
		AvgPerAuthor:       fmt.Sprintf("%.1f", avgPerAuthor),
		Authors:            result.Authors,
		CommitChartOpt:     template.JS(commitOpt),
		LineChartOpt:       template.JS(lineOpt),
		AuthorLineChartOpt:  template.JS(authorLineOpt),
		HourWeekOpt:         template.JS(hourWeekOpt),
		MonthOfYearOpt:      template.JS(monthOfYearOpt),
		YearMonthOpt:        template.JS(yearMonthOpt),
		YearMonthLabels:      result.YearMonthLabels,
		YearMonthData:        result.YearMonthData,
		YearMonthAddedData:   result.YearMonthAddedData,
		YearMonthDeletedData: result.YearMonthDeletedData,
		CumCommitOpt:         template.JS(cumCommitOpt),
		CumAddedOpt:          template.JS(cumAddedOpt),
		MonthAuthorStats:     result.MonthAuthorStats,
		YearAuthorStats:      result.YearAuthorStats,
		FilesChartOpt:        template.JS(fileOpt),
		LocChartOpt:          template.JS(locOpt),
	}

	funcMap := template.FuncMap{
		"percent": func(a, b int) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b) * 100
		},
		"join": func(a []string, sep string) string {
			return strings.Join(a, sep)
		},
		"avgFileSize": func(totalLines, totalFiles int) float64 {
			if totalFiles == 0 {
				return 0
			}
			return float64(totalLines) / float64(totalFiles)
		},
		"age": func(first, last time.Time) string {
			if first.IsZero() || last.IsZero() || last.Before(first) {
				return ""
			}
			d := last.Sub(first)
			days := int(d.Hours() / 24)
			hours := int(d.Hours()) % 24
			mins := int(d.Minutes()) % 60
			secs := int(d.Seconds()) % 60
			return fmt.Sprintf("%d days, %02d:%02d:%02d", days, hours, mins, secs)
		},
	}
	tmpl, err := template.New("report").Funcs(funcMap).Parse(reportTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}