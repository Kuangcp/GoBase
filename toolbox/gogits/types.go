package main

import "time"

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

type FileHotspot struct {
	Path        string
	ModifyCount int
}

type ExtensionStat struct {
	Extension string
	FileCount int
	LineCount int
}

type AuthorDayData struct {
	Name string
	Data []int
}

type PeriodAuthorStat struct {
	Period       string
	TopAuthor    string
	TopCommits   int
	TotalCommits int
	NextTop5     []string
	AuthorCount  int
}

type AnalysisResult struct {
	RepoPath             string
	RepoName             string
	Branch               string
	TotalCommits         int
	TotalAdded           int
	TotalDeleted         int
	ReportStart          time.Time
	ReportEnd            time.Time
	TotalActiveDays      int
	TotalFiles           int
	TotalLinesOfCode     int
	Authors              []AuthorStat
	DateRange            []string
	AuthorSeries         []AuthorDayData
	AddedLineSeries      []int
	DeletedLineSeries    []int
	AuthorAddedSeries    []AuthorDayData
	AuthorDeletedSeries  []AuthorDayData
	AuthorCumCommitSeries []AuthorDayData
	AuthorCumAddedSeries  []AuthorDayData
	AllDayLabels         []string
	HourWeekData         [7][24]int
	MonthOfYearData      [12]int
	YearMonthLabels      []string
	YearMonthData        []int
	YearMonthAddedData   []int
	YearMonthDeletedData []int
	MonthAuthorStats     []PeriodAuthorStat
	YearAuthorStats      []PeriodAuthorStat
	FileChartLabels      []string
	FileChartData        []int
	LocChartLabels       []string
	LocChartData         []int
	LargeFileCount       int
	TodoCount            int
	OldCodeTouchPct      float64
	TestFileCount        int
	AvgFilesPerCommit    float64
	OffHoursCommits      int
	OffHoursAdded        int
	OffHoursDeleted      int
	RecentFileCount      int
	Hotspots             []FileHotspot
	TopLinesFiles         []FileHotspot
	AbandonedPct         float64
	AbandonedLOC         int
	CodeAgeDays          float64
	ActiveWeeks          int
	TotalWeeks           int
	ReleaseCount         int
	GenerationDuration   time.Duration
	ExtensionStats       []ExtensionStat
}

type authorDayStats struct {
	commits int
	added   int
	deleted int
}

const gitTimeFormat = "2006-01-02 15:04:05 -0700"