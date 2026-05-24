package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type CommitInfo struct {
	Hash    string
	Author  string
	Email   string
	Date    time.Time
	Added   int
	Deleted int
}

type AuthorStat struct {
	Name         string
	Email        string
	CommitCount  int
	FirstCommit  time.Time
	LastCommit   time.Time
	ActiveDays   int
	AddedLines   int
	DeletedLines int
}

type AuthorDayData struct {
	Name string
	Data []int
}

type AnalysisResult struct {
	RepoPath          string
	RepoName          string
	Branch            string
	TotalCommits      int
	TotalAdded        int
	TotalDeleted      int
	ReportStart       time.Time
	ReportEnd         time.Time
	TotalActiveDays   int
	TotalFiles        int
	TotalLinesOfCode  int
	Authors           []AuthorStat
	DateRange         []string
	AuthorSeries      []AuthorDayData
	AddedLineSeries   []int
	DeletedLineSeries []int
	AuthorAddedSeries   []AuthorDayData
	AuthorDeletedSeries []AuthorDayData
	HourWeekData        [7][24]int
	MonthOfYearData     [12]int
	YearMonthLabels     []string
	YearMonthData       []int
	YearMonthAddedData  []int
	YearMonthDeletedData []int
	GenerationDuration  time.Duration
}

const gitTimeFormat = "2006-01-02 15:04:05 -0700"

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

	return result, nil
}

func getBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository or git not available: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func getTotalFiles(repoPath string) int {
	cmd := exec.Command("git", "-C", repoPath, "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}

func getTotalLines(repoPath string) int {
	cmd := exec.Command("git", "-C", repoPath, "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	files := strings.Fields(string(out))
	var total int
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(repoPath, f))
		if err != nil {
			continue
		}
		total += bytes.Count(data, []byte{'\n'})
	}
	return total
}

func runGitLog(repoPath string) ([]CommitInfo, error) {
	cmd := exec.Command("git", "-C", repoPath, "log",
		"--format=COMMIT%n%H|%an|%ae|%ai",
		"--numstat", "HEAD")

	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	defer cmd.Wait()

	scanner := bufio.NewScanner(out)
	var commits []CommitInfo
	state := 0
	var current *CommitInfo

	for scanner.Scan() {
		line := scanner.Text()
		switch state {
		case 0:
			if line == "COMMIT" {
				state = 1
			}
		case 1:
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "|", 4)
			if len(parts) >= 4 {
				t, err := time.Parse(gitTimeFormat, parts[3])
				if err != nil {
					t = time.Now()
				}
				current = &CommitInfo{
					Hash:   parts[0],
					Author: parts[1],
					Email:  parts[2],
					Date:   t,
				}
				state = 2
			}
		case 2:
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "COMMIT") {
				if current != nil {
					commits = append(commits, *current)
					current = nil
				}
				state = 1
			} else {
				parts := strings.SplitN(line, "\t", 3)
				if len(parts) >= 3 {
					added, err1 := strconv.Atoi(parts[0])
					deleted, err2 := strconv.Atoi(parts[1])
					if err1 == nil && err2 == nil {
						current.Added += added
						current.Deleted += deleted
					}
				}
			}
		}
	}

	if current != nil {
		commits = append(commits, *current)
	}

	return commits, scanner.Err()
}

type authorDayStats struct {
	commits int
	added   int
	deleted int
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

	totalActiveDates := make(map[string]bool)
	totalAdded, totalDeleted := 0, 0
	var firstCommit, lastCommit time.Time
	var hourWeekData [7][24]int
	var monthOfYearData [12]int
	type ymVal struct {
		commits, added, deleted int
	}
	yearMonthMap := make(map[string]*ymVal)

	for _, c := range commits {
		key := authorKey(c.Author, c.Email)

		stat, ok := authorMap[key]
		if !ok {
			stat = &AuthorStat{
				Name:   c.Author,
				Email:  c.Email,
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

		weekday := (c.Date.Weekday() + 6) % 7
		hour := c.Date.Hour()
		hourWeekData[weekday][hour]++

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
	sort.Strings(ymKeys)
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

	return &AnalysisResult{
		RepoPath:           absPath,
		RepoName:           repoName,
		Branch:             branch,
		TotalCommits:       len(commits),
		TotalAdded:         totalAdded,
		TotalDeleted:       totalDeleted,
		ReportStart:        firstCommit,
		ReportEnd:          lastCommit,
		TotalActiveDays:    totalActiveDays,
		Authors:            authors,
		DateRange:          dateRange,
		AuthorSeries:       authorSeriesList,
		AddedLineSeries:    addedSeries,
		DeletedLineSeries:  deletedSeries,
		AuthorAddedSeries:  authorAddedList,
		AuthorDeletedSeries: authorDeletedList,
		HourWeekData:        hourWeekData,
		MonthOfYearData:     monthOfYearData,
		YearMonthLabels:      ymLabels,
		YearMonthData:        ymData,
		YearMonthAddedData:   ymAdded,
		YearMonthDeletedData: ymDeleted,
	}
}
