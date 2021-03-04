package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/kuangcp/gobase/cuibase"
	"io/ioutil"
	"log"
	"os"
)

var charRankKey = "total_char_rank"
var redisConfigFile = os.Getenv("HOME") + "/.config/app-conf/count-char/redis.json"

var ignoreDirs = []string{".git", ".svn", ".vscode", ".idea", ".gradle",
	"out", "build", "target", "log", "logs", "__pycache__"}
var handleFiles = []string{".md", ".markdown", ".txt"}

var (
	handleFileSuffix string
	ignoreDir        string
)

var info = cuibase.HelpInfo{
	Description: "Count chinese char(UTF8) from file that current dir recursive",
	Version:     "1.0.0",
	VerbLen:     -5,
	ParamLen:    -15,
	Params: []cuibase.ParamInfo{
		{
			Verb:    "-h",
			Param:   "",
			Comment: "Help",
		}, {
			Verb:    "",
			Param:   "",
			Comment: "Count all file chinese char",
		}, {
			Verb:    "-w",
			Param:   "",
			Comment: "Print all chinese char",
			Handler: func(_ []string) {
				showChineseChar(showChar, false)
			},
		}, {
			Verb:    "-f",
			Param:   "",
			Comment: "Read target file",
			Handler: readTargetFile,
		}, {
			Verb:    "-s",
			Param:   "",
			Comment: "Count all file chinese char, show with simplify",
			Handler: func(_ []string) {
				showChineseChar(nil, false)
			},
		}, {
			Verb:    "-a",
			Param:   "file redisKey",
			Comment: "Count chinese char for target file, save. (redis)",
		}, {
			Verb:    "-all",
			Param:   "showNum",
			Comment: "Count, calculate rank data, save. (redis)",
			Handler: readAllSaveIntoRedis,
		}, {
			Verb:    "-del",
			Param:   "",
			Comment: "Del rank data. (redis)",
			Handler: delRank,
		}, {
			Verb:    "-show",
			Param:   "start stop",
			Comment: "Show rank data. (redis)",
			Handler: showRank,
		},
	}}

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
