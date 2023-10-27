package main

import (
	"bytes"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

const (
	redisPrefix = "drink-water:"
)

var (
	iconPath = ".drink-water/msg.svg"
	dbPath   = ".drink-water/db"
	cli      *redis.Client
	xs       []string
)

var (
	addr    string
	db      int
	pass    string
	total   int
	help    bool
	drink   int
	stat    bool
	msgSec  int
	webPort int
)

type (
	dayLine struct {
		day  string
		data []opts.LineData
	}
)

var HelpInfo = ctool.HelpInfo{
	Description: "Drinking water per day",
	Version:     "1.0.1",
	Flags: []ctool.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help info"},
		{Short: "-s", BoolVar: &stat, Comment: "stat"},
	},
	Options: []ctool.ParamVO{
		{Short: "-a", StringVar: &addr, String: "", Value: "host", Comment: "redis address host:port"},
		{Short: "-b", IntVar: &db, Int: 0, Value: "db", Comment: "redis db"},
		{Short: "-w", StringVar: &pass, String: "", Value: "password", Comment: "redis password"},
		{Short: "-d", IntVar: &drink, Int: 0, Value: "ml", Comment: "how much ml drink once"},
		{Short: "-t", IntVar: &total, Int: 0, Value: "path", Comment: "target total ml to drink per day"},
		{Short: "-m", IntVar: &msgSec, Int: 300, Value: "seconds", Comment: "notify msg frequency(seconds)"},
		{Short: "-p", IntVar: &webPort, Int: 33380, Value: "port", Comment: "web server port"},
	},
}

func main() {
	HelpInfo.Parse()
	if help {
		HelpInfo.PrintHelp()
		return
	}

	home, _ := ctool.Home()
	iconPath = home + "/" + iconPath
	dbPath = home + "/" + dbPath

	var opt redis.Options
	opt = redis.Options{Addr: addr, Password: pass, DB: db, PoolSize: 1}
	cli = redis.NewClient(&opt)
	if !IsValidConnection(cli) {
		return
	}

	today := time.Now().Format(ctool.YYYY_MM_DD)
	if stat {
		notifyProgress(today)
		return
	}

	if drink > 0 {
		cli.ZAdd(redisPrefix+today, redis.Z{Member: time.Now().UnixMilli(), Score: float64(drink)})
		notifyProgress(today)
		return
	}

	go func() {
		logger.Info("start web server on localhost:%v", webPort)
		mux := http.NewServeMux()
		mux.HandleFunc("/history", history)
		http.ListenAndServe(fmt.Sprintf(":%v", webPort), mux)
	}()
	for t := range time.NewTicker(time.Second * time.Duration(msgSec)).C {
		now := t.Format(ctool.YYYY_MM_DD)
		notifyProgress(now)
	}
}

func history(writer http.ResponseWriter, request *http.Request) {
	var param struct {
		Frame int
		Start *time.Time `form:"start" fmt:"2006-01-02"`
		End   *time.Time `form:"end" fmt:"2006-01-02"`
	}
	if err := ctool.Unpack(request, &param); err != nil {
		return
	}
	if param.End == nil {
		now := time.Now()
		param.End = &now
	}
	if param.Frame == 0 {
		param.Frame = 1
	}
	xs = []string{}
	hour := param.Frame == 1
	for i := 0; i < 24/param.Frame; i++ {
		if hour {
			xs = append(xs, fmt.Sprint(i+1))
		} else {
			xs = append(xs, fmt.Sprintf("%v - %v", param.Frame*i+1, param.Frame*(i+1)))
		}
	}

	if param.Start == nil {
		today := time.Now().Format(ctool.YYYY_MM_DD)
		RenderChart(writer, dayLine{day: today, data: buildDayData(today, param.Frame)})
	} else {
		t := time.Now()
		todayZero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		var days []dayLine
		for !param.Start.After(todayZero) {
			day := param.Start.Format(ctool.YYYY_MM_DD)
			days = append(days, dayLine{day: day, data: buildDayData(day, param.Frame)})
			date := param.Start.AddDate(0, 0, 1)
			param.Start = &date
		}
		RenderChart(writer, days...)
	}
}

func buildDayData(day string, frame int) []opts.LineData {
	cache := make([]int, 24)
	var items []opts.LineData

	todayKey := redisPrefix + day
	result, err := cli.ZRangeWithScores(todayKey, 0, -1).Result()
	if err != nil {
		logger.Error(err)
		return items
	}
	for _, m := range result {
		ms, err := strconv.ParseInt(fmt.Sprint(m.Member), 10, 64)
		if err != nil {
			logger.Error(err)
			continue
		}
		msTime := time.UnixMilli(ms)
		cache[(msTime.Hour()-1)/frame] += int(m.Score)
	}
	for i := 0; i < len(cache); i++ {
		items = append(items, opts.LineData{Value: cache[i]})
	}

	return items
}

func RenderChart(writer http.ResponseWriter, days ...dayLine) {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "每日饮水",
			Link:  "https://github.com/go-echarts/go-echarts",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis"}),
	)

	line.SetXAxis(xs)
	for _, d := range days {
		line.AddSeries(d.day, d.data).SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: false, ShowSymbol: true, SymbolSize: 15, Symbol: "diamond",
			}),
			charts.WithLabelOpts(opts.Label{
				Show: true,
			}))
	}

	page := components.NewPage()
	page.AddCharts(line)
	page.PageTitle = "每日饮水"
	page.Render(writer)
}

func IsValidConnection(client *redis.Client) bool {
	_, err := client.Ping().Result()
	if err != nil {
		logger.Error("ping redis failed:", client.Options(), err)
		return false
	}
	return true
}

func notifyProgress(today string) {
	todayKey := redisPrefix + today
	result, err := cli.ZRangeWithScores(todayKey, 0, -1).Result()
	if err != nil {
		logger.Error(err)
		return
	}
	var sum int
	for _, m := range result {
		sum += int(m.Score)
	}

	notify(total, sum)
}

func notify(total, progress int) {
	title := fmt.Sprintf("喝水 %dL", total/1000)
	content := fmt.Sprintf("饮用：%dml\n进度：%v%%", progress, progress*100/total)
	execCommand(fmt.Sprintf("notify-send -i %s '%s' '%s' -t %v", iconPath, title, content, 30_000))
}

func execCommand(command string) (string, bool) {
	cmd := exec.Command("/usr/bin/bash", "-c", command)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.Error(err)
		return "", false
	}

	result := out.String()
	return result, true
}
