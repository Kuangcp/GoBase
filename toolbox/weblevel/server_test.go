package weblevel

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
)

var client *WebClient

func init() {
	client = NewClient("localhost", 9066)
}

func TestServer(t *testing.T) {
	newDB, err := leveldb.OpenFile("test-db", nil)
	if err != nil {
		return
	}
	levelDB := NewServer(newDB, 9066)
	levelDB.Bootstrap()
}

func TestRange(t *testing.T) {
	search := client.PrefixSearch("test")
	for k, v := range search {
		fmt.Println(k, v)
	}
}
