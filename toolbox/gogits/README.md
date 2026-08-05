# GoGitStats

Analyze git repository commit history and generate a self-contained HTML report with ECharts interactive visualizations.

## Usage

```bash
gogits [-p repo_path] [-f output_file] [-o] [-g] [-h]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-p` | `.` | Git repository path |
| `-f` | auto | Output HTML report path (auto: `.git/gogits-report.html` if `.git` exists, else `./gogits-report.html`) |
| `-o` | false | Open report in Chrome app mode after generation |
| `-g` | false | Print charts to terminal instead of generating HTML (project overview · last 7 days per-author line changes · recent 30 days commit curve) |
| `-h` | false | Print help |

## Report Contents

**General tab** — Project overview: period, age, active days, total commits/files/LOC, per-author averages.

**Activity tab** — 5 charts:
- Daily Commits (Last 30 Days) — stacked area per author
- Daily Line Changes (Last 30 Days) — added vs deleted
- Hour of Week — commit heatmap (7×24)
- Month of Year — commits by calendar month
- Commits by Year/Month — smooth area + data table with show-more toggle

**Authors tab** — Sortable author stats table + 4 charts:
- Author of Month / Author of Year tables
- Daily Line Changes by Author (Last 30 Days)
- Cumulative Commits per Author
- Cumulative Added Lines per Author

**Files tab** — File count & LOC evolution over time (area charts).

## Build

```bash
go build -ldflags "-X main.buildVersion=$(git describe --tags --always)" -o gogits
```