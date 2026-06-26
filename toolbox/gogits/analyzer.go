package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func Analyze(repoPath string) (*AnalysisResult, error) {
	branch, err := getBranch(repoPath)
	if err != nil {
		return nil, fmt.Errorf("get branch: %w", err)
	}

	commits, err := runGitLog(repoPath)
	if err != nil {
		return nil, fmt.Errorf("git log: %w", err)
	}

	absPath, _ := filepath.Abs(repoPath)
	repoName := filepath.Base(absPath)

	result := buildResult(absPath, repoName, branch, commits)

	result.TotalFiles = getTotalFiles(repoPath)
	result.TotalLinesOfCode = getTotalLines(repoPath)
	result.LargeFileCount = getLargeFileCount(repoPath, 1000)
	result.TodoCount = getTodoCount(repoPath)
	result.TestFileCount = getTestFileCount(repoPath)
	result.AvgFilesPerCommit = getAvgFilesPerCommit(repoPath)
	result.Hotspots = getHotspots(repoPath, 8)
	result.TopLinesFiles = getTopLinesFiles(repoPath, 8)
	result.AbandonedPct, result.AbandonedLOC = getAbandonedData(repoPath, result.TotalLinesOfCode)
	result.CodeAgeDays = getCodeAgeDays(repoPath)
	result.ReleaseCount = getReleaseCount(repoPath)
	result.OldCodeTouchPct = getOldCodeTouchRate(repoPath)
	result.RecentFileCount = getRecentFileCount(repoPath)
	result.FileChartLabels, result.FileChartData = getDailyFileCounts(repoPath)
	result.LocChartLabels, result.LocChartData = getDailyLocCounts(repoPath, result.TotalLinesOfCode)
	result.ExtensionStats = getExtensionStats(repoPath)

	authorMonthFiles, authorMonthCommits, err := extractAuthorMonthlyStats(repoPath)
	if err == nil {
		result.AuthorMonthlyReports = buildAuthorMonthlyReports(result.Authors, authorMonthFiles, authorMonthCommits, repoPath)
	}

	return result, nil
}

