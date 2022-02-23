package main

import (
	"fmt"
	"github.com/kuangcp/logger"
	"io/ioutil"
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {
	open, err := os.Open("/home/kcp/test/ss/c/")
	if err != nil {
		fmt.Println(err)
		return
	}
	logger.Info(open)

	dir, err := ioutil.ReadDir("/home/kcp/test/ss/c/")
	if err != nil {
		fmt.Println(err)
		return
	}
	logger.Info(dir)
}
