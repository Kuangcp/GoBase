package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/cuibase"
)

var totalFile = 0
var totalChineseChar = 0
var fileList = list.New()

var client *redis.Client
var charRankKey = "total_char_rank"

var ignoreDirs = [...]string{
	".git", ".svn", ".vscode", ".idea", ".gradle",
	"out", "build", "target", "log", "logs", "__pycache__",
}
var redisConfigFile = "redis.json"

var handleFiles = [...]string{
	".md", ".markdown", ".txt", ".java", ".groovy", ".go", ".c", ".cpp", ".py",
}

func initRedisClient() {
	config := readRedisConfig()
	if config == nil{
		log.Fatal("config error")
	}
	client = redis.NewClient(config)
}

// 递归遍历目录
func handlerDir(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("occur error: ", err)
		return nil
	}

	if info.IsDir() {
		for _, dir := range ignoreDirs {
			if path == dir {
				return filepath.SkipDir
			}
		}
		return nil
	}
	fileList.PushBack(path)
	return nil
}

func isNeedHandle(filename string) bool {
	for _, fileType := range handleFiles {
		if strings.HasSuffix(filename, fileType) {
			return true
		}
	}
	return false
}

func readToBytes(fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("error: ", err)
		return nil
	}
	// 延迟关闭文件
	defer file.Close()
	//读取文件全部内容
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("occur error: ", err)
		return nil
	}
	return content
}

// increase count in redis
func increaseCharNum(CNChar string) {
	result, e := client.ZIncrBy(charRankKey, 1, CNChar).Result()
	if e == redis.Nil {
		println(result)
	}
}

func showChar(CNChar string) {
	fmt.Printf(CNChar)
}

// 根据字节流  统计中文字符
func handleFile(fileName string, handleChineseChar func(string)) int {
	bytes := readToBytes(fileName)
	var totalChar = 0
	var count = 0
	temp := [3]byte{}
	for _, item := range bytes {
		// 首字节 满足 110xxxxx [4e00,9fa5]  字节的头部分别为 0100 1001
		if count == 0 && item >= 228 && item <= 233 {
			temp[count] = item
			count++
			continue
		}
		// 后续字节 满足 10xxxxxx, 不满足就说明不是汉字了
		if count != 0 && item >= 128 && item <= 191 {
			temp[count] = item
			count++
		} else {
			count = 0
			continue
		}
		// 缓存满了三个字节, 就可以确定为一个汉字了
		if count == 3 {
			if handleChineseChar != nil {
				handleChineseChar(string(temp[0:3]))
			}
			totalChar++
			count = 0
		}
	}
	return totalChar
}

func showCharRank(start int64, stop int64) {
	result, e := client.ZRevRangeWithScores(charRankKey, start, stop).Result()
	if e != nil {
		log.Println("error: ", e)
		return
	}
	for _, a := range result {
		fmt.Printf("%-6v %v->%v %v \n", a.Score, cuibase.Green, cuibase.End, a.Member)
	}
}

func help(params []string) {
	cuibase.PrintTitleDefault("Count chinese char(UTF8) from file that current dir recursive")
	format := cuibase.BuildFormat(-5, -15)
	cuibase.PrintParams(format, []cuibase.ParamInfo{
		{
			Verb:    "-v",
			Param:   "",
			Comment: "show version",
		}, {
			Verb:    "-h",
			Param:   "",
			Comment: "help",
		}, {
			Verb:    "",
			Param:   "",
			Comment: "count all file chinese char",
		}, {
			Verb:    "-w",
			Param:   "",
			Comment: "print all chinese char",
		}, {
			Verb:    "-s",
			Param:   "",
			Comment: "count all file chinese char, show with simplify",
		}, {
			Verb:    "-a",
			Param:   "file redisKey",
			Comment: "count chinese char for target file, save. (redis)",
		}, {
			Verb:    "-all",
			Param:   "showNum",
			Comment: "count, calculate rank data, save. (redis)",
		}, {
			Verb:    "-del",
			Param:   "",
			Comment: "del rank data. (redis)",
		}, {
			Verb:    "-show",
			Param:   "start stop",
			Comment: "show rank data. (redis)",
		},
	})
}

func countWithRedis() {
	// 递归遍历目录 读取所有文件
	err := filepath.Walk("./", handlerDir)
	if err != nil {
		log.Println(err)
	}

	for e := fileList.Front(); e != nil; e = e.Next() {
		fileName := e.Value.(string)
		if isNeedHandle(fileName) {
			handleFile(fileName, increaseCharNum)
		}
	}
}

func showChineseChar(handler func(string), showFileInfo bool) {
	if showFileInfo {
		fmt.Printf("%v%-3v %-5v %-5v %v%v\n", cuibase.Yellow, "No", "Total", "Cur", "File", cuibase.End)
	}
	err := filepath.Walk("./", handlerDir)
	if err != nil {
		log.Println(err)
	}
	for e := fileList.Front(); e != nil; e = e.Next() {
		fileName := e.Value.(string)
		if !isNeedHandle(fileName) {
			continue
		}

		totalFile++
		var total = 0
		total = handleFile(fileName, handler)
		totalChineseChar += total
		if showFileInfo {
			fmt.Printf("%-3v %-5v %v%-5v %v%v \n",
				totalFile, totalChineseChar, cuibase.Green, total, cuibase.End, fileName)
		}
	}
	fmt.Printf("\nTotal characters. files: %v%v%v chars: %v%v%v\n",
		cuibase.Yellow, totalFile, cuibase.End, cuibase.Yellow, totalChineseChar, cuibase.End)
}

func readTargetFile(params []string) {
	cuibase.AssertParamCount(3, "Please input param: filename keyName")

	fileName := params[2]
	charRankKey = params[3]

	log.Printf("%vread file: %v, redis key: %v%v\n",
		cuibase.Green, fileName, charRankKey, cuibase.End)
	initRedisClient()
	handleFile(fileName, increaseCharNum)

	if len(params) == 5 {
		num, err := strconv.ParseInt(params[4], 10, 64)
		if err == nil {
			showCharRank(0, num-1)
			return
		}
	}
	showCharRank(0, 15)
}

func readAllSaveIntoRedis(params []string) {
	initRedisClient()
	countWithRedis()
	if len(params) == 3 {
		num, err := strconv.ParseInt(params[2], 10, 64)
		if err == nil {
			showCharRank(0, num-1)
			return
		}
	}

	showCharRank(0, 10)
}

func showRank(params []string) {
	cuibase.AssertParamCount(3, "Please input all param: start stop")
	start, err1 := strconv.ParseInt(params[2], 10, 64)
	stop, err2 := strconv.ParseInt(params[3], 10, 64)
	if err1 != nil || err2 != nil {
		log.Fatal("please input correct param: start, stop")
	}
	initRedisClient()
	showCharRank(start, stop)
}

func delRank(params []string) {
	initRedisClient()
	client.Del(charRankKey)
	log.Printf("del %v%v%v", cuibase.Green, charRankKey, cuibase.End)
}

func readRedisConfig() *redis.Options{
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

func main() {
	param := os.Args
	if len(param) <= 1 {
		showChineseChar(nil, true)
		return
	}

	cuibase.RunAction(map[string]func(params []string){
		"-h": help,
		"-w": func(params []string) {
			showChineseChar(showChar, false)
		},
		"-s": func(params []string) {
			showChineseChar(nil, false)
		},
		"-f":    readTargetFile,
		"-all":  readAllSaveIntoRedis,
		"-show": showRank,
		"-del":  delRank,
		"-v": func(params []string) {
			println("v1.0.0")
		},
	}, help)
}
