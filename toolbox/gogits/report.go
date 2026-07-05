package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strings"
	"time"
)

type templateData struct {
	RepoPath           string
	RepoName           string
	Branch             string
	Version            string
	GeneratedAt        string
	GenDuration        string
	TotalCommits       int
	TotalAdded         int
	TotalDeleted       int
	TotalLoc           int
	TotalFiles         int
	ReportStart        string
	ReportEnd          string
	AgeDays            int
	TotalActiveDays    int
	ActivePct          string
	AvgPerActive       string
	AvgPerDay          string
	AuthorCount        int
	AvgPerAuthor       string
	RecentMonthCommits int
	BusFactorCount     int
	BusFactorPct       string
	RecentMomentumPct  string
	WorkloadDensity    string
	RecentLocPerCommit string
	ChurnRatio         string
	Authors            []AuthorStat
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
	ExtensionStats      []ExtensionStat

	ActivityGrade  string
	ActivityScore  string
	ScaleGrade     string
	ScaleScore     string
    HealthGrade    string
	HealthScore    string
	Commits30d     int
	ActiveDevs30d  int
	ActiveDays30d  int
	DiversityGrade  string
	DiversityScore  string
	GiniCoefficient  string
	CoreTurnover     string
	NewcomerRate     string
	ActiveLayerActive   int
	ActiveLayerSemi     int
	ActiveLayerDormant  int
	ActiveLayerLost     int
	LargeFileCount      int
	LargeFileRatioStr   string
	TodoCount           int
	TodoDensityStr      string
	TestFileCount       int
	TestRatioStr        string
	OldCodeTouchPct     string
	AvgFilesPerCommitStr string
	AbandonedPct    string
	CodeAgeDays     string
	HotspotCount    int
	RecentFileCount int
	Hotspots        []FileHotspot
	TopLinesFiles   []FileHotspot
	DebtGrade       string
	DebtScore       string
	RhythmGrade     string
	RhythmScore     string
	OffHoursPct       string
	OffHoursCommits   int
	OffHoursCommitStr string
	OffHoursLinePct   string
	OffHoursLineStr   string
	ConsistencyPct   string
	ReleaseCount    int
	ActiveWeeks     int
	TotalWeeks      int
	OverallGrade   string
	OverallScore   string
	RadarChartOpt   template.JS
	AuthorReportsJSON template.JS
}

