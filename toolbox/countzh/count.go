package main

import (
	"container/list"
	"flag"
	"fmt"
	"github.com/kuangcp/gobase/pkg/cuibase"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v7"
)

var totalFile = 0
var totalChineseChar = 0
var fileList = list.New()

var client *redis.Client

func initRedisClient() {
	config := readRedisConfig()
	if config == nil {
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
func handleFile(fileName string, chineseCharHandler func(string)) int {
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
			if chineseCharHandler != nil {
				chineseCharHandler(string(temp[0:3]))
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
		fmt.Printf("%-6v %s %v \n", a.Score, cuibase.Green.Print("->"), a.Member)
	}
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

func showChineseChar(perCharHandler func(string), printFileInfo bool) {
	if printFileInfo {
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
		total = handleFile(fileName, perCharHandler)
		totalChineseChar += total
		if printFileInfo {
			fmt.Printf("%-3v %-5v %s %v \n",
				totalFile, totalChineseChar, cuibase.Green.Printf("%-5v", total), fileName)
		}
	}
	fmt.Printf("\nTotal characters: %v%v%v files  %v%v%v chars \n",
		cuibase.Yellow, totalFile, cuibase.End, cuibase.Yellow, totalChineseChar, cuibase.End)
}

func delRank() {
	initRedisClient()
	client.Del(charRankKey)
	log.Printf("del %v%v%v", cuibase.Green, charRankKey, cuibase.End)
}

func printRank() {
	nums := strings.Split(rankPair, ",")
	if len(nums) == 1 {
		end, err := strconv.ParseInt(nums[0], 10, 64)
		cuibase.CheckIfError(err)
		endIdx = end
	} else if len(nums) == 2 {
		start, err1 := strconv.ParseInt(nums[0], 10, 64)
		cuibase.CheckIfError(err1)
		stop, err2 := strconv.ParseInt(nums[1], 10, 64)
		cuibase.CheckIfError(err2)
		startIdx = start
		endIdx = stop
	} else {
		fmt.Println("error rank format eg: 1,2")
		return
	}

	initRedisClient()
	showCharRank(startIdx, endIdx)
}

func init() {
	flag.StringVar(&handleFileSuffix, "S", "", "")
	flag.StringVar(&ignoreDir, "D", "", "")
	flag.StringVar(&rankPair, "r", "", "")
	flag.StringVar(&targetFile, "f", "", "")
}

func main() {
	info.Parse()
	if help {
		info.PrintHelp()
		return
	}
	if delCache {
		delRank()
		return
	}

	if handleFileSuffix != "" {
		handleFiles = strings.Split(handleFileSuffix, ",")
	}
	if ignoreDir != "" {
		ignoreDirs = strings.Split(handleFileSuffix, ",")
	}

	if rankPair != "" {
		printRank()
		return
	}

	if targetFile != "" {
		log.Printf("%vread file: %v, redis key: %v%v\n",
			cuibase.Green, targetFile, charRankKey, cuibase.End)
		initRedisClient()
		handleFile(targetFile, increaseCharNum)

		showCharRank(0, 15)
		return
	}

	if allFile {
		initRedisClient()
		countWithRedis()
		showCharRank(0, 15)
		return
	}

	if printAllChar {
		showChineseChar(showChar, false)
		return
	}
	if countSummary {
		showChineseChar(nil, false)
		return
	}

	showChineseChar(nil, true)
}
