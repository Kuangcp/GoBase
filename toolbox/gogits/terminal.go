package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/kuangcp/gobase/pkg/ctool"
)

// displayWidth returns the terminal display width of s, counting CJK as 2 columns.
func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		if r > 0x2E7F {
			w += 2
		} else {
			w++
		}
	}
	return w
}

func padRight(s string, n int) string {
	if d := n - displayWidth(s); d > 0 {
		return s + strings.Repeat(" ", d)
	}
	return s
}

func padLeft(s string, n int) string {
	if d := n - displayWidth(s); d > 0 {
		return strings.Repeat(" ", d) + s
	}
	return s
}

func center(s string, n int) string {
	d := n - displayWidth(s)
	if d <= 0 {
		return s
	}
	l := d / 2
	return strings.Repeat(" ", l) + s + strings.Repeat(" ", d-l)
}

func formatThousands(n int) string {
	neg := n < 0
	s := fmt.Sprintf("%d", n)
	if neg {
		s = s[1:]
	}
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	r := strings.Join(parts, ",")
	if neg {
		r = "-" + r
	}
	return r
}

func colorize(c ctool.Color, s string) string {
	if c == "" {
		return s
	}
	return c.Print(s)
}

// RenderProjectOverview builds the terminal project overview box.
func RenderProjectOverview(result *AnalysisResult) string {
	recentMonth := 0
	for _, s := range result.AuthorSeries {
		for _, v := range s.Data {
			recentMonth += v
		}
	}
	recentPct := 0.0
	if result.TotalCommits > 0 {
		recentPct = float64(recentMonth) / float64(result.TotalCommits) * 100
	}

	var names []string
	for _, a := range result.Authors {
		names = append(names, a.Name)
	}

	dur := result.GenerationDuration
	durStr := fmt.Sprintf("%.1fs", dur.Seconds())
	if dur.Seconds() < 1 {
		durStr = fmt.Sprintf("%dms", dur.Milliseconds())
	}

	type rowT struct {
		label string
		value string
		color ctool.Color
		wrap  bool
	}
	participants := strings.Join(names, ", ")
	rows := []rowT{
		{"仓库", result.RepoName + " (" + result.Branch + ")", "", false},
		{"总提交数", formatThousands(result.TotalCommits), ctool.LightBlue, false},
		{"近30天提交", fmt.Sprintf("%s (%.1f%%)", formatThousands(recentMonth), recentPct), ctool.LightYellow, false},
		{"参与人", participants, ctool.LightCyan, true},
		{"生成时间", time.Now().Format("2006-01-02 15:04:05") + " (" + durStr + ")", "", false},
	}

	inner := 0
	for _, r := range rows {
		if r.wrap {
			continue
		}
		if w := displayWidth(r.label) + displayWidth(r.value) + 2; w > inner {
			inner = w
		}
	}
	if inner < displayWidth("项目概览 Project Overview") {
		inner = displayWidth("项目概览 Project Overview")
	}

	var b strings.Builder
	horiz := strings.Repeat("═", inner+2)
	b.WriteString("╔" + horiz + "╗\n")
	b.WriteString("║ " + center("项目概览 Project Overview", inner) + " ║\n")
	b.WriteString("╠" + horiz + "╣\n")
	for _, r := range rows {
		if r.wrap {
			for j, line := range wrapByWidth(r.value, inner-displayWidth(r.label)) {
				label := r.label
				if j > 0 {
					label = strings.Repeat(" ", displayWidth(r.label))
				}
				pad := inner - displayWidth(label) - displayWidth(line)
				if pad < 0 {
					pad = 0
				}
				b.WriteString("║ " + label + strings.Repeat(" ", pad) + colorize(r.color, line) + " ║\n")
			}
			continue
		}
		pad := inner - displayWidth(r.label) - displayWidth(r.value)
		if pad < 0 {
			pad = 0
		}
		b.WriteString("║ " + r.label + strings.Repeat(" ", pad) + colorize(r.color, r.value) + " ║\n")
	}
	b.WriteString("╚" + horiz + "╝")
	return b.String()
}