func safeRatio(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func safeTodoDensity(todoCount, totalLoc int) float64 {
	if totalLoc == 0 {
		return 0
	}
	return float64(todoCount) / (float64(totalLoc) / 1000.0)
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

type activeLayers struct {
	Active     int
	SemiActive int
	Dormant    int
	Lost       int
}

func calcGini(authors []AuthorStat) float64 {
	n := len(authors)
	if n <= 1 {
		return 0
	}
	counts := make([]int, n)
	for i, a := range authors {
		counts[i] = a.CommitCount
	}
	sort.Ints(counts)

	total := 0.0
	for _, c := range counts {
		total += float64(c)
	}
	if total == 0 {
		return 0
	}

	sum := 0.0
	for i, c := range counts {
		sum += float64(i+1) * float64(c)
	}
	return (2.0*sum)/(float64(n)*total) - float64(n+1)/float64(n)
}

func calcCoreTurnover(stats []PeriodAuthorStat) float64 {
	if len(stats) <= 1 {
		return 0.25
	}
	totalDist := 0.0
	pairs := 0
	for i := 0; i < len(stats)-1; i++ {
		prev := make(map[string]bool)
		prev[stats[i].TopAuthor] = true
		for _, n := range stats[i].NextTop5 {
			prev[n] = true
		}
		curr := make(map[string]bool)
		curr[stats[i+1].TopAuthor] = true
		for _, n := range stats[i+1].NextTop5 {
			curr[n] = true
		}
		inter := 0
		for name := range prev {
			if curr[name] {
				inter++
			}
		}
		union := len(prev) + len(curr) - inter
		if union > 0 {
			totalDist += 1.0 - float64(inter)/float64(union)
			pairs++
		}
	}
	if pairs == 0 {
		return 0.25
	}
	return totalDist / float64(pairs)
}

func calcNewcomerRate(authors []AuthorStat, now time.Time) float64 {
	threshold := now.AddDate(0, 0, -180)
	newcomers := 0
	for _, a := range authors {
		if a.FirstCommit.After(threshold) {
			newcomers++
		}
	}
	if len(authors) == 0 {
		return 0
	}
	return float64(newcomers) / float64(len(authors))
}

func calcActiveLayers(authors []AuthorStat, now time.Time) activeLayers {
	var layers activeLayers
	for _, a := range authors {
		days := now.Sub(a.LastCommit).Hours() / 24
		switch {
		case days <= 30:
			layers.Active++
		case days <= 90:
			layers.SemiActive++
		case days <= 365:
			layers.Dormant++
		default:
			layers.Lost++
		}
	}
	return layers
}

func scoreToGrade(score float64) string {
	switch {
	case score >= 90:
		return "S"
	case score >= 80:
		return "A"
	case score >= 65:
		return "B"
	case score >= 45:
		return "C"
	case score >= 20:
		return "D"
	default:
		return "E"
	}
}

func calcActivityScore(commits30d, activeDevs30d, activeDays30d int) (string, float64) {
	commitScore := minFloat(float64(commits30d)/50.0, 1.0) * 50
	devScore := minFloat(float64(activeDevs30d)/5.0, 1.0) * 30
	freqScore := 0.0
	if activeDays30d > 0 {
		freqScore = minFloat(float64(commits30d)/float64(activeDays30d)/3.0, 1.0) * 20
	}
	total := commitScore + devScore + freqScore
	return scoreToGrade(total), total
}

func calcScaleScore(loc, totalFiles, contributors int) (string, float64) {
	locScore := minFloat(float64(loc)/100000.0, 1.0) * 40
	fileScore := minFloat(float64(totalFiles)/200.0, 1.0) * 20
	contribScore := minFloat(float64(contributors)/20.0, 1.0) * 40
	total := locScore + fileScore + contribScore
	return scoreToGrade(total), total
}

func calcHealthScore(largeFileCount, totalFiles, todoCount, totalLoc, testFileCount int, oldCodeTouchPct, avgFilesPerCommit float64) (string, float64) {
	// 1. 大文件率 (Large file ratio ≤ 10% → full marks)
	largeRatio := 0.0
	if totalFiles > 0 {
		largeRatio = float64(largeFileCount) / float64(totalFiles)
	}
	largeScore := (1.0 - minFloat(largeRatio/0.1, 1.0)) * 20

	// 2. TODO 密度 (TODO count per 1000 LOC ≤ 5 → full marks)
	todoDensity := 0.0
	if totalLoc > 0 {
		todoDensity = float64(todoCount) / (float64(totalLoc) / 1000.0)
	}
	todoScore := (1.0 - minFloat(todoDensity/5.0, 1.0)) * 20

	// 3. 测试文件比例 (Test file ratio ≥ 30% → full marks)
	testRatio := 0.0
	if totalFiles > 0 {
		testRatio = float64(testFileCount) / float64(totalFiles)
	}
	testScore := minFloat(testRatio/0.3, 1.0) * 20

	// 4. 老代码维护率 (Old code touch rate: optimal [10%, 60%])
	oldCodeScore := 0.0
	if oldCodeTouchPct < 0.1 {
		oldCodeScore = oldCodeTouchPct / 0.1 * 20
	} else if oldCodeTouchPct <= 0.6 {
		oldCodeScore = 20
	} else {
		oldCodeScore = (1.0 - minFloat((oldCodeTouchPct-0.6)/0.4, 1.0)) * 20
	}

	// 5. 提交粒度 (Avg files per commit: ≤ 5 → full marks, 5-15 linear, ≥ 15 → 0)
	granularityScore := 0.0
	if avgFilesPerCommit <= 5 {
		granularityScore = 20
	} else if avgFilesPerCommit <= 15 {
		granularityScore = (1.0 - (avgFilesPerCommit-5.0)/10.0) * 20
	}

	total := largeScore + todoScore + testScore + oldCodeScore + granularityScore
	return scoreToGrade(total), total
}

func calcDiversityScore(authors []AuthorStat, yearStats []PeriodAuthorStat, now time.Time) (string, float64, float64, float64, float64, activeLayers) {
	gini := calcGini(authors)
	turnover := calcCoreTurnover(yearStats)
	newcomerRate := calcNewcomerRate(authors, now)
	layers := calcActiveLayers(authors, now)

	giniScore := maxFloat(25.0*(1.0-minFloat(gini/0.9, 1.0)), 0)

	turnoverScore := 0.0
	if len(yearStats) <= 1 {
		turnoverScore = 15.0
	} else if turnover <= 0.2 {
		turnoverScore = maxFloat(turnover/0.2*12.5, 0)
	} else if turnover <= 0.5 {
		turnoverScore = 12.5 + (turnover-0.2)/0.3*12.5
	} else {
		turnoverScore = maxFloat(25.0*(1.0-(turnover-0.5)/0.5), 0)
	}

	newcomerScore := 0.0
	if newcomerRate <= 0.3 {
		newcomerScore = newcomerRate / 0.3 * 25.0
	} else {
		newcomerScore = maxFloat(25.0*(1.0-(newcomerRate-0.3)/0.3), 0)
	}

	total := layers.Active + layers.SemiActive + layers.Dormant + layers.Lost
	distributionScore := 0.0
	if total > 0 {
		activeRatio := float64(layers.Active+layers.SemiActive) / float64(total)
		distributionScore = activeRatio * 25.0
	}

	totalScore := giniScore + turnoverScore + newcomerScore + distributionScore
	return scoreToGrade(totalScore), totalScore, gini, turnover*100, newcomerRate*100, layers
}

func calcTechDebtScore(hotspots []FileHotspot, abandonedPct float64, codeAgeDays float64) (string, float64) {
	maxCount := 0
	for _, h := range hotspots {
		if h.ModifyCount > maxCount {
			maxCount = h.ModifyCount
		}
	}
	hotspotScore := (1.0 - minFloat(float64(maxCount)/100.0, 1.0)) * 25
	abandonedScore := (1.0 - minFloat(abandonedPct, 1.0)) * 40
	ageScore := 0.0
	switch {
	case codeAgeDays <= 0:
		ageScore = 0
	case codeAgeDays < 30:
		ageScore = codeAgeDays / 30.0 * 20
	case codeAgeDays <= 730:
		ageScore = 35
	default:
		ageScore = 35.0 * (1.0 - minFloat((codeAgeDays-730)/1095.0, 1.0))
	}
	total := hotspotScore + abandonedScore + ageScore
	return scoreToGrade(total), total
}

func calcRhythmScore(consistencyPct, offHoursPct float64, releaseCount int, ageDays int) (string, float64) {
	consistencyScore := minFloat(consistencyPct/100.0, 1.0) * 35
	expected := float64(ageDays) / 90.0
	if expected < 1 {
		expected = 1
	}
	releaseScore := minFloat(float64(releaseCount)/expected, 1.0) * 25
	longevityScore := 0.0
	if ageDays >= 30 {
		longevityScore = minFloat(float64(ageDays)/365.0, 1.0) * 15
	}
	baseScore := consistencyScore + releaseScore + longevityScore
	penaltyRatio := maxFloat((offHoursPct-0.2)/0.5, 0.0)
	penaltyRatio = minFloat(penaltyRatio, 1.0)
	total := baseScore * (1.0 - penaltyRatio*0.5)
	return scoreToGrade(total), total
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

	recentMonthCommits := 0
	for _, s := range result.AuthorSeries {
		for _, v := range s.Data {
			recentMonthCommits += v
		}
	}

	activeDevs30d := 0
	for _, s := range result.AuthorSeries {
		for _, v := range s.Data {
			if v > 0 {
				activeDevs30d++
				break
			}
		}
	}

	activeDays30d := 0
	for i := range result.DateRange {
		for _, s := range result.AuthorSeries {
			if s.Data[i] > 0 {
				activeDays30d++
				break
			}
		}
	}

	busTotal := 0
	busCount := 0
	for _, a := range result.Authors {
		busTotal += a.CommitCount
		busCount++
		if float64(busTotal) >= float64(result.TotalCommits)*0.8 {
			break
		}
	}
	busPct := 0.0
	if result.TotalCommits > 0 {
		busPct = float64(busTotal) / float64(result.TotalCommits) * 100
	}

	recentPct := 0.0
	if result.TotalCommits > 0 {
		recentPct = float64(recentMonthCommits) / float64(result.TotalCommits) * 100
	}

	workloadDensity := 0.0
	if result.TotalCommits > 0 {
		workloadDensity = float64(result.TotalLinesOfCode) / float64(result.TotalCommits)
	}

	churnRatio := 0.0
	if result.TotalDeleted > 0 {
		churnRatio = float64(result.TotalAdded) / float64(result.TotalDeleted)
	}

	recentLocTotal := 0
	for _, v := range result.AddedLineSeries {
		recentLocTotal += v
	}
	for _, v := range result.DeletedLineSeries {
		recentLocTotal += v
	}
	recentLocPerCommit := 0.0
	if recentMonthCommits > 0 {
		recentLocPerCommit = float64(recentLocTotal) / float64(recentMonthCommits)
	}

	actGrade, actScore := calcActivityScore(recentMonthCommits, activeDevs30d, activeDays30d)
	scaleGrade, scaleScore := calcScaleScore(result.TotalLinesOfCode, result.TotalFiles, authorCount)
	healthGrade, healthScore := calcHealthScore(result.LargeFileCount, result.TotalFiles, result.TodoCount, result.TotalLinesOfCode, result.TestFileCount, result.OldCodeTouchPct, result.AvgFilesPerCommit)
	divGrade, divScore, giniVal, turnoverVal, newcomerVal, layers := calcDiversityScore(result.Authors, result.YearAuthorStats, end)
	debtGrade, debtScore := calcTechDebtScore(result.Hotspots, result.AbandonedPct, result.CodeAgeDays)

	offHoursTotal := 0
	allTotal := 0
	for d := 0; d < 7; d++ {
		for h := 0; h < 24; h++ {
			v := result.HourWeekData[d][h]
			allTotal += v
			if d >= 5 || h < 9 || h >= 18 {
				offHoursTotal += v
			}
		}
	}
	offHoursPct := 0.0
	if allTotal > 0 {
		offHoursPct = float64(offHoursTotal) / float64(allTotal)
	}
	offHoursLines := result.OffHoursAdded + result.OffHoursDeleted
	totalLines := result.TotalAdded + result.TotalDeleted
	offHoursLinePct := 0.0
	if totalLines > 0 {
		offHoursLinePct = float64(offHoursLines) / float64(totalLines) * 100
	}
	consistencyPct := 0.0
	if result.TotalWeeks > 0 {
		consistencyPct = float64(result.ActiveWeeks) / float64(result.TotalWeeks) * 100
	}
	rhythmGrade, rhythmScore := calcRhythmScore(consistencyPct, offHoursPct, result.ReleaseCount, ageDays)

	scores6 := []float64{actScore, scaleScore, healthScore, divScore, debtScore, rhythmScore}
	overallSum := 0.0
	for _, s := range scores6 {
		overallSum += s
	}
	overallScore := overallSum / 6.0
	overallGrade := scoreToGrade(overallScore)

	radarOpt, err := buildRadarChartOption(scores6, overallGrade)
	if err != nil {
		return err
	}

	authorReportsJSON := "{}"
	if len(result.AuthorMonthlyReports) > 0 {
		// Build a structure that JavaScript can easily consume
		type authorEntry struct {
			Name    string              `json:"name"`
			Score   float64             `json:"score"`
			Monthly []AuthorMonthlyStat `json:"monthly"`
		}
		var list []authorEntry
		for _, r := range result.AuthorMonthlyReports {
			list = append(list, authorEntry{Name: r.Name, Score: r.Score, Monthly: r.Monthly})
		}
		b, err := json.Marshal(list)
		if err == nil {
			authorReportsJSON = string(b)
		}
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
		RecentMonthCommits: recentMonthCommits,
		BusFactorCount:     busCount,
		BusFactorPct:       fmt.Sprintf("%.1f", busPct),
		RecentMomentumPct:  fmt.Sprintf("%.2f", recentPct),
		WorkloadDensity:    fmt.Sprintf("%.1f", workloadDensity),
		RecentLocPerCommit: fmt.Sprintf("%.1f", recentLocPerCommit),
		ChurnRatio:         fmt.Sprintf("%.1f", churnRatio),
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
		ExtensionStats:       result.ExtensionStats,
		ActivityGrade:     actGrade,
		ActivityScore:     fmt.Sprintf("%.0f", actScore),
		ScaleGrade:        scaleGrade,
		ScaleScore:        fmt.Sprintf("%.0f", scaleScore),
		HealthGrade:       healthGrade,
		HealthScore:       fmt.Sprintf("%.0f", healthScore),
		LargeFileCount:      result.LargeFileCount,
		LargeFileRatioStr:   fmt.Sprintf("%.1f", safeRatio(float64(result.LargeFileCount), float64(result.TotalFiles))*100),
		TodoCount:           result.TodoCount,
		TodoDensityStr:      fmt.Sprintf("%.1f", safeTodoDensity(result.TodoCount, result.TotalLinesOfCode)),
		TestFileCount:       result.TestFileCount,
		TestRatioStr:        fmt.Sprintf("%.1f", safeRatio(float64(result.TestFileCount), float64(result.TotalFiles))*100),
		OldCodeTouchPct:     fmt.Sprintf("%.1f", result.OldCodeTouchPct*100),
		AvgFilesPerCommitStr: fmt.Sprintf("%.1f", result.AvgFilesPerCommit),
		Commits30d:        recentMonthCommits,
		ActiveDevs30d:     activeDevs30d,
		ActiveDays30d:     activeDays30d,
		DiversityGrade:    divGrade,
		DiversityScore:    fmt.Sprintf("%.0f", divScore),
		GiniCoefficient:   fmt.Sprintf("%.2f", giniVal),
		CoreTurnover:      fmt.Sprintf("%.0f", turnoverVal),
		NewcomerRate:      fmt.Sprintf("%.0f", newcomerVal),
		ActiveLayerActive: layers.Active,
		ActiveLayerSemi:   layers.SemiActive,
		ActiveLayerDormant: layers.Dormant,
		ActiveLayerLost:   layers.Lost,
		DebtGrade:         debtGrade,
		DebtScore:         fmt.Sprintf("%.0f", debtScore),
		RhythmGrade:       rhythmGrade,
		RhythmScore:       fmt.Sprintf("%.0f", rhythmScore),
		OffHoursPct:       fmt.Sprintf("%.1f", offHoursPct*100),
		OffHoursCommits:   result.OffHoursCommits,
		OffHoursCommitStr: fmt.Sprintf("%d / %d", result.OffHoursCommits, result.TotalCommits),
		OffHoursLinePct:   fmt.Sprintf("%.1f", offHoursLinePct),
		OffHoursLineStr:   fmt.Sprintf("%d / %d", offHoursLines, totalLines),
		ConsistencyPct:    fmt.Sprintf("%.1f", consistencyPct),
		ReleaseCount:      result.ReleaseCount,
		ActiveWeeks:       result.ActiveWeeks,
		TotalWeeks:        result.TotalWeeks,
		AbandonedPct:      fmt.Sprintf("%.1f", result.AbandonedPct*100),
		CodeAgeDays:       fmt.Sprintf("%.0f", result.CodeAgeDays),
		HotspotCount:      len(result.Hotspots),
		RecentFileCount:   result.RecentFileCount,
		Hotspots:          result.Hotspots,
		TopLinesFiles:     result.TopLinesFiles,
		OverallGrade:      overallGrade,
		OverallScore:      fmt.Sprintf("%.0f", overallScore),
		RadarChartOpt:     template.JS(radarOpt),
		AuthorReportsJSON: template.JS(authorReportsJSON),
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
		"sub": func(a, b int) int {
			return a - b
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