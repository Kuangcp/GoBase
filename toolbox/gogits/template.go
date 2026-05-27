package main

const reportTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Git Report - {{.RepoName}}</title>
<script src="https://cdn.jsdelivr.net/npm/echarts@5/dist/echarts.min.js"></script>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
:root {
  --bg: #0f0f23;
  --bg-card: #1a1a2e;
  --bg-hover: rgba(255,255,255,.06);
  --text: #e0e0e0;
  --text-secondary: #aaa;
  --text-muted: #777;
  --border: #2a2a4a;
  --shadow: 0 1px 4px rgba(0,0,0,.4);
}
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  background: var(--bg); color: var(--text); }
body.light {
  --bg: #f5f5f5;
  --bg-card: #fff;
  --bg-hover: #f8f9fa;
  --text: #333;
  --text-secondary: #555;
  --text-muted: #888;
  --border: #e9ecef;
  --shadow: 0 1px 3px rgba(0,0,0,.1);
}
.header { background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  color: #fff; padding: 28px 0; }
.header .wrap { max-width: 1200px; margin: 0 auto; padding: 0 20px; }
.header h1 { font-size: 22px; font-weight: 600; }
.header p { font-size: 13px; opacity: .75; margin-top: 4px; }
.wrap { max-width: 1200px; margin: 0 auto; padding: 0 20px; }

.theme-toggle { background: rgba(255,255,255,.15); border: 1px solid rgba(255,255,255,.3);
  color: #fff; padding: 5px 14px; border-radius: 4px; cursor: pointer; font-size: 13px;
  margin-top: 8px; transition: background .2s; }
.theme-toggle:hover { background: rgba(255,255,255,.25); }

.tabs { display: flex; gap: 0; border-bottom: 2px solid var(--border); margin: 20px 0 24px; }
.tab { padding: 10px 28px; cursor: pointer; border: none; background: none;
  font-size: 14px; font-weight: 500; color: var(--text-muted);
  border-bottom: 2px solid transparent; margin-bottom: -2px; transition: all .2s; }
.tab:hover { color: var(--text); }
.tab.active { color: var(--text); border-bottom-color: var(--text); background: var(--bg-hover); }

.tab-content { display: none; }
.tab-content.active { display: block; }

.section { background: var(--bg-card); border-radius: 8px; padding: 24px; margin-bottom: 24px;
  box-shadow: var(--shadow); }
.section h2 { font-size: 17px; margin-bottom: 16px; color: var(--text); }

table { width: 100%; border-collapse: collapse; font-size: 13px; }
th { background: var(--bg-hover); text-align: left; padding: 9px 12px; font-weight: 600;
  border-bottom: 2px solid var(--border); white-space: nowrap; }
td { padding: 9px 12px; border-bottom: 1px solid var(--border); }
tr:hover td { background: var(--bg-hover); }
.num { text-align: right; font-variant-numeric: tabular-nums; }
.nowrap { white-space: nowrap; }
.col-author { max-width: 120px; width: 120px; }
.col-commits { min-width: 140px; }
.col-next5 { max-width: 200px; }
th.sortable { cursor: pointer; user-select: none; }
th.sortable:hover { background: var(--border); }
th.sortable::after { content: " \25B4\25BE"; font-size: 10px; opacity: .4; }
th.sortable.asc::after { content: " \25B4"; opacity: 1; }
th.sortable.desc::after { content: " \25BE"; opacity: 1; }
.added { color: #4caf50; font-weight: 600; }
.deleted { color: #f44336; font-weight: 600; }
.chart-box { width: 100%; height: 420px; }
.extra-row.hidden { display: none; }
.collapse-toggle { display: block; width: 100%; padding: 8px; text-align: center;
  background: var(--bg-hover); border: 1px solid var(--border); border-radius: 4px;
  cursor: pointer; font-size: 13px; color: #5470c6; margin-top: 8px; }
.collapse-toggle:hover { background: var(--border); }

.stats-grid { display: grid; grid-template-columns: auto 1fr; gap: 2px 32px;
  font-size: 14px; line-height: 2; }
.stats-label { color: var(--text-muted); white-space: nowrap; text-align: right; }
.stats-value { color: var(--text); font-weight: 500; }

.section-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; margin-bottom: 24px; }
.section-grid .section { margin-bottom: 0; }

.overall-badge { display: flex; flex-direction: column; align-items: center;
  padding: 10px 20px; border-radius: 10px; min-width: 100px;
  background: var(--bg-card); border: 2px solid var(--border);
  box-shadow: var(--shadow); }
.overall-label { font-size: 10px; opacity: .7; margin-bottom: 1px;
  text-transform: uppercase; letter-spacing: 1px; }
.overall-badge .badge-grade { font-size: 28px; }
.overall-badge .badge-score { font-size: 10px; }
.overall-badge.overall-S { border-color: #ffd700; background: linear-gradient(135deg, rgba(255,215,0,0.08), transparent); }
.overall-badge.overall-A { border-color: #4caf50; background: linear-gradient(135deg, rgba(76,175,80,0.08), transparent); }
.overall-badge.overall-B { border-color: #5470c6; background: linear-gradient(135deg, rgba(84,112,198,0.08), transparent); }
.overall-badge.overall-C { border-color: #ff9800; background: linear-gradient(135deg, rgba(255,152,0,0.08), transparent); }
.overall-badge.overall-D { border-color: #f44336; background: linear-gradient(135deg, rgba(244,67,54,0.08), transparent); }
.overall-badge.overall-E { border-color: #9e9e9e; background: linear-gradient(135deg, rgba(158,158,158,0.08), transparent); }

.radar-grid { display: grid; grid-template-columns: auto 1fr 1fr; gap: 24px; align-items: center; }
.radar-stats { display: flex; flex-direction: column; gap: 10px; }
.radar-item { display: flex; align-items: center; gap: 8px; font-size: 14px; }
.radar-item .badge-dot { flex-shrink: 0; }
.radar-item .radar-name { flex: 1; }
.radar-item .radar-grade { width: 44px; font-weight: 600; white-space: nowrap; }
.radar-item .radar-score { width: 56px; text-align: right; font-weight: 500; white-space: nowrap; }
.badge { display: flex; flex-direction: column; align-items: center; padding: 16px 20px;
  border-radius: 12px; min-width: 110px; background: var(--bg-card);
  border: 2px solid var(--border); box-shadow: var(--shadow); transition: transform .2s; }
.badge:hover { transform: translateY(-3px); }
.badge-grade { font-size: 38px; font-weight: 800; line-height: 1; }
.badge-label { font-size: 13px; margin-top: 6px; opacity: .85; }
.badge-score { font-size: 11px; margin-top: 3px; opacity: .65; }
.badge.badge-S { border-color: #ffd700; }
.badge-S .badge-grade { color: #ffd700; }
.badge.badge-A { border-color: #4caf50; }
.badge-A .badge-grade { color: #4caf50; }
.badge.badge-B { border-color: #5470c6; }
.badge-B .badge-grade { color: #5470c6; }
.badge.badge-C { border-color: #ff9800; }
.badge-C .badge-grade { color: #ff9800; }
.badge.badge-D { border-color: #f44336; }
.badge-D .badge-grade { color: #f44336; }
.badge.badge-E { border-color: #9e9e9e; }
.badge-E .badge-grade { color: #9e9e9e; }
.badge-dot { display: inline-block; width: 10px; height: 10px; border-radius: 50%;
  margin-right: 6px; vertical-align: middle; }
.badge-dot-S { background: #ffd700; }
.badge-dot-A { background: #4caf50; }
.badge-dot-B { background: #5470c6; }
.badge-dot-C { background: #ff9800; }
.badge-dot-D { background: #f44336; }
.badge-dot-E { background: #9e9e9e; }
.section-title { display: flex; align-items: center; font-size: 17px; margin-bottom: 16px; color: var(--text); }

@media (max-width: 640px) {
  .chart-box { height: 280px; }
  table { font-size: 11px; }
  th, td { padding: 5px 6px; }
  .stats-grid { font-size: 12px; gap: 2px 16px; }
  .tab { padding: 8px 14px; font-size: 13px; }
  .badge { min-width: 80px; padding: 12px 14px; }
  .badge-grade { font-size: 28px; }
  .section-grid { grid-template-columns: 1fr; }
  .radar-grid { grid-template-columns: 1fr; }
}
</style>
</head>
<body>

<div class="header">
<div class="wrap">
<h1>{{.RepoName}}</h1>
<p>{{.RepoPath}} &nbsp;·&nbsp; {{.Branch}} branch &nbsp;·&nbsp; Generated {{.GeneratedAt}}</p>
<button class="theme-toggle" onclick="toggleTheme()">Light Mode</button>
</div>
</div>

<div class="wrap">
<div class="tabs">
<button class="tab active" data-tab="general">General</button>
<button class="tab" data-tab="activity">Activity</button>
<button class="tab" data-tab="authors">Authors</button>
<button class="tab" data-tab="files">Files</button>
</div>

<div id="general" class="tab-content active">

<div class="section">
<div class="radar-grid">
<div class="overall-badge overall-{{.OverallGrade}}">
<span class="overall-label">综合评级</span>
<span class="badge-grade">{{.OverallGrade}}</span>
<span class="badge-score">{{.OverallScore}} 分</span>
</div>
<div><div id="radarChart" class="chart-box" style="height:300px"></div></div>
<div>
<h2 style="margin-bottom:16px;font-size:15px">各维度评分</h2>
<div class="radar-stats">
<div class="radar-item"><span class="badge-dot badge-dot-{{.ActivityGrade}}"></span><span class="radar-name">活跃度</span><span class="radar-grade">{{.ActivityGrade}} 级</span><span class="radar-score">{{.ActivityScore}} 分</span></div>
<div class="radar-item"><span class="badge-dot badge-dot-{{.ScaleGrade}}"></span><span class="radar-name">项目规模</span><span class="radar-grade">{{.ScaleGrade}} 级</span><span class="radar-score">{{.ScaleScore}} 分</span></div>
<div class="radar-item"><span class="badge-dot badge-dot-{{.HealthGrade}}"></span><span class="radar-name">代码健康</span><span class="radar-grade">{{.HealthGrade}} 级</span><span class="radar-score">{{.HealthScore}} 分</span></div>
<div class="radar-item"><span class="badge-dot badge-dot-{{.DiversityGrade}}"></span><span class="radar-name">协作多样</span><span class="radar-grade">{{.DiversityGrade}} 级</span><span class="radar-score">{{.DiversityScore}} 分</span></div>
<div class="radar-item"><span class="badge-dot badge-dot-{{.DebtGrade}}"></span><span class="radar-name">技术债</span><span class="radar-grade">{{.DebtGrade}} 级</span><span class="radar-score">{{.DebtScore}} 分</span></div>
<div class="radar-item"><span class="badge-dot badge-dot-{{.RhythmGrade}}"></span><span class="radar-name">研发节奏</span><span class="radar-grade">{{.RhythmGrade}} 级</span><span class="radar-score">{{.RhythmScore}} 分</span></div>
</div>
</div>
</div>
</div>

<div class="section-grid">
<div class="section">
<div class="section-title"><span class="badge-dot badge-dot-{{.ActivityGrade}}"></span>活跃度 · {{.ActivityGrade}} 级 ({{.ActivityScore}} 分)</div>
<div class="stats-grid">
<span class="stats-label">总提交数</span><span class="stats-value">{{.TotalCommits}}</span>
<span class="stats-label">近30天提交</span><span class="stats-value">{{.Commits30d}}</span>
<span class="stats-label">近30天活跃开发者</span><span class="stats-value">{{.ActiveDevs30d}} 人</span>
<span class="stats-label">近30天活跃天数</span><span class="stats-value">{{.ActiveDays30d}} 天</span>
<span class="stats-label">提交强度</span><span class="stats-value">{{.AvgPerActive}} 次/活跃日 · {{.AvgPerDay}} 次/日历日</span>
<span class="stats-label">近期动量</span><span class="stats-value">近30天 {{.RecentMonthCommits}} 次提交 ({{.RecentMomentumPct}}%)</span>
</div>
</div>
<div class="section">
<div class="section-title"><span class="badge-dot badge-dot-{{.ScaleGrade}}"></span>项目规模 · {{.ScaleGrade}} 级 ({{.ScaleScore}} 分)</div>
<div class="stats-grid">
<span class="stats-label">代码行数</span><span class="stats-value">{{.TotalLoc}}</span>
<span class="stats-label">文件总数</span><span class="stats-value">{{.TotalFiles}}</span>
<span class="stats-label">贡献者总数</span><span class="stats-value">{{.AuthorCount}} 人</span>
<span class="stats-label">人均提交数</span><span class="stats-value">{{.AvgPerAuthor}} 次/人</span>
<span class="stats-label">工作量密度</span><span class="stats-value">{{.WorkloadDensity}} LOC/提交</span>
<span class="stats-label">Bus Factor</span><span class="stats-value">{{.BusFactorCount}} 位作者控制 {{.BusFactorPct}}% 的提交</span>
</div>
</div>
</div>

<div class="section-grid">
<div class="section">
<div class="section-title"><span class="badge-dot badge-dot-{{.HealthGrade}}"></span>代码健康 · {{.HealthGrade}} 级 ({{.HealthScore}} 分)</div>
<div class="stats-grid">
<span class="stats-label">大文件数 (>1000行)</span><span class="stats-value">{{.LargeFileCount}} 个</span>
<span class="stats-label">TODO/FIXME 标记</span><span class="stats-value">{{.TodoCount}} 处</span>
<span class="stats-label">新增/删除比</span><span class="stats-value">{{.TotalAdded}} / {{.TotalDeleted}} ({{.ChurnRatio}}:1)</span>
<span class="stats-label">活跃天数占比</span><span class="stats-value">{{.ActivePct}}% ({{.TotalActiveDays}}/{{.AgeDays}} 天有提交)</span>
</div>
</div>
<div class="section">
<div class="section-title"><span class="badge-dot badge-dot-{{.DebtGrade}}"></span>技术债 · {{.DebtGrade}} 级 ({{.DebtScore}} 分)</div>
<div class="stats-grid">
<span class="stats-label">热点文件 (近90天)</span><span class="stats-value">{{.HotspotCount}} 个</span>
<span class="stats-label">遗弃代码 (1年未动)</span><span class="stats-value">{{.AbandonedPct}}%</span>
<span class="stats-label">代码平均年龄</span><span class="stats-value">{{.CodeAgeDays}} 天</span>
</div>
<div style="margin-top:12px;font-size:12px;color:var(--text-muted)">修改最频繁的文件 Top {{.HotspotCount}}</div>
<div style="margin-top:6px;font-size:13px;line-height:2">
{{range .Hotspots}}<div style="display:flex;gap:8px"><span style="color:var(--text-muted);min-width:36px;text-align:right;font-weight:600">({{.ModifyCount}})</span><span>{{.Path}}</span></div>{{end}}
</div>
</div>
</div>

<div class="section-grid">
<div class="section">
<div class="section-title"><span class="badge-dot badge-dot-{{.RhythmGrade}}"></span>研发节奏 · {{.RhythmGrade}} 级 ({{.RhythmScore}} 分)</div>
<div class="stats-grid">
<span class="stats-label">提交连续性</span><span class="stats-value">{{.ConsistencyPct}}% ({{.ActiveWeeks}}/{{.TotalWeeks}} 周有提交)</span>
<span class="stats-label">版本发布数</span><span class="stats-value">{{.ReleaseCount}} 个标签</span>
<span class="stats-label">加班疲劳度</span><span class="stats-value">{{.OffHoursPct}}% (非工作时间提交)</span>
<span class="stats-label">项目持续天数</span><span class="stats-value">{{.AgeDays}} 天</span>
</div>
</div>
<div class="section">
<div class="section-title"><span class="badge-dot badge-dot-{{.DiversityGrade}}"></span>协作多样 · {{.DiversityGrade}} 级 ({{.DiversityScore}} 分)</div>
<div class="stats-grid">
<span class="stats-label">贡献者分布</span><span class="stats-value">{{.BusFactorCount}} 位核心作者控制 {{.BusFactorPct}}% 的提交</span>
<span class="stats-label">团队纵深</span><span class="stats-value">{{sub .AuthorCount .BusFactorCount}} 人超出 Bus Factor 核心</span>
<span class="stats-label">总贡献者</span><span class="stats-value">{{.AuthorCount}} 人</span>
<span class="stats-label">人均提交数</span><span class="stats-value">{{.AvgPerAuthor}} 次/人</span>
</div>
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
<table style="margin-top:16px" id="yearMonthTable">
<thead>
<tr>
<th class="num">Month</th>
<th class="num">Commits</th>
<th class="num added">++</th>
<th class="num deleted">--</th>
</tr>
</thead>
<tbody>
{{range $i, $label := .YearMonthLabels}}
<tr class="{{if ge $i 15}}extra-row hidden{{end}}">
<td class="num">{{$label}}</td>
<td class="num">{{index $.YearMonthData $i}}</td>
<td class="num added">{{index $.YearMonthAddedData $i}}</td>
<td class="num deleted">{{index $.YearMonthDeletedData $i}}</td>
</tr>
{{end}}
</tbody>
</table>
{{if gt (len .YearMonthLabels) 15}}
<button class="collapse-toggle" onclick="toggleCollapse(this)" data-table="yearMonthTable">Show more ({{sub (len .YearMonthLabels) 15}})</button>
{{end}}
</div>
</div>

<div id="authors" class="tab-content">
<div class="section">
<h2>Author Statistics</h2>
<table id="authorStatsTable">
<thead>
<tr>
<th data-sort="name">Author</th>
<th>Email</th>
<th class="num" data-sort="commits">Commits</th>
<th class="num" data-sort="age">Age</th>
<th data-sort="first">First</th>
<th data-sort="last">Last</th>
<th class="num" data-sort="active">Active Days</th>
<th class="num added" data-sort="added">++</th>
<th class="num deleted" data-sort="deleted">--</th>
</tr>
</thead>
<tbody>
{{range .Authors}}
<tr>
<td>{{.Name}}</td>
<td>{{.Email}}</td>
<td class="num">{{.CommitCount}}</td>
<td class="num">{{age .FirstCommit .LastCommit}}</td>
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
<h2>Author of Month</h2>
<table id="monthAuthorTable">
<thead>
<tr>
<th>Month</th>
<th class="col-author">Author</th>
<th class="num col-commits nowrap">Commits (%)</th>
<th class="col-next5">Next top 5</th>
<th class="num">Number of authors</th>
</tr>
</thead>
<tbody>
{{range $i, $stat := .MonthAuthorStats}}
<tr class="{{if ge $i 15}}extra-row hidden{{end}}">
<td class="nowrap">{{$stat.Period}}</td>
<td class="col-author">{{$stat.TopAuthor}}</td>
<td class="num col-commits nowrap">{{printf "%4d" $stat.TopCommits}}/{{printf "%-4d" $stat.TotalCommits}} ({{printf "%.2f" (percent $stat.TopCommits $stat.TotalCommits)}}%)</td>
<td class="col-next5">{{join $stat.NextTop5 ", "}}</td>
<td class="num">{{$stat.AuthorCount}}</td>
</tr>
{{end}}
</tbody>
</table>
{{if gt (len .MonthAuthorStats) 15}}
<button class="collapse-toggle" onclick="toggleCollapse(this)" data-table="monthAuthorTable">Show more ({{sub (len .MonthAuthorStats) 15}})</button>
{{end}}
</div>
<div class="section">
<h2>Author of Year</h2>
<table>
<thead>
<tr>
<th>Year</th>
<th class="col-author">Author</th>
<th class="num col-commits nowrap">Commits (%)</th>
<th class="col-next5">Next top 5</th>
<th class="num">Number of authors</th>
</tr>
</thead>
<tbody>
{{range .YearAuthorStats}}
<tr>
<td class="nowrap">{{.Period}}</td>
<td class="col-author">{{.TopAuthor}}</td>
<td class="num col-commits nowrap">{{printf "%4d" .TopCommits}}/{{printf "%-4d" .TotalCommits}} ({{printf "%.2f" (percent .TopCommits .TotalCommits)}}%)</td>
<td class="col-next5">{{join .NextTop5 ", "}}</td>
<td class="num">{{.AuthorCount}}</td>
</tr>
{{end}}
</tbody>
</table>
</div>
<div class="section">
<h2>Daily Line Changes by Author (Last 30 Days)</h2>
<div id="authorLineChart" class="chart-box"></div>
</div>
<div class="section">
<h2>Commits per Author (Cumulative)</h2>
<div id="cumCommitChart" class="chart-box"></div>
</div>
<div class="section">
<h2>Cumulated Added Lines of Code per Author</h2>
<div id="cumAddedChart" class="chart-box"></div>
</div>
</div>

<div id="files" class="tab-content">
<div class="section">
<h2>File Statistics</h2>
<div class="stats-grid">
<span class="stats-label">Total files</span><span class="stats-value">{{.TotalFiles}}</span>
<span class="stats-label">Total lines</span><span class="stats-value">{{.TotalLoc}}</span>
<span class="stats-label">Average file size</span><span class="stats-value">{{printf "%.2f" (avgFileSize .TotalLoc .TotalFiles)}} bytes</span>
</div>
</div>
<div class="section">
<h2>Files Count Over Time</h2>
<div id="fileChart" class="chart-box"></div>
</div>
<div class="section">
<h2>Extensions</h2>
<table>
<thead>
<tr>
<th>Extension</th>
<th class="num">Files (%)</th>
<th class="num">Lines (%)</th>
<th class="num">Lines/file</th>
</tr>
</thead>
<tbody>
{{range .ExtensionStats}}
<tr>
<td>{{if eq .Extension ""}}(empty){{else}}.{{.Extension}}{{end}}</td>
<td class="num">{{.FileCount}} ({{printf "%.2f" (percent .FileCount $.TotalFiles)}}%)</td>
<td class="num">{{.LineCount}} ({{printf "%.2f" (percent .LineCount $.TotalLoc)}}%)</td>
<td class="num">{{printf "%.0f" (avgFileSize .LineCount .FileCount)}}</td>
</tr>
{{end}}
</tbody>
</table>
</div>
<div class="section">
<h2>Lines of Code Over Time</h2>
<div id="locChart" class="chart-box"></div>
</div>
</div>

</div>

<script>
const colors = ['#5470c6','#91cc75','#fac858','#ee6666','#73c0de','#3ba272','#fc8452','#9a60b4','#ea7ccc',
  '#97b552','#95706d','#dc69aa','#07a2a4','#9a7fd1','#588dd5','#f5994e','#c05050','#59678c','#c9ab00'];

var allCharts = [];

// restore saved theme before creating charts
(function() {
  if (localStorage.getItem('theme') === 'light') {
    document.body.classList.add('light');
    var btn = document.querySelector('.theme-toggle');
    if (btn) btn.textContent = 'Dark Mode';
  }
})();

function addChart(el, opt) {
  var c = echarts.init(el);
  c.setOption(opt);
  allCharts.push(c);
  return c;
}

var commitChart = addChart(document.getElementById('commitChart'), {{.CommitChartOpt}});
var lineChart = addChart(document.getElementById('lineChart'), {{.LineChartOpt}});
var hourWeekChart = addChart(document.getElementById('hourWeekChart'), {{.HourWeekOpt}});
var monthOfYearChart = addChart(document.getElementById('monthOfYearChart'), {{.MonthOfYearOpt}});
var yearMonthChart = addChart(document.getElementById('yearMonthChart'), {{.YearMonthOpt}});
var radarChart = addChart(document.getElementById('radarChart'), {{.RadarChartOpt}});

// apply theme to initial charts
updateChartTheme();

function chartAxisColor(isLight) {
  return isLight ? '#333' : '#e0e0e0';
}
function updateChartTheme() {
  var color = chartAxisColor(document.body.classList.contains('light'));
  allCharts.forEach(function(c) {
    var opt = c.getOption();
    if (!opt.textStyle) opt.textStyle = {};
    opt.textStyle.color = color;
    if (!opt.xAxis) opt.xAxis = {};
    if (!opt.yAxis) opt.yAxis = {};
    [opt.xAxis, opt.yAxis].forEach(function(axis) {
      if (!axis) return;
      var list = Array.isArray(axis) ? axis : [axis];
      list.forEach(function(a) {
        a.axisLabel = a.axisLabel || {};
        a.axisLabel.color = color;
      });
    });
    if (opt.legend) {
      var legs = Array.isArray(opt.legend) ? opt.legend : [opt.legend];
      legs.forEach(function(l) {
        l.textStyle = l.textStyle || {};
        l.textStyle.color = color;
      });
    }
    if (opt.radar) {
      var radar = Array.isArray(opt.radar) ? opt.radar : [opt.radar];
      radar.forEach(function(r) {
        if (r.axisName) {
          r.axisName.textStyle = r.axisName.textStyle || {};
          r.axisName.textStyle.color = color;
        }
      });
    }
    c.setOption(opt, { notMerge: true });
  });
}

var authorLineChart = null;
var cumCommitChart = null;
var cumAddedChart = null;
var fileChart = null;
var locChart = null;
function initAuthorChart() {
  if (authorLineChart) return;
  authorLineChart = echarts.init(document.getElementById('authorLineChart'));
  authorLineChart.setOption({{.AuthorLineChartOpt}});
  allCharts.push(authorLineChart);
  cumCommitChart = echarts.init(document.getElementById('cumCommitChart'));
  cumCommitChart.setOption({{.CumCommitOpt}});
  allCharts.push(cumCommitChart);
  cumAddedChart = echarts.init(document.getElementById('cumAddedChart'));
  cumAddedChart.setOption({{.CumAddedOpt}});
  allCharts.push(cumAddedChart);
  updateChartTheme();
}
function initFileChart() {
  if (fileChart) return;
  fileChart = echarts.init(document.getElementById('fileChart'));
  fileChart.setOption({{.FilesChartOpt}});
  allCharts.push(fileChart);
  updateChartTheme();
}
function initLocChart() {
  if (locChart) return;
  locChart = echarts.init(document.getElementById('locChart'));
  locChart.setOption({{.LocChartOpt}});
  allCharts.push(locChart);
  updateChartTheme();
}

document.querySelectorAll('.tab').forEach(function(tab) {
  tab.addEventListener('click', function() {
    document.querySelectorAll('.tab').forEach(function(t) { t.classList.remove('active'); });
    document.querySelectorAll('.tab-content').forEach(function(c) { c.classList.remove('active'); });
    tab.classList.add('active');
    document.getElementById(tab.dataset.tab).classList.add('active');
    if (tab.dataset.tab === 'authors') initAuthorChart();
    if (tab.dataset.tab === 'files') { initFileChart(); initLocChart(); }
    setTimeout(function() {
      commitChart.resize();
      lineChart.resize();
      hourWeekChart.resize();
      monthOfYearChart.resize();
      yearMonthChart.resize();
      radarChart.resize();
      if (authorLineChart) authorLineChart.resize();
      if (cumCommitChart) cumCommitChart.resize();
      if (cumAddedChart) cumAddedChart.resize();
      if (fileChart) fileChart.resize();
      if (locChart) locChart.resize();
    }, 100);
  });
});

function toggleCollapse(btn) {
  var rows = document.getElementById(btn.getAttribute('data-table')).querySelectorAll('.extra-row');
  var expanded = btn.getAttribute('data-expanded') === 'true';
  rows.forEach(function(r) { r.classList.toggle('hidden', expanded); });
  btn.textContent = expanded ? 'Show more (' + rows.length + ')' : 'Collapse';
  btn.setAttribute('data-expanded', expanded ? 'false' : 'true');
}

window.addEventListener('resize', function(){
  commitChart.resize();
  lineChart.resize();
  hourWeekChart.resize();
  monthOfYearChart.resize();
  yearMonthChart.resize();
  radarChart.resize();
  if (authorLineChart) authorLineChart.resize();
  if (cumCommitChart) cumCommitChart.resize();
  if (cumAddedChart) cumAddedChart.resize();
  if (fileChart) fileChart.resize();
  if (locChart) locChart.resize();
});

(function() {
  var table = document.getElementById('authorStatsTable');
  if (!table) return;
  var tbody = table.querySelector('tbody');
  var allThs = table.querySelectorAll('thead th');
  var sortThs = table.querySelectorAll('thead th[data-sort]');
  sortThs.forEach(function(th) {
    th.classList.add('sortable');
    th.addEventListener('click', function() {
      var colIdx = Array.prototype.indexOf.call(allThs, th);
      var key = th.getAttribute('data-sort');
      var isAsc = th.classList.contains('asc');
      sortThs.forEach(function(t) { t.classList.remove('asc', 'desc'); });
      th.classList.add(isAsc ? 'desc' : 'asc');
      var rows = Array.prototype.slice.call(tbody.querySelectorAll('tr'));
      rows.sort(function(a, b) {
        var va = a.querySelectorAll('td')[colIdx].textContent.trim();
        var vb = b.querySelectorAll('td')[colIdx].textContent.trim();
        if (key === 'commits' || key === 'active' || key === 'added' || key === 'deleted') {
          return (parseInt(va) || 0) - (parseInt(vb) || 0);
        }
        if (key === 'age') {
          return parseAge(va) - parseAge(vb);
        }
        return va.localeCompare(vb);
      });
      if (isAsc) rows.reverse();
      rows.forEach(function(r) { tbody.appendChild(r); });
    });
  });
  function parseAge(s) {
    var m = s.match(/(\d+) days?, (\d+):(\d+):(\d+)/);
    if (!m) return 0;
    return parseInt(m[1])*86400 + parseInt(m[2])*3600 + parseInt(m[3])*60 + parseInt(m[4]);
  }
})();

function toggleTheme() {
  document.body.classList.toggle('light');
  var btn = document.querySelector('.theme-toggle');
  btn.textContent = document.body.classList.contains('light') ? 'Dark Mode' : 'Light Mode';
  localStorage.setItem('theme', document.body.classList.contains('light') ? 'light' : 'dark');
  updateChartTheme();
}
</script>
</body>
</html>`