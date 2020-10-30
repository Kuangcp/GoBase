package main

import (
	"fmt"
	"sync"
)

func input(ch chan string) {
	defer wg.Done()
	defer close(ch)
	var input string
	fmt.Println("Enter 'EOF' to shut down: ")
	for {
		_, err := fmt.Scanf("%s", &input)
		if err != nil {
			fmt.Println("Read input err: ", err.Error())
			break
		}
		if input == "EOF" {
			fmt.Println("Bye!")
			break
		}
		ch <- input
	}
}

func output(ch chan string) {
	defer wg.Done()
	for value := range ch {
		fmt.Println("Your input: ", value)
	}
}

var wg sync.WaitGroup

func main() {
	ch := make(chan string)
	wg.Add(2)
	go input(ch)  // 读取输入
	go output(ch) // 输出到命令行
	wg.Wait()
	fmt.Println("Exit!")
}