// wrapByWidth splits s into lines, each fitting within maxW display columns.
// It prefers to break on ", " so names stay whole, and hard-wraps any single
// segment that is still too wide.
func wrapByWidth(s string, maxW int) []string {
	if maxW <= 0 {
		return []string{s}
	}
	segments := strings.Split(s, ", ")
	var lines []string
	cur := ""
	for _, seg := range segments {
		sep := ""
		if cur != "" {
			sep = ", "
		}
		if displayWidth(cur)+displayWidth(sep)+displayWidth(seg) <= maxW {
			cur += sep + seg
			continue
		}
		if cur != "" {
			lines = append(lines, cur)
			cur = ""
		}
		if displayWidth(seg) > maxW {
			runes := []rune(seg)
			start := 0
			w := 0
			for i, r := range runes {
				cw := 1
				if r > 0x2E7F {
					cw = 2
				}
				if w+cw > maxW && i > start {
					lines = append(lines, string(runes[start:i]))
					start = i
					w = 0
				}
				w += cw
			}
			if start < len(runes) {
				lines = append(lines, string(runes[start:]))
			}
		} else {
			cur = seg
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}

type weekEntry struct {
	name     string
	commits  int
	added    int
	deleted  int
}

type weekDay struct {
	date    string
	entries []weekEntry
}

func netStr(n int) string {
	if n == 0 {
		return "0"
	}
	return fmt.Sprintf("%+d", n)
}

func borderRow(widths []int, left, mid, right string) string {
	var b strings.Builder
	b.WriteString(left)
	for i, w := range widths {
		b.WriteString(strings.Repeat("─", w+2))
		if i < len(widths)-1 {
			b.WriteString(mid)
		}
	}
	b.WriteString(right)
	return b.String()
}

// RenderAuthorWeekTable builds the last-7-days per-author change table.
func RenderAuthorWeekTable(result *AnalysisResult) string {
	n := len(result.DateRange)
	if n == 0 {
		return ""
	}
	startIdx := n - 7
	if startIdx < 0 {
		startIdx = 0
	}

	weekdays := []string{"日", "一", "二", "三", "四", "五", "六"}
	var days []weekDay
	for i := startIdx; i < n; i++ {
		label := result.DateRange[i]
		t, _ := time.Parse("2006-01-02", label)
		date := t.Format("01-02") + " 周" + weekdays[int(t.Weekday())]
		wd := weekDay{date: date}
		for si := range result.AuthorAddedSeries {
			if i >= len(result.AuthorAddedSeries[si].Data) {
				continue
			}
			commits := 0
			if si < len(result.AuthorSeries) && i < len(result.AuthorSeries[si].Data) {
				commits = result.AuthorSeries[si].Data[i]
			}
			added := result.AuthorAddedSeries[si].Data[i]
			deleted := 0
			if si < len(result.AuthorDeletedSeries) && i < len(result.AuthorDeletedSeries[si].Data) {
				deleted = result.AuthorDeletedSeries[si].Data[i]
			}
			if commits > 0 || added > 0 || deleted > 0 {
				wd.entries = append(wd.entries, weekEntry{name: result.AuthorAddedSeries[si].Name, commits: commits, added: added, deleted: deleted})
			}
		}
		days = append(days, wd)
	}

	dateW := displayWidth("日期")
	authorW := displayWidth("人员")
	numW := displayWidth("净变更")
	for _, d := range days {
		if w := displayWidth(d.date); w > dateW {
			dateW = w
		}
		for _, e := range d.entries {
			if w := displayWidth(e.name); w > authorW {
				authorW = w
			}
			for _, s := range []string{fmt.Sprintf("%d", e.commits), fmt.Sprintf("%d", e.added), fmt.Sprintf("%d", e.deleted), netStr(e.added - e.deleted)} {
				if w := displayWidth(s); w > numW {
					numW = w
				}
			}
		}
	}
	if w := displayWidth("当日合计"); w > authorW {
		authorW = w
	}

	widths := []int{dateW, authorW, numW, numW, numW, numW}
	headers := []string{"日期", "人员", "提交数", "新增", "删除", "净变更"}

	type rowT struct {
		cells []string
		sep   bool
	}
	var rows []rowT

	headerCells := make([]string, len(headers))
	for i, h := range headers {
		headerCells[i] = colorize(ctool.LightWhite, padRight(h, widths[i]))
	}
	rows = append(rows, rowT{cells: headerCells, sep: true})

	for _, d := range days {
		if len(d.entries) == 0 {
			cells := []string{padRight(d.date, dateW), padRight("无提交", authorW), "", "", "", ""}
			rows = append(rows, rowT{cells: cells, sep: true})
			continue
		}
		dayAdded, dayDeleted, dayCommits := 0, 0, 0
		firstIdx := len(rows)
		for _, e := range d.entries {
			dayCommits += e.commits
			dayAdded += e.added
			dayDeleted += e.deleted
			net := e.added - e.deleted
			netColor := ctool.Green
			if net < 0 {
				netColor = ctool.Red
			}
			cells := []string{
				"",
				padRight(e.name, authorW),
				colorize(ctool.LightCyan, padLeft(fmt.Sprintf("%d", e.commits), numW)),
				colorize(ctool.Green, padLeft(fmt.Sprintf("%d", e.added), numW)),
				colorize(ctool.Red, padLeft(fmt.Sprintf("%d", e.deleted), numW)),
				colorize(netColor, padLeft(netStr(net), numW)),
			}
			rows = append(rows, rowT{cells: cells})
		}
		rows[firstIdx].cells[0] = padRight(d.date, dateW)

		subCells := []string{
			"",
			colorize(ctool.LightYellow, padRight("当日合计", authorW)),
			colorize(ctool.LightCyan, padLeft(fmt.Sprintf("%d", dayCommits), numW)),
			colorize(ctool.Green, padLeft(fmt.Sprintf("%d", dayAdded), numW)),
			colorize(ctool.Red, padLeft(fmt.Sprintf("%d", dayDeleted), numW)),
			colorize(ctool.LightYellow, padLeft(netStr(dayAdded-dayDeleted), numW)),
		}
		rows = append(rows, rowT{cells: subCells, sep: true})
	}

	title := "最近7天代码变更 (" + result.DateRange[startIdx] + " ~ " + result.DateRange[n-1] + ")"
	var b strings.Builder
	b.WriteString(title + "\n")
	b.WriteString(borderRow(widths, "┌", "┬", "┐") + "\n")
	for i, r := range rows {
		var cells []string
		for j, w := range widths {
			c := ""
			if j < len(r.cells) {
				c = r.cells[j]
			}
			cells = append(cells, padRight(c, w))
		}
		b.WriteString("│ " + strings.Join(cells, " │ ") + " │\n")
		if r.sep && i < len(rows)-1 {
			b.WriteString(borderRow(widths, "├", "┼", "┤") + "\n")
		}
	}
	b.WriteString(borderRow(widths, "└", "┴", "┘"))
	return b.String()
}
