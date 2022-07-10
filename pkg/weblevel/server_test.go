package weblevel

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var client Client

func init() {
	client = NewClient("localhost", 33742)
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
	s := "20220626_021829_Masquerade-マスカレード-葵つかさアサ芸SEXY女優写真集[172P]_172P"

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