func buildResult(absPath, repoName, branch string, commits []CommitInfo) *AnalysisResult {
	now := time.Now()
	start := now.AddDate(0, 0, -29)

	var dateRange []string
	for d := start; !d.After(now); d = d.AddDate(0, 0, 1) {
		dateRange = append(dateRange, d.Format("2006-01-02"))
	}

	dateInRange := make(map[string]bool, len(dateRange))
	for _, d := range dateRange {
		dateInRange[d] = true
	}

	authorKey := func(name, email string) string {
		return name + "|" + email
	}

	authorMap := make(map[string]*AuthorStat)
	authorActiveDates := make(map[string]map[string]bool)
	authorDaily := make(map[string]map[string]*authorDayStats)
	dailyAdded := make(map[string]int)
	dailyDeleted := make(map[string]int)
	authorAllDayCommit := make(map[string]map[string]int)
	authorAllDayAdded := make(map[string]map[string]int)

	totalActiveDates := make(map[string]bool)
	totalAdded, totalDeleted := 0, 0
	offHoursCommits, offHoursAdded, offHoursDeleted := 0, 0, 0
	var firstCommit, lastCommit time.Time
	activeWeeks := make(map[string]bool)
	var hourWeekData [7][24]int
	var monthOfYearData [12]int
	type ymVal struct {
		commits, added, deleted int
	}
	yearMonthMap := make(map[string]*ymVal)
	monthAuthorCommits := make(map[string]map[string]int)
	yearAuthorCommits := make(map[string]map[string]int)

	for _, c := range commits {
		key := authorKey(c.Author, c.Email)

		stat, ok := authorMap[key]
		if !ok {
			stat = &AuthorStat{
				Name:        c.Author,
				Email:       c.Email,
				FirstCommit: c.Date,
				LastCommit:  c.Date,
			}
			authorMap[key] = stat
		}

		stat.CommitCount++
		stat.AddedLines += c.Added
		stat.DeletedLines += c.Deleted

		if c.Date.Before(stat.FirstCommit) {
			stat.FirstCommit = c.Date
		}
		if c.Date.After(stat.LastCommit) {
			stat.LastCommit = c.Date
		}

		totalAdded += c.Added
		totalDeleted += c.Deleted

		if firstCommit.IsZero() || c.Date.Before(firstCommit) {
			firstCommit = c.Date
		}
		if lastCommit.IsZero() || c.Date.After(lastCommit) {
			lastCommit = c.Date
		}

		dateStr := c.Date.Format("2006-01-02")
		year, week := c.Date.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)
		activeWeeks[weekKey] = true

		weekday := (c.Date.Weekday() + 6) % 7
		hour := c.Date.Hour()
		hourWeekData[weekday][hour]++
		if weekday >= 5 || hour < 9 || hour >= 18 {
			offHoursCommits++
			offHoursAdded += c.Added
			offHoursDeleted += c.Deleted
		}

		month := c.Date.Month() - 1
		monthOfYearData[month]++
		ymKey := c.Date.Format("2006-01")
		ym, ok := yearMonthMap[ymKey]
		if !ok {
			ym = &ymVal{}
			yearMonthMap[ymKey] = ym
		}
		ym.commits++
		ym.added += c.Added
		ym.deleted += c.Deleted
		if monthAuthorCommits[ymKey] == nil {
			monthAuthorCommits[ymKey] = make(map[string]int)
		}
		monthAuthorCommits[ymKey][key]++

		yearKey := c.Date.Format("2006")
		if yearAuthorCommits[yearKey] == nil {
			yearAuthorCommits[yearKey] = make(map[string]int)
		}
		yearAuthorCommits[yearKey][key]++

		if authorAllDayCommit[key] == nil {
			authorAllDayCommit[key] = make(map[string]int)
			authorAllDayAdded[key] = make(map[string]int)
		}
		authorAllDayCommit[key][dateStr]++
		authorAllDayAdded[key][dateStr] += c.Added

		if authorActiveDates[key] == nil {
			authorActiveDates[key] = make(map[string]bool)
		}
		authorActiveDates[key][dateStr] = true
		totalActiveDates[dateStr] = true

		if dateInRange[dateStr] {
			if authorDaily[key] == nil {
				authorDaily[key] = make(map[string]*authorDayStats)
			}
			dc, ok := authorDaily[key][dateStr]
			if !ok {
				dc = &authorDayStats{}
				authorDaily[key][dateStr] = dc
			}
			dc.commits++
			dc.added += c.Added
			dc.deleted += c.Deleted
			dailyAdded[dateStr] += c.Added
			dailyDeleted[dateStr] += c.Deleted
		}
	}

	for key, stat := range authorMap {
		stat.ActiveDays = len(authorActiveDates[key])
	}

	totalActiveDays := len(totalActiveDates)
	totalWeeks := 0
	if !firstCommit.IsZero() && !lastCommit.IsZero() {
		fy, fw := firstCommit.ISOWeek()
		ly, lw := lastCommit.ISOWeek()
		totalWeeks = (ly-fy)*52 + (lw - fw) + 1
		if totalWeeks < 1 {
			totalWeeks = 1
		}
	}

	var authors []AuthorStat
	for _, stat := range authorMap {
		authors = append(authors, *stat)
	}
	sort.Slice(authors, func(i, j int) bool {
		if authors[i].CommitCount != authors[j].CommitCount {
			return authors[i].CommitCount > authors[j].CommitCount
		}
		return authors[i].Name < authors[j].Name
	})

	addedSeries := make([]int, len(dateRange))
	deletedSeries := make([]int, len(dateRange))

	var authorSeriesList []AuthorDayData
	var authorAddedList []AuthorDayData
	var authorDeletedList []AuthorDayData

	for _, a := range authors {
		key := authorKey(a.Name, a.Email)
		commitData := make([]int, len(dateRange))
		addData := make([]int, len(dateRange))
		delData := make([]int, len(dateRange))
		for i, d := range dateRange {
			if ds, ok := authorDaily[key][d]; ok {
				commitData[i] = ds.commits
				addData[i] = ds.added
				delData[i] = ds.deleted
			}
		}
		authorSeriesList = append(authorSeriesList, AuthorDayData{Name: a.Name, Data: commitData})
		authorAddedList = append(authorAddedList, AuthorDayData{Name: a.Name, Data: addData})
		authorDeletedList = append(authorDeletedList, AuthorDayData{Name: a.Name, Data: delData})
	}

	for i, d := range dateRange {
		addedSeries[i] = dailyAdded[d]
		deletedSeries[i] = dailyDeleted[d]
	}

	ymKeys := make([]string, 0, len(yearMonthMap))
	for k := range yearMonthMap {
		ymKeys = append(ymKeys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(ymKeys)))
	ymLabels := make([]string, len(ymKeys))
	ymData := make([]int, len(ymKeys))
	ymAdded := make([]int, len(ymKeys))
	ymDeleted := make([]int, len(ymKeys))
	for i, k := range ymKeys {
		ymLabels[i] = k
		ymData[i] = yearMonthMap[k].commits
		ymAdded[i] = yearMonthMap[k].added
		ymDeleted[i] = yearMonthMap[k].deleted
	}

	allDays, allDayLabels := dayRange(firstCommit, lastCommit)

	var cumCommitSeries []AuthorDayData
	var cumAddedSeries []AuthorDayData
	for _, a := range authors {
		key := authorKey(a.Name, a.Email)
		cc := 0
		ca := 0
		commitData := make([]int, len(allDays))
		addedData := make([]int, len(allDays))
		for i, t := range allDays {
			ds := t.Format("2006-01-02")
			cc += authorAllDayCommit[key][ds]
			ca += authorAllDayAdded[key][ds]
			commitData[i] = cc
			addedData[i] = ca
		}
		cumCommitSeries = append(cumCommitSeries, AuthorDayData{Name: a.Name, Data: commitData})
		cumAddedSeries = append(cumAddedSeries, AuthorDayData{Name: a.Name, Data: addedData})
	}

	monthKeys := make([]string, 0, len(monthAuthorCommits))
	for k := range monthAuthorCommits {
		monthKeys = append(monthKeys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(monthKeys)))
	var monthStats []PeriodAuthorStat
	for _, mk := range monthKeys {
		ac := monthAuthorCommits[mk]
		total := 0
		for _, c := range ac {
			total += c
		}
		topAuthor, topCommits := "", 0
		var others []struct{ name string; cnt int }
		for k, cnt := range ac {
			if cnt > topCommits {
				topCommits = cnt
				topAuthor = k
			}
			others = append(others, struct{ name string; cnt int }{k, cnt})
		}
		sort.Slice(others, func(i, j int) bool {
			if others[i].cnt != others[j].cnt {
				return others[i].cnt > others[j].cnt
			}
			return others[i].name < others[j].name
		})
		var nextTop5 []string
		for _, o := range others {
			if o.name == topAuthor {
				continue
			}
			parts := strings.SplitN(o.name, "|", 2)
			nextTop5 = append(nextTop5, parts[0])
			if len(nextTop5) >= 5 {
				break
			}
		}
		parts := strings.SplitN(topAuthor, "|", 2)
		monthStats = append(monthStats, PeriodAuthorStat{
			Period:       mk,
			TopAuthor:    parts[0],
			TopCommits:   topCommits,
			TotalCommits: total,
			NextTop5:     nextTop5,
			AuthorCount:  len(ac),
		})
	}

	yearKeys := make([]string, 0, len(yearAuthorCommits))
	for k := range yearAuthorCommits {
		yearKeys = append(yearKeys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(yearKeys)))
	var yearStats []PeriodAuthorStat
	for _, yk := range yearKeys {
		ac := yearAuthorCommits[yk]
		total := 0
		for _, c := range ac {
			total += c
		}
		topAuthor, topCommits := "", 0
		var others []struct{ name string; cnt int }
		for k, cnt := range ac {
			if cnt > topCommits {
				topCommits = cnt
				topAuthor = k
			}
			others = append(others, struct{ name string; cnt int }{k, cnt})
		}
		sort.Slice(others, func(i, j int) bool {
			if others[i].cnt != others[j].cnt {
				return others[i].cnt > others[j].cnt
			}
			return others[i].name < others[j].name
		})
		var nextTop5 []string
		for _, o := range others {
			if o.name == topAuthor {
				continue
			}
			parts := strings.SplitN(o.name, "|", 2)
			nextTop5 = append(nextTop5, parts[0])
			if len(nextTop5) >= 5 {
				break
			}
		}
		parts := strings.SplitN(topAuthor, "|", 2)
		yearStats = append(yearStats, PeriodAuthorStat{
			Period:       yk,
			TopAuthor:    parts[0],
			TopCommits:   topCommits,
			TotalCommits: total,
			NextTop5:     nextTop5,
			AuthorCount:  len(ac),
		})
	}

	return &AnalysisResult{
		RepoPath:            absPath,
		RepoName:            repoName,
		Branch:              branch,
		TotalCommits:        len(commits),
		TotalAdded:          totalAdded,
		TotalDeleted:        totalDeleted,
		OffHoursCommits:     offHoursCommits,
		OffHoursAdded:       offHoursAdded,
		OffHoursDeleted:     offHoursDeleted,
		ReportStart:         firstCommit,
		ReportEnd:           lastCommit,
		TotalActiveDays:     totalActiveDays,
		Authors:             authors,
		DateRange:           dateRange,
		AuthorSeries:        authorSeriesList,
		AddedLineSeries:     addedSeries,
		DeletedLineSeries:   deletedSeries,
		AuthorAddedSeries:   authorAddedList,
		AuthorDeletedSeries: authorDeletedList,
		HourWeekData:        hourWeekData,
		MonthOfYearData:     monthOfYearData,
		YearMonthLabels:     ymLabels,
		YearMonthData:       ymData,
		YearMonthAddedData:  ymAdded,
		YearMonthDeletedData: ymDeleted,
		AuthorCumCommitSeries: cumCommitSeries,
		AuthorCumAddedSeries:  cumAddedSeries,
		AllDayLabels:          allDayLabels,
		ActiveWeeks:           len(activeWeeks),
		TotalWeeks:            totalWeeks,
		MonthAuthorStats:      monthStats,
		YearAuthorStats:       yearStats,
	}
}

