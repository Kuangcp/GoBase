package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/kuangcp/logger"
	"strings"
)

func main() {
	concatNumber()
}

func concatNumber() {
	lines := readClipList()
	result := strings.Join(lines, ",")
	fmt.Println(result)
	err := clipboard.WriteAll(result)
	if err != nil {
		logger.Fatal(err)
	}
}

func readClipList() []string {
	last, err := clipboard.ReadAll()
	if err != nil {
		logger.Fatal(err)
	}

	last = strings.TrimSpace(last)
	if last == "" {
		logger.Fatal("Empty")
	}
	lines := strings.Split(last, "\n")
	return lines
}
