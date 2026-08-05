package main

import (
	"strings"
	"testing"
	"time"
)

func sampleResult() *AnalysisResult {
	now := time.Now()
	var dates []string
	for i := 6; i >= 0; i-- {
		dates = append(dates, now.AddDate(0, 0, -i).Format("2006-01-02"))
	}

	added := []AuthorDayData{
		{Name: "alice", Data: []int{10, 0, 5, 0, 2, 0, 3}},
		{Name: "bob", Data: []int{0, 0, 0, 0, 0, 1, 0}},
	}
	deleted := []AuthorDayData{
		{Name: "alice", Data: []int{1, 0, 0, 0, 1, 0, 0}},
		{Name: "bob", Data: []int{0, 0, 0, 0, 0, 0, 0}},
	}
	commits := []AuthorDayData{
		{Name: "alice", Data: []int{1, 0, 1, 0, 1, 0, 1}},
		{Name: "bob", Data: []int{0, 0, 0, 0, 0, 1, 0}},
	}

	return &AnalysisResult{
		RepoPath:           "/tmp/demo",
		RepoName:           "demo",
		Branch:             "main",
		TotalCommits:       7,
		Authors:            []AuthorStat{{Name: "alice", CommitCount: 5}, {Name: "bob", CommitCount: 2}},
		DateRange:          dates,
		AuthorSeries:       commits,
		AuthorAddedSeries:  added,
		AuthorDeletedSeries: deleted,
	}
}

func TestRenderProjectOverview(t *testing.T) {
	out := RenderProjectOverview(sampleResult())

	for _, want := range []string{"项目概览", "总提交数", "近30天提交", "参与人", "7", "alice, bob", "仓库"} {
		if !strings.Contains(out, want) {
			t.Errorf("overview missing %q in:\n%s", want, out)
		}
	}
}

func TestRenderAuthorWeekTable(t *testing.T) {
	result := sampleResult()
	// day0: alice added10/deleted1 (net +9), bob added12/deleted10 (net +2)
	result.AuthorAddedSeries[1].Data[0] = 12
	result.AuthorDeletedSeries[1].Data[0] = 10
	result.AuthorSeries[1].Data[0] = 2
	out := RenderAuthorWeekTable(result)

	for _, want := range []string{"最近7天代码变更", "日期", "人员", "提交数", "新增", "删除", "净变更", "当日合计", "alice", "bob", "无提交"} {
		if !strings.Contains(out, want) {
			t.Errorf("table missing %q in:\n%s", want, out)
		}
	}

	for _, want := range []string{"+9", "+5", "+3", "+1"} {
		if !strings.Contains(out, want) {
			t.Errorf("table missing net %q", want)
		}
	}

	// on day0 alice (+9) must appear above bob (+2): net-sorted descending
	ai := strings.Index(out, "alice")
	bi := strings.Index(out, "bob")
	if ai < 0 || bi < 0 || ai > bi {
		t.Errorf("expected alice (net +9) before bob (net +2) on day0")
	}
}

func TestRenderOverviewParticipantsWrap(t *testing.T) {
	result := sampleResult()
	result.Authors = []AuthorStat{
		{Name: "alice", CommitCount: 1}, {Name: "bob", CommitCount: 1}, {Name: "carol", CommitCount: 1},
		{Name: "dave", CommitCount: 1}, {Name: "eve", CommitCount: 1}, {Name: "frank", CommitCount: 1},
		{Name: "grace", CommitCount: 1}, {Name: "heidi", CommitCount: 1}, {Name: "ivan", CommitCount: 1},
		{Name: "judy", CommitCount: 1}, {Name: "mallory", CommitCount: 1}, {Name: "oscar", CommitCount: 1},
		{Name: "peggy", CommitCount: 1}, {Name: "trent", CommitCount: 1},
	}
	out := RenderProjectOverview(result)

	for _, want := range []string{"alice", "judy", "trent"} {
		if !strings.Contains(out, want) {
			t.Errorf("overview missing %q in:\n%s", want, out)
		}
	}
	if got := strings.Count(out, "参与人"); got != 1 {
		t.Errorf("expected one 参与人 label, got %d", got)
	}
	// every box line must fit within 80 columns even with many participants
	for _, line := range strings.Split(out, "\n") {
		if w := displayWidth(line); w > 80 {
			t.Errorf("box line too wide (%d): %q", w, line)
		}
	}
}

func TestWrapByWidth(t *testing.T) {
	cases := []struct {
		in   string
		maxW int
		want string
	}{
		{"alice, bob, carol, dave, eve", 10, "alice, bob|carol|dave, eve"},
		{"alice, bob, carol", 26, "alice, bob, carol"},
		{"", 10, ""},
	}
	for _, c := range cases {
		if got := strings.Join(wrapByWidth(c.in, c.maxW), "|"); got != c.want {
			t.Errorf("wrapByWidth(%q, %d) = %q, want %q", c.in, c.maxW, got, c.want)
		}
	}
}

func TestRenderCommitCurve(t *testing.T) {
	out := RenderCommitCurve(sampleResult())
	if !strings.Contains(out, "最近30天提交曲线") {
		t.Errorf("curve missing caption:\n%s", out)
	}
}

func TestDisplayWidth(t *testing.T) {
	cases := map[string]int{
		"abc":      3,
		"净变更":      6,
		"08-01 周六": 10,
		"":         0,
	}
	for s, want := range cases {
		if got := displayWidth(s); got != want {
			t.Errorf("displayWidth(%q) = %d, want %d", s, got, want)
		}
	}
}

func TestFormatThousands(t *testing.T) {
	cases := map[int]string{
		0:     "0",
		999:   "999",
		1000:  "1,000",
		1234567: "1,234,567",
		-1234: "-1,234",
	}
	for n, want := range cases {
		if got := formatThousands(n); got != want {
			t.Errorf("formatThousands(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestNetStr(t *testing.T) {
	cases := map[int]string{0: "0", 75: "+75", -15: "-15"}
	for n, want := range cases {
		if got := netStr(n); got != want {
			t.Errorf("netStr(%d) = %q, want %q", n, got, want)
		}
	}
}