func buildAuthorMonthlyReports(authors []AuthorStat, authorMonthFiles map[string]map[string][]string, authorMonthCommits map[string]map[string]int, repoPath string) []AuthorMonthlyReport {
	allMonths := make(map[string]bool)
	for _, mf := range authorMonthFiles {
		for m := range mf {
			allMonths[m] = true
		}
	}
	for _, mc := range authorMonthCommits {
		for m := range mc {
			allMonths[m] = true
		}
	}
	var monthList []string
	for m := range allMonths {
		monthList = append(monthList, m)
	}
	sort.Strings(monthList)

	totalFilesByMonth := getTotalFilesByMonth(repoPath, monthList)

	var reports []AuthorMonthlyReport
	for _, a := range authors {
		author := a.Name
		monthFiles := authorMonthFiles[author]
		monthCommits := authorMonthCommits[author]

		cumFiles := make(map[string]bool)
		var monthly []AuthorMonthlyStat
		for _, m := range monthList {
			if files, ok := monthFiles[m]; ok {
				for _, f := range files {
					cumFiles[f] = true
				}
			}
			fileSet := totalFilesByMonth[m]
			intersection := 0
			totalCount := 0
			if fileSet != nil {
				totalCount = len(fileSet)
				for f := range cumFiles {
					if fileSet[f] {
						intersection++
					}
				}
			}
			breadth := 0.0
			if totalCount > 0 {
				breadth = float64(intersection) / float64(totalCount)
			}
			commits := 0
			if monthCommits != nil {
				commits = monthCommits[m]
			}
			monthly = append(monthly, AuthorMonthlyStat{
				Month:      m,
				Commits:    commits,
				CumFiles:   intersection,
				TotalFiles: totalCount,
				Breadth:    breadth,
			})
		}

		score := 0.0
		if len(monthly) > 0 {
			last := monthly[len(monthly)-1]
			normActivity := minFloat(float64(last.Commits)/50.0, 1.0)
			score = last.Breadth*100*0.5 + normActivity*50.0
		}
		reports = append(reports, AuthorMonthlyReport{
			Name:    author,
			Score:   score,
			Monthly: monthly,
		})
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Score > reports[j].Score
	})
	return reports
}