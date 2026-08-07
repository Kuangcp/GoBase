package web

import (
	"testing"
	"time"
)

func micros(dayOffset, hour, min int) int64 {
	base := time.Date(2026, 8, 1, 0, 0, 0, 0, time.Local)
	t := base.AddDate(0, 0, dayOffset).Add(time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute)
	return t.UnixNano() / 1000
}

func assertWorkMinutes(t *testing.T, got, want []int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %d, want %d", i, got[i], want[i])
		}
	}
}

func TestComputeWorkTime_SingleDay(t *testing.T) {
	r := []dayTimeRange{
		{first: micros(0, 9, 0), last: micros(0, 17, 0)},
		{},
		{},
	}
	assertWorkMinutes(t, computeWorkTime(r, nil), []int{480, 0, 0})
}

func TestComputeWorkTime_CrossMidnight_WithAfternoon(t *testing.T) {
	r := []dayTimeRange{
		{first: micros(0, 22, 0), last: micros(0, 23, 0), late: true},
		{first: micros(1, 0, 0), last: micros(1, 18, 0), lastMorning: micros(1, 2, 0), firstNoon: micros(1, 13, 0)},
		{},
	}
	assertWorkMinutes(t, computeWorkTime(r, nil), []int{240, 300, 0})
}

func TestComputeWorkTime_CrossMidnight_NoAfternoon(t *testing.T) {
	r := []dayTimeRange{
		{first: micros(0, 22, 0), last: micros(0, 23, 0), late: true},
		{first: micros(1, 0, 0), last: micros(1, 2, 0), lastMorning: micros(1, 2, 0)},
		{},
	}
	assertWorkMinutes(t, computeWorkTime(r, nil), []int{240, 0, 0})
}

func TestComputeWorkTime_EarlyMorningNotAbsorbed(t *testing.T) {
	r := []dayTimeRange{
		{first: micros(0, 0, 30), last: micros(0, 1, 30), lastMorning: micros(0, 1, 30)},
		{},
	}
	assertWorkMinutes(t, computeWorkTime(r, nil), []int{60, 0})
}

func TestComputeWorkTime_DirtyDataGuard(t *testing.T) {
	r := []dayTimeRange{
		{first: micros(0, 0, 0), last: micros(0, 23, 0), late: true},
		{first: micros(1, 0, 0), last: micros(1, 18, 0), lastMorning: micros(1, 10, 0), firstNoon: micros(1, 13, 0)},
		{},
	}
	assertWorkMinutes(t, computeWorkTime(r, nil), []int{0, 300, 0})
}

func TestComputeWorkTime_ConsecutiveOvernights(t *testing.T) {
	r := []dayTimeRange{
		{first: micros(0, 13, 0), last: micros(0, 23, 0), late: true, firstNoon: micros(0, 13, 0)},
		{first: micros(1, 0, 0), last: micros(1, 23, 0), lastMorning: micros(1, 2, 0), firstNoon: micros(1, 13, 0), late: true},
		{first: micros(2, 0, 0), last: micros(2, 2, 0), lastMorning: micros(2, 2, 0)},
	}
	assertWorkMinutes(t, computeWorkTime(r, nil), []int{780, 780, 0})
}

// 前日深夜(22点)与次日首键(10点)间隔长(11h)，不是真加班，不应跨日扩展
func TestComputeWorkTime_LongGapNoCrossMidnight(t *testing.T) {
	r := []dayTimeRange{
		{first: micros(0, 9, 0), last: micros(0, 22, 30), late: true},
		{first: micros(1, 10, 0), last: micros(1, 21, 0), lastMorning: micros(1, 10, 30), firstNoon: micros(1, 10, 30)},
		{},
	}
	assertWorkMinutes(t, computeWorkTime(r, nil), []int{810, 660, 0})
}

// 即使相隔很近，只要提供的日期串不连续(中间缺了一天)，也不做跨夜计算
func TestComputeWorkTime_NonConsecutiveDaysNoCross(t *testing.T) {
	r := []dayTimeRange{
		{first: micros(0, 22, 30), last: micros(0, 23, 0), late: true, lastMorning: micros(0, 23, 0)},
		{first: micros(2, 0, 30), last: micros(2, 2, 0), lastMorning: micros(2, 2, 0)},
	}
	days := []string{"2026:08:01", "2026:08:03"}
	got := computeWorkTime(r, days)
	// 不连续: 第0日不扩展到第1日凌晨; 第1日不被第0日吸收
	assertWorkMinutes(t, got, []int{30, 90})
}
