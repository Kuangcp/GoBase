package main

import (
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/kuangcp/gobase/pkg/sizedpool"
	"github.com/kuangcp/logger"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	portStr string
	hostStr string
	con     int
	sorted  bool

	total int
)

func init() {
	flag.StringVar(&portStr, "p", "80", "port")
	flag.StringVar(&hostStr, "h", "", "host")
	flag.IntVar(&con, "c", 0, "parallel count")
	flag.BoolVar(&sorted, "s", false, "sorted output")
}

func main() {
	flag.Parse()
	logger.Info("Start", hostStr, portStr)
	start, end, err := parsePort()
	if err != nil {
		return
	}
	total = end - start + 1

	if con != 0 {
		parallelScan(start, end)
		return
	}

	if end != 0 {
		bar := pb.StartNew(total)
		for i := start; i <= end; i++ {
			bar.Increment()
			pingS(hostStr, fmt.Sprint(i))
		}
		bar.Finish()
	} else {
		pingS(hostStr, fmt.Sprint(start))
	}
}

func parallelScan(start, end int) {
	noLimit := con == -1

	var doPing = func(i int, ps chan int) {
		if sorted {
			if ping(hostStr, fmt.Sprint(i)) {
				ps <- i
			}
		} else {
			pingS(hostStr, fmt.Sprint(i))
		}
	}

	// 并不一定无限制的效率更高，创建太多协程反而使得调度更耗时吞吐量下降明细
	if noLimit {
		var ps chan int
		if sorted {
			ps = make(chan int, end-start+1)
		}

		var w sync.WaitGroup
		w.Add(total)
		for i := start; i <= end; i++ {
			i := i
			go func() {
				doPing(i, ps)
				w.Done()
			}()
		}
		w.Wait()
		readResult(ps)
	} else {
		var ps chan int
		if sorted {
			ps = make(chan int, end-start+1)
		}
		group, _ := sizedpool.New(sizedpool.PoolOption{Size: con})
		for i := start; i <= end; i++ {
			i := i
			group.Run(func() {
				doPing(i, ps)
			})
		}

		group.Wait()
		readResult(ps)
	}
}

func readResult(ps chan int) {
	if !sorted {
		return
	}

	close(ps)
	var res []int
	for p := range ps {
		res = append(res, p)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})
	for _, re := range res {
		fmt.Printf("%v\t%s\t\n", re, "open")
	}
}

func parsePort() (int, int, error) {
	portRange := strings.Contains(portStr, "-")
	if portRange {
		pair := strings.Split(portStr, "-")
		start, err := strconv.Atoi(pair[0])
		if err != nil {
			logger.Error(err)
			return 0, 0, err
		}
		end, err := strconv.Atoi(pair[1])
		if err != nil {
			logger.Error(err)
			return 0, 0, err
		}
		return start, end, nil
	} else {
		start, err := strconv.Atoi(portStr)
		if err != nil {
			logger.Error(err)
			return 0, 0, err
		}
		return start, 0, nil
	}
}

func pingS(host, port string) {
	if ping(host, port) {
		fmt.Printf("%v\t%s\t\n", port, "open")
	}
}

func ping(host, port string) bool {
	addr := fmt.Sprintf("%s:%s", host, port)
	_, err := net.Dial("tcp", addr)
	if err == nil {
		//fmt.Printf("%v\t%s\t\n", port, "open")
		return true
	} else {
		//fmt.Printf("%v\t%s\t\n", port, "closed")
		return false
	}
}
