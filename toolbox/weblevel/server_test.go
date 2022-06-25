package weblevel

import (
	"fmt"
	"log"
	"testing"
)

var client Client

func init() {
	client = NewClient("localhost", 33742)
}

func TestServer(t *testing.T) {
	levelDB, err := NewServer(&Options{Port: 9066, DBPath: "test-db"})
	if err != nil {
		log.Println(err)
		return
	}
	levelDB.Bootstrap()
}

func TestRange(t *testing.T) {
	search := client.PrefixSearch("page:")
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
	val, err := client.Get("test-a")
	if err != nil {
		t.Failed()
	}
	fmt.Println(val)
}
