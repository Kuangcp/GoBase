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
	allPort bool
	hostStr string
	con     int
	help    bool

	total int
)

func init() {
	flag.StringVar(&portStr, "p", "80", "port")
	flag.StringVar(&hostStr, "i", "", "host")
	flag.IntVar(&con, "c", 0, "parallel count")
	flag.BoolVar(&allPort, "P", false, "all port")
	flag.BoolVar(&help, "h", false, "help info")
}

func main() {
	flag.Parse()
	if help {
		logger.Info("netscan -i 192.168.1.3")
		logger.Info("netscan -i 192.168.1.3 -p 443")
		logger.Info("netscan -i 192.168.1.3 -P -c 10")
		logger.Info("netscan -i 192.168.1.3 -p 5000-8000 -c 10")
		return
	}

	logger.Info("Start scan", hostStr, portStr)
	start, end, err := parsePort()
	if err != nil {
		return
	}
	total = end - start + 1

	if con != 0 {
		parallelScan(start, end)
		return
	}

	if end == 0 {
		pingS(hostStr, fmt.Sprint(start))
	} else {
		var ps = make(chan int, total)
		bar := pb.StartNew(total)
		for i := start; i <= end; i++ {
			if ping(hostStr, fmt.Sprint(i)) {
				ps <- i
			}
			bar.Increment()
		}
		bar.Finish()
		readResult(ps)
	}
}

func parallelScan(start, end int) {
	noLimit := con == -1

	// 并不一定无限制的效率更高，创建太多协程反而使得调度更耗时吞吐量下降明细
	if noLimit {
		var ps = make(chan int, 1000)

		var w sync.WaitGroup
		w.Add(total)
		bar := pb.Full.Start(total)
		for i := start; i <= end; i++ {
			i := i
			go func() {
				if ping(hostStr, fmt.Sprint(i)) {
					ps <- i
				}
				bar.Increment()
				w.Done()
			}()
		}
		w.Wait()
		bar.Finish()
		readResult(ps)
	} else {
		var ps = make(chan int, 1000)
		group, _ := sizedpool.New(sizedpool.PoolOption{Size: con})
		bar := pb.Full.Start(total)
		for i := start; i <= end; i++ {
			i := i
			group.Run(func() {
				if ping(hostStr, fmt.Sprint(i)) {
					ps <- i
				}
				bar.Increment()
			})
		}

		group.Wait()
		readResult(ps)
		bar.Finish()
	}
}

func readResult(ps chan int) {
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
	if allPort {
		return 1, 65535, nil
	}
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
		fmt.Printf("%v %s \n", port, "open")
	}
}

func ping(host, port string) bool {
	addr := fmt.Sprintf("%s:%s", host, port)
	_, err := net.Dial("tcp", addr)
	return err == nil
}
