package main

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	start := time.Now()
	for i := 0; i < 1000; i++ {
		rePackage = regexp.MustCompile(`(?m)^\s*package\s+([\w.]+)\s*;`)
		reImport = regexp.MustCompile(`(?m)^\s*import\s+(?:static\s+)?([\w.]+)(?:\.\*)?\s*;`)
		reWord = regexp.MustCompile(`\b([A-Z][a-zA-Z0-9]*)\b`)      // 简单识别首字母大写的标识符（类名）
		reController = regexp.MustCompile(`@(?:Rest)?Controller\b`) // @Controller 或 @RestController
	}
	fmt.Println(time.Now().Sub(start))
}
