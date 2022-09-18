package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Lofanmi/chinese-calendar-golang/calendar"
	"github.com/kuangcp/gobase/pkg/ctool"
	"strings"
	"time"
)

type (
	Lunar struct {
		Day        int    `json:"day"`
		DayAlias   string `json:"day_alias"`
		Month      int    `json:"month"`
		MonthAlias string `json:"month_alias"`
		Year       int    `json:"year"`
		YearAlias  string `json:"year_alias"`
	}
	Solar struct {
		Day       int    `json:"day"`
		Month     int    `json:"month"`
		WeekAlias string `json:"week_alias"`
		WeekDay   int    `json:"week_number"`
		Year      int    `json:"year"`
	}
	LunarCal struct {
		Lunar Lunar `json:"lunar"`
		Solar Solar `json:"solar"`
	}
)

var now = time.Now()
var (
	startMonth int
	month      int
)

func main() {
	flag.IntVar(&startMonth, "s", 0, "month cursor")
	flag.IntVar(&month, "m", 1, "month count")

	flag.Parse()

	//var lunar = toLunar(time.Now())
	//fmt.Println(lunar)

	firstDay := now.AddDate(0, startMonth, -now.Day()+1)
	for i := 0; i < month; i++ {
		date := firstDay.AddDate(0, startMonth+i, 0)
		//fmt.Println(date)
		fmt.Println(buildMonthBlock(date))
	}
}

func buildMonthBlock(first time.Time) string {
	end := first.AddDate(0, 1, 0)
	firstLine := true
	result := ""
	footLine := ""
	for first.Before(end) {
		tmpLunar := toLunar(first)
		weekDay := tmpLunar.Solar.WeekDay
		if firstLine {
			firstLine = false
			result += buildTitle(tmpLunar)
			result += fmt.Sprint(strings.Repeat("      ", (weekDay+6)%7))
			footLine += fmt.Sprint(strings.Repeat("      ", (weekDay+6)%7))
		}
		if weekDay == 1 {
			result += "\n"
			result += footLine + "\n"
			footLine = ""
		}
		if first.Equal(now) {
			result += ctool.Yellow.Print(fmt.Sprintf("%5v ", first.Day()))
			footLine += "  " + ctool.Yellow.Print(getDay(tmpLunar.Lunar))
		} else {
			result += fmt.Sprintf("%5v ", first.Day())
			footLine += "  " + ctool.LightBlue.Print(getDay(tmpLunar.Lunar))
		}

		first = first.AddDate(0, 0, 1)
	}
	return result + "\n" + footLine + "\n"
}

func getDay(lunar Lunar) string {
	if lunar.DayAlias == "初一" {
		return lunar.MonthAlias
	} else {
		return lunar.DayAlias
	}
}

func toLunar(t time.Time) LunarCal {
	cal := calendar.ByTimestamp(t.Unix())
	jsonBt, _ := cal.ToJSON()

	var lunar LunarCal
	//fmt.Println(string(jsonBt))
	_ = json.Unmarshal(jsonBt, &lunar)
	return lunar
}

func buildTitle(now LunarCal) string {
	var s = []string{"一", "二", "三", "四", "五", "六", "日"}
	return fmt.Sprintf("%18v  %v\n    %v\n", now.Solar.Year, now.Solar.Month,
		ctool.Green.Print(strings.Join(s, "    ")))
}
