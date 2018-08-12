package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"
	"github.com/go-redis/redis"
)

// WARN: 该程序只对 采用 UTF-8 编码的文件保证统计正确

var totalFile = 0
var wordDetail = 0 // 默认 0:输出路径, 1:所有字, 2:只有统计 3 分析
var printf = fmt.Printf

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
	var handleFileList = [...]string{
		".md", ".markdown", ".txt", ".java", ".groovy", ".go", ".c", ".cpp", ".py",
	}
	for _, fileType := range handleFileList {
		if strings.HasSuffix(path, fileType) {
			if wordDetail == 3 {
				param := []string{"", "", path, "a"}
				analysisTotalChar(param)
			} else {
				countNumByFile(path)
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

// TODO 分析 所有汉字
func analysisTotalChar(params []string) {
	fileName := params[2]
	keyName := params[3]
	fileAsBytes := readFileAsBytes(fileName)
	countChineseChar(fileAsBytes, addCharNum, keyName)
}

func addCharNum(keyName string, CNChar string) {
	result, e := client.ZIncrBy(keyName, 1, CNChar).Result()
	if e == redis.Nil {
		println(result)
	}
}

// 读取参数, 设置全局变量的值
func handlerArgs(verb string, param []string) {
	switch verb {
	case "-h":
		var format = "%-5v %-10v %v \n"
		printf(format, "-h", "", "帮助")
		printf(format, "-w", "", "输出所有汉字")
		printf(format, "-s", "", "简洁输出总字数")
		os.Exit(0)
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
		os.Exit(0)
	case "-al":
		wordDetail = 3
		break
	}

}

func main() {
	// 递归遍历目录
	var argLen = len(os.Args)
	//printf("len of args:%d\n", argLen)
	//for i, v := range os.Args {
	//	printf("args[%d]=%s\n", i, v)
	//}

	if argLen == 2 {
		handlerArgs(os.Args[1], nil)
	}
	if argLen > 2 {
		handlerArgs(os.Args[1], os.Args)
	}
	filepath.Walk("./", handlerDir)
	printf("\n %v Total characters of all files : %v %v %v \n", purple, yellow, totalFile, end)

	if wordDetail == 3 {
		result, e := client.ZRevRangeWithScores("a", 0, 100).Result()
		for _, a := range result {
			printf("%6v-%v %v \n", a.Score, a.Member, e)
		}
	}
}
