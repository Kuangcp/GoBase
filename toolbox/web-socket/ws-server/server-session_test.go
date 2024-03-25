package main

import (
	"fmt"
	"log"
	"regexp"
	"testing"
)

func TestMatch(t *testing.T) {
	compile, err := regexp.Compile("((\\w+)( |(\\.\\.\\.)|(-))?)+\\w")
	if err != nil {
		log.Println(err)
	}
	result := compile.FindAllString("aaa-bbb 是该 llllll 模式 vvvvv 校 ooooooo", -1)

	for i, sub := range result {
		fmt.Println(i, sub)
	}
}
