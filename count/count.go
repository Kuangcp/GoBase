package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"
)
// 汉字以及中文标点 Unicode : 0x4E00 - 0x9FA5 
// 该程序只对 UTF-8 编码的文件保证统计正确

// TODO 去除中文标点
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
	fmt.Println(fileName, " total char : ", total)
	totalFile += total
}
func showchar(origin []byte) int{
	var total int = 0
	temp := [3]byte{}
	var count int = 0
	for _, item := range origin {

		// 满足 110xxxxx
		if count == 0 && item >= 224 && item <= 239{
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
	fmt.Println("\nall file char : ", totalFile)
}