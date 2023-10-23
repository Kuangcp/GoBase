package main

import (
	"bytes"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"os/exec"
	"time"
)

var (
	iconPath = ".drink-water/msg.svg"
	dbPath   = ".drink-water/db"
)

var (
	addr   string
	db     int
	pass   string
	total  int
	help   bool
	drink  int
	stat   bool
	msgSec int
)

var HelpInfo = ctool.HelpInfo{
	Description: "Drink water per day",
	Version:     "1.0.0",
	Flags: []ctool.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help info"},
		{Short: "-s", BoolVar: &stat, Comment: "stat"},
	},
	Options: []ctool.ParamVO{
		{Short: "-a", StringVar: &addr, String: "", Value: "host", Comment: "redis host"},
		{Short: "-db", IntVar: &db, Int: 0, Value: "db", Comment: "redis db"},
		{Short: "-w", StringVar: &pass, String: "", Value: "password", Comment: "redis password"},
		{Short: "-d", IntVar: &drink, Int: 0, Value: "ml", Comment: "how much ml drink once"},
		{Short: "-t", IntVar: &total, Int: 0, Value: "path", Comment: "target total ml to drink per day"},
		{Short: "-m", IntVar: &msgSec, Int: 300, Value: "s", Comment: "notify msg second"},
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
	Conn := redis.NewClient(&opt)
	today := time.Now().Format(ctool.YYYY_MM_DD)

	if stat {
		notifyProgress(today, Conn)
		return
	}

	if drink > 0 {
		Conn.ZAdd("drink-water:"+today, redis.Z{Member: time.Now().UnixMilli(), Score: float64(drink)})
		notifyProgress(today, Conn)
		return
	}

	for t := range time.NewTicker(time.Second * time.Duration(msgSec)).C {
		now := t.Format(ctool.YYYY_MM_DD)
		notifyProgress(now, Conn)
	}
}

func notifyProgress(today string, Conn *redis.Client) {
	todayKey := "drink-water:" + today
	result, err := Conn.ZRangeWithScores(todayKey, 0, -1).Result()
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
