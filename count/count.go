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
	"github.com/kuangcp/gobase/cuibase"
)

// WARN: 只对 采用 UTF-8 编码的文件保证统计正确

var totalFile = 0
var totalCNChar = 0
var fileList = list.New()

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
func handleFile(fileName string, handleCNChar func(string)) int {
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
		fmt.Printf("%-6v %v->%v %v \n", a.Score, cuibase.Green, cuibase.End, a.Member)
	}
}

func help(params []string) {
	cuibase.PrintTitleDefault("Count chinese char from file that current dir recursive")
	format := cuibase.BuildFormat(-5, -15)
	cuibase.PrintParams(format, []cuibase.ParamInfo{
		{
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
			Comment: "count chinese char for target file, save in redis",
		}, {
			Verb:    "-all",
			Param:   "showNum",
			Comment: "count, calculate rank data, save in redis",
		}, {
			Verb:    "-del",
			Param:   "",
			Comment: "del redis rank data",
		}, {
			Verb:    "-show",
			Param:   "start stop",
			Comment: "show rank data from redis ",
		},
	})
}

func countWithRedis() {
	// 递归遍历目录 读取所有文件
	err := filepath.Walk("./", handlerDir)
	if err != nil {
		log.Print(err)
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
		log.Print(err)
	}
	for e := fileList.Front(); e != nil; e = e.Next() {
		fileName := e.Value.(string)
		if isNeedHandle(fileName) {
			totalFile++
			var total = 0
			total = handleFile(fileName, handler)
			totalCNChar += total
			if showFileInfo {
				fmt.Printf("%-3v %-5v %v%-5v %v%v \n", totalFile, totalCNChar, cuibase.Green, total, cuibase.End, fileName)
			}
		}
	}
	fmt.Printf("\nTotal characters. files: %v%v%v chars: %v%v%v\n", cuibase.Yellow, totalFile, cuibase.End, cuibase.Yellow, totalCNChar, cuibase.End)
}

func readTargetFile(params []string) {
	cuibase.AssertParamCount(3, "Please input param: filename keyName")

	fileName := params[2]
	charRankKey = params[3]

	log.Printf("%v read file: %v, redis key: %v %v ", cuibase.Green, fileName, charRankKey, cuibase.End)
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
