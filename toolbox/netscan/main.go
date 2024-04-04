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
	"time"
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
		fmt.Println("Usage: ")
		fmt.Println("    netscan -i 192.168.1.0/24 -c 100")
		fmt.Println("    netscan -i 192.168.1.3")
		fmt.Println("    netscan -i 192.168.1.3 -p 443")
		fmt.Println("    netscan -i 192.168.1.3 -P -c 10")
		fmt.Println("    netscan -i 192.168.1.3 -p 5000-8000 -c 10")
		return
	}

	logger.Info("Start scan", hostStr, portStr)
	start, end, err := parsePort()
	if err != nil {
		return
	}
	total = end - start + 1

	if strings.Contains(hostStr, "/") {
		scanHostRange(end, start)
		return
	}

	var bar *pb.ProgressBar
	if end != 0 {
		bar = pb.Full.Start(total)
	}
	scanHost(start, end, hostStr, bar)
	if bar != nil {
		bar.Finish()
	}
}

func scanHostRange(end int, start int) {
	_, ipNet, err := net.ParseCIDR(hostStr)
	if err != nil {
		logger.Error(err)
		return
	}

	startHost, endHost := AddressRange(ipNet)
	logger.Info("scan range:", startHost, endHost)
	cur := startHost
	a, _ := sizedpool.New(sizedpool.PoolOption{Size: 100})

	var ips []string
	for !cur.Equal(endHost) {
		cur = Inc(cur)
		ips = append(ips, cur.String())
	}
	var bar *pb.ProgressBar
	if end != 0 {
		bar = pb.Full.Start(total * len(ips))
	}
	for _, host := range ips {
		a.Run(func() {
			//fmt.Println("======", cur.String())
			scanHost(start, end, host, bar)
		})
	}
	a.Wait()
	if bar != nil {
		bar.Finish()
	}
	return
}

func scanHost(start, end int, host string, bar *pb.ProgressBar) {
	if con != 0 {
		parallelScan(start, end, host, bar)
		return
	}

	if end == 0 {
		pingS(host, fmt.Sprint(start))
	} else {
		var ps = make(chan int, total)
		for i := start; i <= end; i++ {
			if ping(host, fmt.Sprint(i)) {
				ps <- i
			}
			bar.Increment()
		}
		readResult(host, ps)
	}
}

func parallelScan(start, end int, host string, bar *pb.ProgressBar) {
	noLimit := con == -1

	// 并不一定无限制开启协程的方式会效率更高，创建太多协程反而使得调度更耗时吞吐量下降明显
	// 但是目标端网络延迟高且大量端口未开启时创建大量协程的效率会更高，因为大量的协程会阻塞在IO 出让了CPU执行，相当于网络延迟被忽视了
	if noLimit {
		var ps = make(chan int, 1000)
		var w sync.WaitGroup
		w.Add(total)
		for i := start; i <= end; i++ {
			i := i
			go func() {
				if ping(host, fmt.Sprint(i)) {
					ps <- i
				}
				bar.Increment()
				w.Done()
			}()
		}
		w.Wait()
		readResult(host, ps)
	} else {
		var ps = make(chan int, 1000)
		group, _ := sizedpool.New(sizedpool.PoolOption{Size: con})
		for i := start; i <= end; i++ {
			i := i
			group.Run(func() {
				if ping(host, fmt.Sprint(i)) {
					ps <- i
				}
				bar.Increment()
			})
		}

		group.Wait()
		readResult(host, ps)
	}
}

func readResult(host string, ps chan int) {
	close(ps)
	var res []int
	for p := range ps {
		res = append(res, p)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})
	for _, re := range res {
		fmt.Printf("%v %v\t%s\t\n", host, re, "open")
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
		fmt.Printf("%v %v %s \n", host, port, "open")
	}
}

func ping(host, port string) bool {
	addr := fmt.Sprintf("%s:%s", host, port)
	_, err := net.DialTimeout("tcp", addr, time.Second*1)
	return err == nil
}
