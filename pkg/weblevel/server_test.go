package weblevel

import (
	"crypto/md5"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"log"
	"testing"
	"time"
)

var client Client

func init() {
	client = NewClient("localhost", 33742)
}

// 对比内存中目录和磁盘目录性能差异
// sudo mkdir /mnt/tmp
// sudo mount -t tmpfs -o size=100m tmpfs /mnt/tmp
func TestDiffServer(t *testing.T) {
	memDB, err := NewServer(&Options{Port: 33740, DBPath: "/mnt/tmp/mem-db"})
	if err != nil {
		log.Println(err)
		return
	}
	go memDB.Bootstrap()

	levelDB, err := NewServer(&Options{Port: 33741, DBPath: "disk-db"})
	if err != nil {
		log.Println(err)
		return
	}
	levelDB.Bootstrap()
}

func putData(cli *WebClient) {
	for i := 0; i < 900000; i++ {
		iStr := fmt.Sprint("do", i)
		sum := md5.Sum([]byte(iStr))
		cli.Set(iStr, fmt.Sprintf("%x", sum))
	}
}

// 写入有点点差异，读取
//
//	52.812s 46% mem write
//	57.363s 50% disk write
//	 1947ms  1% mem
//	 1887ms  1% disk

// dd if=/dev/zero bs=1M count=50 of=b
//50+0 records in
//50+0 records out
//52428800 bytes (52 MB, 50 MiB) copied, 0.0208666 s, 2.5 GB/s

// dd if=/dev/zero bs=1M count=50 of=b
//50+0 records in
//50+0 records out
//52428800 bytes (52 MB, 50 MiB) copied, 0.0120049 s, 4.4 GB/s
func TestInitDiffServer(t *testing.T) {
	memCli := NewClient("localhost", 33740)
	diskCli := NewClient("localhost", 33741)

	watch := ctool.NewStopWatch()

	watch.Start("mem write")
	putData(memCli)
	watch.Stop()

	watch.Start("disk write")
	putData(diskCli)
	watch.Stop()

	watch.Start("mem")
	memCli.PrefixSearch("do")
	watch.Stop()

	watch.Start("disk")
	diskCli.PrefixSearch("do")
	watch.Stop()
	fmt.Println(watch.PrettyPrint())
}

func TestServer(t *testing.T) {
	levelDB, err := NewServer(&Options{Port: 33742, DBPath: "test-db"})
	if err != nil {
		log.Println(err)
		return
	}
	levelDB.Bootstrap()
}

func TestRange(t *testing.T) {
	search := client.PrefixSearch("do")
	for k, v := range search {
		fmt.Println(k, v)
	}
}

func TestWebClient_Sets(t *testing.T) {
	client.Sets(map[string]string{
		"test-a": "bbb",
		"test-b": "bbb",
	})
}

func TestWebClient_Get(t *testing.T) {
	val, err := client.Get("test-a1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
}

func TestRollingTxt(t *testing.T) {
	s := "202206"

	runes := []rune(s)
	showWin := 20
	cursor := 0
	l := len(runes)

	for i := 0; i < 200; i++ {
		cursor = cursor % l
		v2 := cursor + showWin
		if v2 > l {
			fmt.Print(string(runes[cursor:]) + "   " + string(runes[:v2-l]) + "\r\r")
		} else {
			fmt.Print(string(runes[cursor:v2-1]) + "\r\r")
		}
		cursor++
		time.Sleep(time.Millisecond * 150)
	}
}
