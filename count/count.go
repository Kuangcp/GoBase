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

func walkfunc(path string, info os.FileInfo, err error) error {
	if(info.IsDir()){
		return nil
	}
	if(strings.Contains(path, ".git")){
		return nil
	}
	// fmt.Println(path)
	readfile(path)
    return nil
}
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
	var total int = showchar(content)
	// fmt.Println("\n ", fileName, " total char : ", total)
	// fmt.Println(fileName, " -> chinese char : ", total)
	fmt.Printf("%-60v -> chinese char = \033[0;32m %v \033[0m \n", fileName, total)
	totalFile += total
}
func showchar(origin []byte) int{
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
			// fmt.Printf(string(temp[0:3]))
			total ++
			count = 0
		}
	}
	return total
}

func main() {
	filepath.Walk("./", walkfunc)
	fmt.Println()
	fmt.Println("all file chinese char : ", totalFile)
}