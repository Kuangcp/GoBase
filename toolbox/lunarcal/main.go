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
	LunarMonth struct {
		leapMark   bool
		title      string
		weeks      []string
		lunarWeeks []string
	}
)

var now = time.Now()
var s = []string{"一", "二", "三", "四", "五", "六", "日"}
var weekDay = ctool.Green.Print(strings.Join(s, "   "))
var (
	startMonth  int
	monthNumber int
)

func (l *LunarMonth) toString() string {
	s := len(l.weeks)
	block := ""
	for i := 0; i < s; i++ {
		block += l.weeks[i] + "\n" + l.lunarWeeks[i] + "\n"
	}
	return fmt.Sprintf("%v\n  %v\n\n%v", l.title, weekDay, block)
}

func main() {
	flag.IntVar(&startMonth, "s", 0, "month cursor")
	flag.IntVar(&monthNumber, "n", 1, "month number")
	flag.Parse()

	firstDay := now.AddDate(0, startMonth, -now.Day()+1)
	// one month one block
	//for i := 0; i < month; i++ {
	//	date := firstDay.AddDate(0, startMonth+i, 0)
	//	fmt.Println(buildMonthBlock(date).toString())
	//}

	var list []*LunarMonth
	for i := 0; i < monthNumber; i++ {
		date := firstDay.AddDate(0, startMonth+i, 0)
		list = append(list, buildMonthBlock(date))
	}

	var blockRightSplit = " │"
	var block = []string{"", "", "", "", "", "", "", "", "", "", "", "", "", ""}
	for i := range list {
		lunarMonth := list[i]
		block[0] += lunarMonth.title + blockRightSplit
		block[1] += "   " + weekDay + blockRightSplit

		i2 := len(lunarMonth.weeks)
		for j := 0; j < i2; j++ {
			block[2*j+2] += buildWeekLineBlock(lunarMonth.weeks[j], blockRightSplit)
			block[2*j+2+1] += buildWeekLineBlock(lunarMonth.lunarWeeks[j], blockRightSplit)
		}

		if i%3 == 2 {
			printBlock(block)
			block = []string{"", "", "", "", "", "", "", "", "", "", "", "", "", ""}
		}
	}
	printBlock(block)
}

func buildWeekLineBlock(block, split string) string {
	if block == "" {
		return strings.Repeat(" ", 35) + split
	} else {
		return block + split
	}
}

func printBlock(block []string) {
	s := ""
	for _, v := range block {
		if v == "" {
			continue
		}
		s += v + "\n"
	}

	if s != "" {
		fmt.Print(s)
		count := strings.Count(block[0], "│")
		fmt.Println(strings.Repeat("─", 37*count))
	}
}

func buildMonthBlock(first time.Time) *LunarMonth {
	end := first.AddDate(0, 1, 0)
	firstLine := true
	result := ""
	footLine := ""

	var weeks []string
	var lunarWeeks []string

	month := &LunarMonth{
		weeks:      weeks,
		lunarWeeks: lunarWeeks,
	}

	for first.Before(end) {
		tmpLunar := toLunar(first)
		weekDay := tmpLunar.Solar.WeekDay
		if firstLine {
			firstLine = false
			month.title = buildTitle(first)
			result += fmt.Sprint(strings.Repeat("     ", (weekDay+6)%7))
			footLine += fmt.Sprint(strings.Repeat("     ", (weekDay+6)%7))
		}
		if weekDay == 1 {
			month.weeks = append(month.weeks, result)
			month.lunarWeeks = append(month.lunarWeeks, footLine)
			footLine = ""
			result = ""
		}

		if first.Equal(now) {
			result += ctool.Yellow.Print(fmt.Sprintf("%4v ", first.Day()))
			footLine += ctool.Yellow.Print(getDay(tmpLunar.Lunar, month))
		} else {
			result += fmt.Sprintf("%4v ", first.Day())
			footLine += ctool.LightBlue.Print(getDay(tmpLunar.Lunar, month))
		}

		first = first.AddDate(0, 0, 1)
	}

	actualDay := first.AddDate(0, 0, -1)
	weekDay := int(actualDay.Weekday())
	result += fmt.Sprint(strings.Repeat("     ", (7-weekDay)%7))
	footLine += fmt.Sprint(strings.Repeat("     ", (7-weekDay)%7))

	if result != "" {
		month.weeks = append(month.weeks, result)
		month.lunarWeeks = append(month.lunarWeeks, footLine)
	}
	return month
}

func getDay(lunar Lunar, m *LunarMonth) string {
	if lunar.DayAlias == "初一" {
		if strings.Contains(lunar.MonthAlias, "闰") {
			m.leapMark = true
			return lunar.MonthAlias
		}
		return " " + lunar.MonthAlias
	} else if m.leapMark {
		m.leapMark = false
		return lunar.DayAlias
	} else {
		return " " + lunar.DayAlias
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

func buildTitle(first time.Time) string {
	return fmt.Sprintf("%23v%-12v", first.Format(ctool.YYYY_MM), "")
}
