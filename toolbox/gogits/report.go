package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"time"
)

const reportTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Git Report - {{.RepoName}}</title>
<script src="https://cdn.jsdelivr.net/npm/echarts@5/dist/echarts.min.js"></script>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  background: #f5f5f5; color: #333; }
.header { background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  color: #fff; padding: 28px 0; }
.header .wrap { max-width: 1200px; margin: 0 auto; padding: 0 20px; }
.header h1 { font-size: 22px; font-weight: 600; }
.header p { font-size: 13px; opacity: .75; margin-top: 4px; }
.wrap { max-width: 1200px; margin: 0 auto; padding: 0 20px; }

.tabs { display: flex; gap: 0; border-bottom: 2px solid #e9ecef; margin: 20px 0 24px; }
.tab { padding: 10px 28px; cursor: pointer; border: none; background: none;
  font-size: 14px; font-weight: 500; color: #888;
  border-bottom: 2px solid transparent; margin-bottom: -2px; transition: all .2s; }
.tab:hover { color: #1a1a2e; }
.tab.active { color: #1a1a2e; border-bottom-color: #1a1a2e; background: rgba(26,26,46,.04); }

.tab-content { display: none; }
.tab-content.active { display: block; }

.section { background: #fff; border-radius: 8px; padding: 24px; margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0,0,0,.1); }
.section h2 { font-size: 17px; margin-bottom: 16px; color: #1a1a2e; }

table { width: 100%; border-collapse: collapse; font-size: 13px; }
th { background: #f8f9fa; text-align: left; padding: 9px 12px; font-weight: 600;
  border-bottom: 2px solid #e9ecef; white-space: nowrap; }
td { padding: 9px 12px; border-bottom: 1px solid #e9ecef; }
tr:hover td { background: #f8f9fa; }
.num { text-align: right; font-variant-numeric: tabular-nums; }
.added { color: #4caf50; font-weight: 600; }
.deleted { color: #f44336; font-weight: 600; }
.chart-box { width: 100%; height: 420px; }

.stats-grid { display: grid; grid-template-columns: auto 1fr; gap: 2px 32px;
  font-size: 14px; line-height: 2; }
.stats-label { color: #888; white-space: nowrap; text-align: right; }
.stats-value { color: #333; font-weight: 500; }

@media (max-width: 640px) {
  .chart-box { height: 280px; }
  table { font-size: 11px; }
  th, td { padding: 5px 6px; }
  .stats-grid { font-size: 12px; gap: 2px 16px; }
  .tab { padding: 8px 14px; font-size: 13px; }
}
</style>
</head>
<body>

<div class="header">
<div class="wrap">
<h1>{{.RepoName}}</h1>
<p>{{.RepoPath}} &nbsp;·&nbsp; {{.Branch}} branch &nbsp;·&nbsp; Generated {{.GeneratedAt}}</p>
</div>
</div>

<div class="wrap">
<div class="tabs">
<button class="tab active" data-tab="general">General</button>
<button class="tab" data-tab="activity">Activity</button>
<button class="tab" data-tab="authors">Authors</button>
</div>

<div id="general" class="tab-content active">
<div class="section">
<h2>Overview</h2>
<div class="stats-grid">
<span class="stats-label">Project name</span><span class="stats-value">{{.RepoName}}</span>
<span class="stats-label">Generated</span><span class="stats-value">{{.GeneratedAt}} (in {{.GenDuration}})</span>
<span class="stats-label">Generator</span><span class="stats-value">gogits {{.Version}}</span>
<span class="stats-label">Report Period</span><span class="stats-value">{{.ReportStart}} to {{.ReportEnd}}</span>
<span class="stats-label">Age</span><span class="stats-value">{{.AgeDays}} days, {{.TotalActiveDays}} active days ({{.ActivePct}}%)</span>
<span class="stats-label">Total Files</span><span class="stats-value">{{.TotalFiles}}</span>
<span class="stats-label">Total Lines of Code</span><span class="stats-value">{{.TotalLoc}} ({{.TotalAdded}} added, {{.TotalDeleted}} removed)</span>
<span class="stats-label">Total Commits</span><span class="stats-value">{{.TotalCommits}} (avg {{.AvgPerActive}} per active day, {{.AvgPerDay}} per all days)</span>
<span class="stats-label">Authors</span><span class="stats-value">{{.AuthorCount}} (avg {{.AvgPerAuthor}} commits per author)</span>
</div>
</div>
</div>

<div id="activity" class="tab-content">
<div class="section">
<h2>Daily Commits (Last 30 Days)</h2>
<div id="commitChart" class="chart-box"></div>
</div>
<div class="section">
<h2>Daily Line Changes (Last 30 Days)</h2>
<div id="lineChart" class="chart-box"></div>
</div>
<div class="section">
<h2>Hour of Week</h2>
<div id="hourWeekChart" class="chart-box" style="height:340px"></div>
</div>
<div class="section">
<h2>Month of Year</h2>
<div id="monthOfYearChart" class="chart-box" style="height:360px"></div>
</div>
<div class="section">
<h2>Commits by Year/Month</h2>
<div id="yearMonthChart" class="chart-box"></div>
</div>
</div>

<div id="authors" class="tab-content">
<div class="section">
<h2>Author Statistics</h2>
<table>
<thead>
<tr>
<th>Author</th>
<th>Email</th>
<th class="num">Commits</th>
<th>First</th>
<th>Last</th>
<th class="num">Active Days</th>
<th class="num added">++</th>
<th class="num deleted">--</th>
</tr>
</thead>
<tbody>
{{range .Authors}}
<tr>
<td>{{.Name}}</td>
<td>{{.Email}}</td>
<td class="num">{{.CommitCount}}</td>
<td>{{.FirstCommit.Format "2006-01-02"}}</td>
<td>{{.LastCommit.Format "2006-01-02"}}</td>
<td class="num">{{.ActiveDays}}</td>
<td class="num added">{{.AddedLines}}</td>
<td class="num deleted">{{.DeletedLines}}</td>
</tr>
{{end}}
</tbody>
</table>
</div>
<div class="section">
<h2>Daily Line Changes by Author (Last 30 Days)</h2>
<div id="authorLineChart" class="chart-box"></div>
</div>
</div>

</div>

<script>
const colors = ['#5470c6','#91cc75','#fac858','#ee6666','#73c0de','#3ba272','#fc8452','#9a60b4','#ea7ccc',
  '#97b552','#95706d','#dc69aa','#07a2a4','#9a7fd1','#588dd5','#f5994e','#c05050','#59678c','#c9ab00'];

var commitChart = echarts.init(document.getElementById('commitChart'));
commitChart.setOption({{.CommitChartOpt}});

var lineChart = echarts.init(document.getElementById('lineChart'));
lineChart.setOption({{.LineChartOpt}});

var hourWeekChart = echarts.init(document.getElementById('hourWeekChart'));
hourWeekChart.setOption({{.HourWeekOpt}});

var monthOfYearChart = echarts.init(document.getElementById('monthOfYearChart'));
monthOfYearChart.setOption({{.MonthOfYearOpt}});

var yearMonthChart = echarts.init(document.getElementById('yearMonthChart'));
yearMonthChart.setOption({{.YearMonthOpt}});

var authorLineChart = null;
function initAuthorChart() {
  if (authorLineChart) return;
  authorLineChart = echarts.init(document.getElementById('authorLineChart'));
  authorLineChart.setOption({{.AuthorLineChartOpt}});
}

document.querySelectorAll('.tab').forEach(function(tab) {
  tab.addEventListener('click', function() {
    document.querySelectorAll('.tab').forEach(function(t) { t.classList.remove('active'); });
    document.querySelectorAll('.tab-content').forEach(function(c) { c.classList.remove('active'); });
    tab.classList.add('active');
    document.getElementById(tab.dataset.tab).classList.add('active');
    if (tab.dataset.tab === 'authors') initAuthorChart();
    setTimeout(function() {
      commitChart.resize();
      lineChart.resize();
      hourWeekChart.resize();
      monthOfYearChart.resize();
      yearMonthChart.resize();
      if (authorLineChart) authorLineChart.resize();
    }, 100);
  });
});

window.addEventListener('resize', function(){
  commitChart.resize();
  lineChart.resize();
  hourWeekChart.resize();
  monthOfYearChart.resize();
  yearMonthChart.resize();
  if (authorLineChart) authorLineChart.resize();
});
</script>
</body>
</html>`

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
	}

	tmpl, err := template.New("report").Parse(reportTemplate)
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

func buildCommitChartOption(result *AnalysisResult) (string, error) {
	type seriesItem struct {
		Name      string   `json:"name"`
		Type      string   `json:"type"`
		Stack     string   `json:"stack"`
		AreaStyle struct{} `json:"areaStyle"`
		Data      []int    `json:"data"`
	}

	var series []seriesItem
	for _, s := range result.AuthorSeries {
		item := seriesItem{
			Name:  s.Name,
			Type:  "line",
			Stack: "Total",
			Data:  s.Data,
		}
		item.AreaStyle = struct{}{}
		series = append(series, item)
	}

	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"legend": map[string]interface{}{
			"data":   authorNames(result.AuthorSeries),
			"bottom": 0,
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "15%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type":        "category",
			"data":        result.DateRange,
			"boundaryGap": false,
		},
		"yAxis": map[string]interface{}{
			"type": "value",
		},
		"series": series,
	}

	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildLineChartOption(result *AnalysisResult) (string, error) {
	type seriesItem struct {
		Name      string            `json:"name"`
		Type      string            `json:"type"`
		Stack     string            `json:"stack"`
		AreaStyle map[string]string `json:"areaStyle"`
		Data      []int             `json:"data"`
		LineStyle map[string]string `json:"lineStyle,omitempty"`
		ItemStyle map[string]string `json:"itemStyle,omitempty"`
	}

	series := []seriesItem{
		{
			Name: "Added", Type: "line", Stack: "Total",
			AreaStyle: map[string]string{"color": "#4caf50"},
			LineStyle: map[string]string{"color": "#4caf50"},
			ItemStyle: map[string]string{"color": "#4caf50"},
			Data:      result.AddedLineSeries,
		},
		{
			Name: "Deleted", Type: "line", Stack: "Total",
			AreaStyle: map[string]string{"color": "#f44336"},
			LineStyle: map[string]string{"color": "#f44336"},
			ItemStyle: map[string]string{"color": "#f44336"},
			Data:      result.DeletedLineSeries,
		},
	}

	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"legend": map[string]interface{}{
			"data":   []string{"Added", "Deleted"},
			"bottom": 0,
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "15%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category", "data": result.DateRange, "boundaryGap": false,
		},
		"yAxis": map[string]interface{}{"type": "value"},
		"series": series,
	}

	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildAuthorLineChartOption(result *AnalysisResult) (string, error) {
	type seriesItem struct {
		Name      string   `json:"name"`
		Type      string   `json:"type"`
		Stack     string   `json:"stack"`
		AreaStyle struct{} `json:"areaStyle"`
		Data      []int    `json:"data"`
	}

	var series []seriesItem
	for i, s := range result.AuthorAddedSeries {
		total := make([]int, len(s.Data))
		for j := range s.Data {
			total[j] = s.Data[j]
			if j < len(result.AuthorDeletedSeries[i].Data) {
				total[j] += result.AuthorDeletedSeries[i].Data[j]
			}
		}
		item := seriesItem{
			Name:  s.Name,
			Type:  "line",
			Stack: "Total",
			Data:  total,
		}
		item.AreaStyle = struct{}{}
		series = append(series, item)
	}

	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"legend": map[string]interface{}{
			"data":   authorNames(result.AuthorAddedSeries),
			"bottom": 0,
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "15%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category", "data": result.DateRange, "boundaryGap": false,
		},
		"yAxis": map[string]interface{}{"type": "value"},
		"series": series,
	}

	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildHourWeekChartOption(result *AnalysisResult) (string, error) {
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	var hours []string
	for i := 0; i < 24; i++ {
		hours = append(hours, fmt.Sprintf("%02d", i))
	}

	maxVal := 0
	var data [][3]int
	for d := 0; d < 7; d++ {
		for h := 0; h < 24; h++ {
			v := result.HourWeekData[d][h]
			if v > maxVal {
				maxVal = v
			}
			data = append(data, [3]int{h, d, v})
		}
	}

	opt := map[string]interface{}{
		"tooltip": map[string]interface{}{
			"position": "top",
		},
		"grid": map[string]interface{}{
			"left":     "2%",
			"right":    "4%",
			"bottom":   "15%",
			"top":      "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category",
			"data": hours,
			"splitArea": map[string]interface{}{"show": true},
		},
		"yAxis": map[string]interface{}{
			"type": "category",
			"data": days,
			"splitArea": map[string]interface{}{"show": true},
		},
		"visualMap": map[string]interface{}{
			"min":       0,
			"max":       maxVal,
			"calculable": true,
			"orient":    "horizontal",
			"left":      "center",
			"bottom":    0,
			"inRange": map[string]interface{}{
				"color": []string{"#f5f5f5", "#c6e48b", "#7bc96f", "#239a3b", "#196127"},
			},
		},
		"series": []map[string]interface{}{
			{
				"name": "Commits",
				"type": "heatmap",
				"data": data,
				"label": map[string]interface{}{
					"show": false,
				},
				"emphasis": map[string]interface{}{
					"itemStyle": map[string]interface{}{
						"shadowBlur":   10,
						"shadowColor":  "rgba(0,0,0,0.5)",
					},
				},
			},
		},
	}

	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildMonthOfYearChartOption(result *AnalysisResult) (string, error) {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "10%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category", "data": months,
		},
		"yAxis": map[string]interface{}{
			"type": "value",
		},
		"series": []map[string]interface{}{
			{
				"type": "bar",
				"data": result.MonthOfYearData[:],
				"itemStyle": map[string]interface{}{
					"color": "#5470c6",
				},
			},
		},
	}

	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildYearMonthChartOption(result *AnalysisResult) (string, error) {
	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"legend": map[string]interface{}{
			"data":   []string{"Commits"},
			"bottom": 0,
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "15%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category", "data": result.YearMonthLabels,
			"boundaryGap": false,
		},
		"yAxis": map[string]interface{}{
			"type": "value",
		},
		"series": []map[string]interface{}{
			{
				"name": "Commits",
				"type": "line",
				"smooth": true,
				"data": result.YearMonthData,
				"areaStyle": map[string]interface{}{
					"color": "rgba(84,112,198,0.2)",
				},
				"lineStyle": map[string]interface{}{
					"color": "#5470c6",
				},
				"itemStyle": map[string]interface{}{
					"color": "#5470c6",
				},
			},
		},
	}

	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func authorNames(series []AuthorDayData) []string {
	names := make([]string, len(series))
	for i, s := range series {
		names[i] = s.Name
	}
	return names
}
