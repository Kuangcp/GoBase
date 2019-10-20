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
var wordDetail = 0 // 默认 0:输出路径, 1:所有字, 2:只有统计 3 统计加分析
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

var client = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6667",
	Password: "", // no password set
	DB:       0,  // use default DB
})

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
func handleFile(fileName string, handleChar func(string)) int {
	bytes := readFileToBytes(fileName)
	var total = 0
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
		// 满了三个字节, 就可以确定为一个汉字了
		if count == 3 {
			if handleChar != nil {
				handleChar(string(temp[0:3]))
			}
			total++
			count = 0
		}
	}
	return total
}

func showCharRank(start int64, stop int64) {
	result, e := client.ZRevRangeWithScores(charRankKey, start, stop).Result()
	if e != nil {
		println("error: ", e)
		return
	}
	for _, a := range result {
		printf("%-6v -> %v \n", a.Score, a.Member)
	}
}

func printParam(verb string, param string, comment string) {
	var format = "   %v %-5v %v %-10v %v %v\n"
	printf(format, green, verb, yellow, param, end, comment)
}

func help() {
	printf("  count %v <verb> %v <param> \n", green, yellow)
	printParam("", "", "遍历并统计所有文件汉字数")
	printParam("-h", "", "帮助")
	printParam("-w", "", "输出所有汉字")
	printParam("-s", "", "简洁输出总字数")
	printParam("-a", "", "统计单个文件字数, 存入 redis")
	printParam("-all", "", "统计字数 列出排行, 存入 redis")
	printParam("-del", "", "删除 redis 排行数据")
	printParam("-show", "start stop", "近列出排行, 读取 redis中的 zset 结构")
}

// 参数构成: 0 文件 1 参数 2 参数
func main() {
	param := os.Args
	if len(param) > 1 {
		verb := param[1]
		switch verb {
		case "-h":
			help()
			return
		case "-w":
			wordDetail = 1
			break
		case "-s":
			wordDetail = 2
			break
		case "-a":
			if len(param) < 4 {
				println("please input all param: filename, keyName")
				os.Exit(1)
			}

			fileName := param[2]
			keyName := param[3]
			log.Println(green, fileName, keyName, end)
			charRankKey = keyName
			handleFile(fileName, increaseCharNum)
			return
		case "-all":
			wordDetail = 3
			break
		case "-show":
			if len(param) < 4 {
				println("please input all param: start, stop")
				os.Exit(1)
			}
			start, err1 := strconv.ParseInt(param[2], 10, 64)
			stop, err2 := strconv.ParseInt(param[3], 10, 64)
			if err1 != nil || err2 != nil {
				println("please input correct param: start, stop")
				os.Exit(1)
			}
			showCharRank(start, stop)
			os.Exit(0)
		case "-del":
			client.Del(charRankKey)
			return
		}
	}

	// 递归遍历目录 读取所有文件
	filepath.Walk("./", handlerDir)
	for e := fileList.Front(); e != nil; e = e.Next() {
		fileName := e.Value.(string)
		// log.Println(yellow, fileName, end)
		if isNeedHandle(fileName) {
			if wordDetail == 3 {
				handleFile(fileName, increaseCharNum)
			} else {
				var total = 0
				if wordDetail == 1 {
					total = handleFile(fileName, showChar)
				} else {
					total = handleFile(fileName, nil)
				}
				// printf("%-50v-> chinese char =\033[0;32m %v\033[0m \n", fileName, total)
				if wordDetail == 0 {
					printf("%-5v %v %-5v %v %v \n", total, green, total, end, fileName)
				}
				totalFile += total
			}
		}
	}

	println()
	log.Printf("%v Total characters of all files : %v %v %v \n", purple, yellow, totalFile, end)

	if wordDetail != 3 {
		return
	}
	showCharRank(0, 10)
}
