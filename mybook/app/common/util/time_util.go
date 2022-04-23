package util

import (
	"strconv"
	"strings"
	"time"
)

//判断时间是当年的第几周
func WeekByDate(t time.Time) int {
	year, week := WeekOfYearByDate(t)
	return year*100 + week
}

//判断时间是当年的第几周
func WeekOfYearByDate(t time.Time) (int, int) {
	yearDay := t.YearDay()
	yearFirstDay := t.AddDate(0, 0, -yearDay+1)
	firstDayInWeek := int(yearFirstDay.Weekday())

	//今年第一周有几天
	firstWeekDays := 1
	if firstDayInWeek != 0 {
		firstWeekDays = 7 - firstDayInWeek + 1
	}
	var week int
	if yearDay <= firstWeekDays {
		week = 1
	} else {
		week = (yearDay-firstWeekDays)/7 + 2
	}
	return t.Year(), week
}

func MonthByDate(t time.Time) int {
	return t.Year()*100 + int(t.Month())
}

func StartDayByYearWeek(yearWeek string) string {
	split := strings.Split(yearWeek, "-")
	return StartDayByWeekAndYear(split[0], split[1])
}

func StartDayByWeekAndYear(year, week string) string {
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return ""
	}
	weekInt, err := strconv.Atoi(week)
	if err != nil {
		return ""
	}

	firstDay := time.Date(yearInt, 1, 1, 0, 0, 0, 0, time.Local)
	date := firstDay.AddDate(0, 0, weekInt*7)
	return year + "_" + week + " " + date.Format("01-02")
}
