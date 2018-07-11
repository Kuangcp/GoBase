package main

import (
	"fmt"
	"os"
	"io/ioutil"
)

func show_char(origin []byte){
	temp := [3]byte{}
	var count int = 0	

	for _, item := range origin {

		// if item >= 161 && item <= 254{
		if item > 127{
			temp[count] = item
			count++
		}
		if count == 3 {
			fmt.Println("ch char : ", temp, string(temp[0:3]))
			count = 0
		}
	}
}

func main() {
	f, err := os.Open("a.md")
	if err != nil {
		fmt.Println("os Open error: ", err)
		return
	}
	defer f.Close()

	//读取文件全部内容
	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("ioutil ReadAll error: ", err)
		return
	}
	// fmt.Println("content: ", b)
	show_char(b)
}