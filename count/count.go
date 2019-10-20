package main

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

// WARN: 只对 采用 UTF-8 编码的文件保证统计正确

var totalFile = 0
var totalCNChar = 0
var fileList = list.New()

var printf = fmt.Printf
var green = "\033[0;32m"
var yellow = "\033[0;33m"
var purple = "\033[0;35m"
var end = "\033[0m"

var charRankKey = "total_char_rank"

var ignoreDirList = [...]string{
	".git", ".svn", ".vscode", ".idea", ".gradle",
	"out", "build", "target", "log", "logs", "__pycache__",
}
var handleFileList = [...]string{
	".md", ".markdown", ".txt", ".java", ".groovy", ".go", ".c", ".cpp", ".py",
}

var client *redis.Client

func initRedisClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6666",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// 递归遍历目录
func handlerDir(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println("occur error: ", err)
		return nil
	}

	if info.IsDir() {
		for _, dir := range ignoreDirList {
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
	for _, fileType := range handleFileList {
		if strings.HasSuffix(filename, fileType) {
			return true
		}
	}
	return false
}

func readFileToBytes(fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("error: ", err)
		return nil
	}
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
	printf(CNChar)
}

// 根据字节流  统计中文字符
func handleFile(fileName string, handleCNChar func(string)) int {
	bytes := readFileToBytes(fileName)
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
			if handleCNChar != nil {
				handleCNChar(string(temp[0:3]))
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
		println("error: ", e)
		return
	}
	for _, a := range result {
		printf("%-6v %v->%v %v \n", a.Score, green, end, a.Member)
	}
}

func printParam(verb string, param string, comment string) {
	var format = "   %v %-5v %v %-15v %v %v\n"
	printf(format, green, verb, yellow, param, end, comment)
}

func help() {
	printf("  count %v <verb> %v <param> \n", green, yellow)
	printParam("-h", "", "帮助")
	printParam("", "", "遍历并统计所有文件汉字数")
	printParam("-w", "", "输出所有汉字")
	printParam("-s", "", "简洁输出总字数")
	printParam("-a", "file redisKey", "统计单个文件字数, 存入 redis")
	printParam("-all", "showNum", "统计字数 列出排行, 存入 redis")
	printParam("-del", "", "删除 redis 排行数据")
	printParam("-show", "start stop", "近列出排行, 读取 redis中的 zset 结构")
}

func countWithRedis() {
	// 递归遍历目录 读取所有文件
	filepath.Walk("./", handlerDir)
	for e := fileList.Front(); e != nil; e = e.Next() {
		fileName := e.Value.(string)
		if isNeedHandle(fileName) {
			handleFile(fileName, increaseCharNum)
		}
	}
}

func showAllCNChar(handleCNChar func(string), showFileInfo bool) {
	if showFileInfo {
		printf("%v%-3v %-5v %-5v %v%v\n", yellow, "No", "Total", "Cur", "File", end)
	}
	filepath.Walk("./", handlerDir)
	for e := fileList.Front(); e != nil; e = e.Next() {
		fileName := e.Value.(string)
		if isNeedHandle(fileName) {
			totalFile++
			var total = 0
			total = handleFile(fileName, handleCNChar)
			totalCNChar += total
			if showFileInfo {
				printf("%-3v %-5v %v%-5v %v%v \n", totalFile, totalCNChar, green, total, end, fileName)
			}
		}
	}
	printf("\nTotal characters. files: %v%v%v chars: %v%v%v\n", yellow, totalFile, end, yellow, totalCNChar, end)
}

// 参数构成: 0 文件 1 参数 2 参数
func main() {
	param := os.Args
	if len(param) <= 1 {
		showAllCNChar(nil, true)
		return
	}

	switch param[1] {
	case "-h":
		help()
	case "-w":
		showAllCNChar(showChar, false)
	case "-s":
		showAllCNChar(nil, false)
	case "-f":
		if len(param) < 4 {
			log.Fatal("please input param: filename keyName")
		}

		fileName := param[2]
		charRankKey = param[3]

		log.Printf("%v read file: %v, redis key: %v %v ", green, fileName, charRankKey, end)
		initRedisClient()
		handleFile(fileName, increaseCharNum)

		if len(param) == 5 {
			num, err := strconv.ParseInt(param[4], 10, 64)
			if err == nil {
				showCharRank(0, num-1)
				return
			}
		}
		showCharRank(0, 15)
	case "-all":
		initRedisClient()
		countWithRedis()
		if len(param) == 3 {
			num, err := strconv.ParseInt(param[2], 10, 64)
			if err == nil {
				showCharRank(0, num-1)
				return
			}
		}

		showCharRank(0, 10)
	case "-show":
		if len(param) < 4 {
			log.Fatal("please input all param: start, stop")
		}
		start, err1 := strconv.ParseInt(param[2], 10, 64)
		stop, err2 := strconv.ParseInt(param[3], 10, 64)
		if err1 != nil || err2 != nil {
			log.Fatal("please input correct param: start, stop")
		}
		initRedisClient()
		showCharRank(start, stop)
	case "-del":
		initRedisClient()
		client.Del(charRankKey)
		log.Printf("del %v%v%v", green, charRankKey, end)
	default:
		help()
		return
	}
}
