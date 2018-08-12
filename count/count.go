package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"
	"strconv"
	"github.com/go-redis/redis"
)

// WARN: 该程序只对 采用 UTF-8 编码的文件保证统计正确

var totalFile = 0
var wordDetail = 0 // 默认 0:输出路径, 1:所有字, 2:只有统计 3 统计加分析
var printf = fmt.Printf
var analysisKey = "total_char_rank"

var green = "\033[0;32m"
var yellow = "\033[0;33m"
var purple = "\033[0;35m"
var end = "\033[0m"

var client = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6666",
	Password: "", // no password set
	DB:       0,  // use default DB
})

// 往递归遍历目录 作为参数传入的函数 
func handlerDir(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println("occur error: ", err)
		return nil
	}
	var ignoreDirList = [...]string{
		".git", ".svn", ".vscode", ".idea", ".gradle", "out", "build", "target", "log", "logs", "__pycache__",
	}
	if info.IsDir() {
		for _, dir := range ignoreDirList {
			if path == dir {
				return filepath.SkipDir
			}
		}
		return nil
	}
	return handleFile(path)
}

func handleFile(filename string) error {
	var handleFileList = [...]string{
		".md", ".markdown", ".txt", ".java", ".groovy", ".go", ".c", ".cpp", ".py",
	}
	for _, fileType := range handleFileList {
		if strings.HasSuffix(filename, fileType) {
			if wordDetail == 3 {
				param := []string{"", "", filename, analysisKey}
				analysisTotalChar(param)
			} else {
				countNumByFile(filename)
			}
			return nil
		}
	}
	return nil
}

func readFileAsBytes(fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}
	defer file.Close()
	//读取文件全部内容
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("occur error: ", err)
		return nil
	}
	return content
}

// 根据文件路径读取文件字节
func countNumByFile(fileName string) {
	var content = readFileAsBytes(fileName)
	var total = countChineseChar(content, nil, "")
	// printf("%-50v-> chinese char =\033[0;32m %v\033[0m \n", fileName, total)
	if wordDetail == 0 {
		printf("%-5v %v %-5v %v %v \n", total, green, total, end, fileName)
	}
	totalFile += total
}

// 根据字节流  统计中文字符
func countChineseChar(origin []byte, countChar func(string, string), keyName string) int {
	var total = 0
	temp := [3]byte{}
	var count = 0
	for _, item := range origin {
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
			if wordDetail == 1 {
				printf(string(temp[0:3]))
			}
			if countChar != nil {
				countChar(keyName, string(temp[0:3]))
			}
			total ++
			count = 0
		}
	}
	return total
}

func analysisTotalChar(params []string) {
	fileName := params[2]
	keyName := params[3]
	fileAsBytes := readFileAsBytes(fileName)
	countChineseChar(fileAsBytes, increaseCharNum, keyName)
}

func increaseCharNum(keyName string, CNChar string) {
	result, e := client.ZIncrBy(keyName, 1, CNChar).Result()
	if e == redis.Nil {
		println(result)
	}
}

func showCharRank(start int64, stop int64) {
	result, e := client.ZRevRangeWithScores(analysisKey, start, stop).Result()
	if e != nil {
		println("occur error ")
		return
	}
	for _, a := range result {
		printf("%-6v -> %v \n", a.Score, a.Member)
	}
}

func help(){
	var format = "%-5v %-10v %v \n"
	printf(format, "-h", "", "帮助")
	printf(format, "-w", "", "输出所有汉字")
	printf(format, "-s", "", "简洁输出总字数")
	printf(format, "-all", "", "统计字数,列出排行")
	printf(format, "-del", "", "删除排行数据")
	printf(format, "-show", "start stop", "近列出排行(redis中的 zset 结构)")
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
			analysisTotalChar(param)
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
			client.Del(analysisKey)
			return
		}
	}

	// 递归遍历目录
	filepath.Walk("./", handlerDir)
	printf("\n %v Total characters of all files : %v %v %v \n", purple, yellow, totalFile, end)

	if wordDetail != 3 {
		return
	}
	showCharRank(0, 10)
}
