package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"
)
// WARN: 该程序只对 采用 UTF-8 编码的文件保证统计正确

// TODO 分析用字的数据, 得到一个统计报表

var totalFile int = 0
var wordDetail bool = false
var printf = fmt.Printf
var println = fmt.Println

// 往递归遍历目录 作为参数传入的函数 
func handlerDir(path string, info os.FileInfo, err error) error {
	if(info.IsDir()){
		return nil
	}
	var ignoreList = [...]string{
		".md",
		".markdown", 
		".txt", 
	}
	for _, fileType := range ignoreList {
		if strings.HasSuffix(path, fileType) {
			readfile(path)
			return nil
		}
	}
    return nil
}
// 根据文件路径读取文件字节
func readfile(fileName string){
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("os Open error: ", err)
		return
	}
	defer file.Close()
	//读取文件全部内容
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("ioutil ReadAll error: ", err)
		return
	}
	var total int = countChineseChar(content)
	// printf("%-50v-> chinese char =\033[0;32m %v\033[0m \n", fileName, total)
	if !wordDetail{
		printf(">\033[0;32m %-5v\033[0m %v\n", total, fileName)
	}
	totalFile += total
}
// 计算中文字符
func countChineseChar(origin []byte) int{
	var total int = 0
	temp := [3]byte{}
	var count int = 0
	for _, item := range origin {
		// 满足 110xxxxx 4e00 9fa5  头部分别为 0100 1001 
		if count == 0 && item >= 228 && item <= 233{
			temp[count] = item
			count++
		}
		// 满足 10xxxxxx
		if count != 0 && item >= 128 && item <= 191{
			temp[count] = item 
			count++
		}
		if count == 3 {
			// fmt.Println("ch char : ", temp, string(temp[0:3]))
			if wordDetail {
				printf(string(temp[0:3]))
			}
			total ++
			count = 0
		}
	}
	return total
}
// TODO 分析 所有汉字
func analysisTotalChar(){

}
// 读取参数 
func handlerArgs(verb string, param string){
	switch verb {
		case "-h":
			var format string = "%-5v %-10v %v \n"
			printf(format, "-h", "", "帮助")
			printf(format, "-w", "", "输出所有汉字")
			os.Exit(0)
		case "-w":
			wordDetail = true
			break
	}
}

func main() {
	// 递归遍历目录
	var argLen int = len(os.Args)
	// printf("len of args:%d\n", argLen)
	// for i, v := range os.Args {
	// 	printf("args[%d]=%s\n", i, v)
	// }

	if argLen == 2 {
		handlerArgs(os.Args[1], "")
	}
	if argLen == 3 {
		handlerArgs(os.Args[1], os.Args[2])
	}
	filepath.Walk("./", handlerDir)
	printf("\n\033[0;35m Total characters of all files : \033[0;33m %v \033[0m \n", totalFile)
}