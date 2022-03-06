package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-redis/redis/v7"
	"github.com/kuangcp/gobase/pkg/cuibase"
)

var charRankKey = "count:total_char_rank"
var redisConfigFile = os.Getenv("HOME") + "/.config/app-conf/count-char/redis.json"

var ignoreDirs = []string{".git", ".svn", ".vscode", ".idea", ".gradle",
	"out", "build", "target", "log", "logs", "__pycache__"}
var handleFiles = []string{".md", ".markdown", ".txt"}

var (
	handleFileSuffix string
	ignoreDir        string
	rankPair         string
	targetFile       string

	printAllChar bool
	countSummary bool
	help         bool
	delCache     bool
	allFile      bool
)
var (
	startIdx int64 = 0
	endIdx   int64 = 0
)

var info = cuibase.HelpInfo{
	Description: "Count chinese char(UTF8) from file that current dir recursive",
	Version:     "1.1.0",
	ValueLen:    12,
	Flags: []cuibase.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help"},
		{Short: "-s", BoolVar: &countSummary, Comment: "count all files, print summary"},
		{Short: "-p", BoolVar: &printAllChar, Comment: "print all chinese char"},
		{Short: "-d", BoolVar: &delCache, Comment: "delete rank data"},
		{Short: "-a", BoolVar: &allFile, Comment: "count chinese char on current dir"},
	},
	Options: []cuibase.ParamVO{
		{Short: "-c", Value: "start,end", Comment: "cursor"},
		{Short: "-f", Value: "file", Comment: "file"},

		{Short: "-S", Value: "suffix", Comment: ""},
		{Short: "-D", Value: "dir", Comment: ""},
	},
}

func readRedisConfig() *redis.Options {
	config := redis.Options{}
	data, e := ioutil.ReadFile(redisConfigFile)
	if e != nil {
		log.Fatal("read config file failed")
		return nil
	}
	e = json.Unmarshal(data, &config)
	if e != nil {
		log.Fatal("unmarshal config file failed")
		return nil
	}
	log.Println("redis ", config.Addr, config.DB)
	return &config
}
