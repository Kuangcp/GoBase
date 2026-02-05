package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// 1. 同步阻塞：顺序执行，阻塞等待
func fetchSync(url string) string {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func SyncMain() {
	start := time.Now()
	fmt.Println(fetchSync("http://localhost:8080/s"))
	fmt.Println(fetchSync("http://localhost:8080/s"))
	fmt.Println(fetchSync("http://localhost:8080/s"))
	fmt.Printf("总耗时: %v\n", time.Since(start)) // ~6s
}

// 2. 异步非阻塞：协程内部并发，主 goroutine 线性收结果
func fetchAsync(url string) <-chan string {
	out := make(chan string)
	go func() {
		resp, _ := http.Get(url)
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		out <- string(body)
	}()
	return out
}

func AsyncMain() {
	start := time.Now()
	// 1. 并发发起请求（goroutine 内部并发）
	ch1 := fetchAsync("http://localhost:8080/s")
	ch2 := fetchAsync("http://localhost:8080/s")
	ch3 := fetchAsync("http://localhost:8080/s")

	// 2. 主 goroutine 并发收结果（非阻塞）
	fmt.Println(<-ch1)
	fmt.Println(<-ch2)
	fmt.Println(<-ch3)
	fmt.Printf("总耗时: %v\n", time.Since(start))
}

func main() {
	go func() {
		http.HandleFunc("/s", func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(time.Duration(2) * time.Second)
			writer.Write([]byte("ok"))
		})
		http.ListenAndServe(":8080", nil)
	}()
	time.Sleep(time.Duration(1) * time.Second)
	fmt.Println("start")
	SyncMain()
	AsyncMain()
}
