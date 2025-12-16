package main

import (
	"flag"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"
)

var (
	pmapFile    string
	onlyArena   bool
	onlySummary bool
)

func main() {
	flag.StringVar(&pmapFile, "pmap", "", "pmap file to use")
	flag.BoolVar(&onlyArena, "o", false, "only arena")
	flag.BoolVar(&onlySummary, "s", false, "only summary")
	flag.Parse()

	var last []string
	cnt := 0
	var lastRow string

	var lines []string
	if pmapFile != "" {
		lines = ctool.ReadStrLines(pmapFile, nil)
	} else {
		stat, _ := os.Stdin.Stat()
		if stat.Mode()&os.ModeNamedPipe == os.ModeNamedPipe {
			// 从管道读取全部内容
			bytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println("读取管道内容失败:", err)
				return
			}
			all := string(bytes)
			//fmt.Println("从管道读取的内容:", all)
			xt := strings.Split(all, "\n")
			//fmt.Println(result)
			for _, o := range xt {
				lines = append(lines, o+"\n")
			}
		}
	}
	//fmt.Println(len(lines))
	start := false
	for _, row := range lines {
		if !start {
			start = strings.HasPrefix(row, "Address")
			continue
		}
		cols := strings.Fields(row)
		if len(cols) == 0 {
			continue
		}
		if len(last) == 0 {
			last = cols
			lastRow = row
			continue
		}
		curAddr := cols[0]
		preAddr := last[0]

		a := new(big.Int)
		b := new(big.Int)
		a.SetString(curAddr, 16)
		b.SetString(preAddr, 16)

		dif := new(big.Int).Sub(a, b)
		blockKib := dif.Int64() / 1024

		//fmt.Println("_____", curAddr, preAddr, blockKib)

		p := false
		if len(last) > 1 && fmt.Sprint(blockKib) == last[1] {
			mm, _ := strconv.Atoi(cols[1])
			if int64(mm)+blockKib == int64(65536) {
				if !onlySummary {
					fmt.Printf("%sArena %s%s", ctool.Cyan, lastRow, ctool.White)
					fmt.Printf("%sArena %s%s", ctool.Cyan, row, ctool.White)
					fmt.Println()
				}
				p = true
				cnt += 1
			}
		}
		if !p && !onlyArena && !onlySummary {
			fmt.Printf("%s%s%s", ctool.White, row, ctool.White)
		}

		//fmt.Println(strings.Join(cols, "_"))
		last = cols
		lastRow = row
	}
	fmt.Println("count:", cnt)
}
