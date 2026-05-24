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
.nowrap { white-space: nowrap; }
.col-author { max-width: 120px; width: 120px; }
.col-commits { min-width: 140px; }
.col-next5 { max-width: 200px; }
th.sortable { cursor: pointer; user-select: none; }
th.sortable:hover { background: #e9ecef; }
th.sortable::after { content: " \25B4\25BE"; font-size: 10px; opacity: .4; }
th.sortable.asc::after { content: " \25B4"; opacity: 1; }
th.sortable.desc::after { content: " \25BE"; opacity: 1; }
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
<button class="tab" data-tab="files">Files</button>
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
<table style="margin-top:16px">
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
<tr>
<td class="num">{{$label}}</td>
<td class="num">{{index $.YearMonthData $i}}</td>
<td class="num added">{{index $.YearMonthAddedData $i}}</td>
<td class="num deleted">{{index $.YearMonthDeletedData $i}}</td>
</tr>
{{end}}
</tbody>
</table>
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
<table>
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
{{range .MonthAuthorStats}}
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
<h2>Lines of Code Over Time</h2>
<div id="locChart" class="chart-box"></div>
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
var cumCommitChart = null;
var cumAddedChart = null;
var fileChart = null;
var locChart = null;
function initAuthorChart() {
  if (authorLineChart) return;
  authorLineChart = echarts.init(document.getElementById('authorLineChart'));
  authorLineChart.setOption({{.AuthorLineChartOpt}});
  cumCommitChart = echarts.init(document.getElementById('cumCommitChart'));
  cumCommitChart.setOption({{.CumCommitOpt}});
  cumAddedChart = echarts.init(document.getElementById('cumAddedChart'));
  cumAddedChart.setOption({{.CumAddedOpt}});
}
function initFileChart() {
  if (fileChart) return;
  fileChart = echarts.init(document.getElementById('fileChart'));
  fileChart.setOption({{.FilesChartOpt}});
}
function initLocChart() {
  if (locChart) return;
  locChart = echarts.init(document.getElementById('locChart'));
  locChart.setOption({{.LocChartOpt}});
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
      if (authorLineChart) authorLineChart.resize();
      if (cumCommitChart) cumCommitChart.resize();
      if (cumAddedChart) cumAddedChart.resize();
      if (fileChart) fileChart.resize();
      if (locChart) locChart.resize();
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
</script>
</body>
</html>`