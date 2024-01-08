package ctool

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

var self = func(s string) string {
	return s
}

// ReadLines filter and map
func ReadLines[T any](filename string, filterFunc func(string) bool, mapFunc func(string) T) []T {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		log.Println("Open file error!", err)
		return nil
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if stat.Size() == 0 {
		log.Printf("file:%s is empty", filename)
		return nil
	}

	var result []T

	buf := bufio.NewReader(file)
	end := false
	for !end {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				end = true
			} else {
				log.Println("Read file error!", err)
				return nil
			}
		}

		if filterFunc != nil && !filterFunc(line) {
			continue
		}

		if mapFunc != nil {
			val := mapFunc(line)
			result = append(result, val)
		}
	}
	return result
}

func ReadStrLines(filename string, filterFunc func(string) bool) []string {
	return ReadLines(filename, filterFunc, self)
}

func ReadStrLinesNoFilter(filename string) []string {
	return ReadLines(filename, nil, self)
}

// ReadLinesNoFilter only map
func ReadLinesNoFilter[T any](filename string, mapFunc func(string) T) []T {
	return ReadLines(filename, nil, mapFunc)
}

func readSplitLine(filename string, split string) [][]string {
	lines := ReadLines(filename, nil, func(s string) []string {
		return strings.Split(strings.ReplaceAll(strings.TrimSpace(s), "\"", ""), split)
	})

	return lines
}

func ReadCsvLines(filename string) [][]string {
	return readSplitLine(filename, ",")
}
func ReadTsvLines(filename string) [][]string {
	return readSplitLine(filename, "\t")
}

// IsFileExist relative or absolute path
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
