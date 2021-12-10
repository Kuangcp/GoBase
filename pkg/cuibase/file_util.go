package cuibase

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// ReadFileLines filter and map
func ReadFileLines(filename string,
	filterFunc func(string) bool,
	mapFunc func(string) interface{}) []interface{} {
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

	var result []interface{}

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
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
			if val == nil {
				continue
			}
			result = append(result, val)
		} else {
			if line == "" {
				continue
			}
			result = append(result, line)
		}
	}
	return result
}

// ReadFileLinesNoFilter only map
func ReadFileLinesNoFilter(filename string, mapFunc func(string) interface{}) []interface{} {
	return ReadFileLines(filename, nil, mapFunc)
}

type BufferWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func NewWriterIgnoreError(filePath string, truncate bool) *BufferWriter {
	writer, err := NewWriter(filePath, truncate)
	if err != nil {
		fmt.Println(err)
	}
	return writer
}

func NewWriter(filePath string, truncate bool) (*BufferWriter, error) {
	var fileMod int
	if truncate {
		fileMod = os.O_WRONLY | os.O_TRUNC
	} else {
		fileMod = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	}

	file, err := os.OpenFile(filePath, fileMod, 0666)
	if err != nil {
		return nil, err
	}

	fmt.Println(file.Name())
	buffer := bufio.NewWriter(file)
	return &BufferWriter{file: file, writer: buffer}, nil
}

func (w *BufferWriter) Write(val []byte) (nn int, err error) {
	return w.writer.Write(val)
}

func (w *BufferWriter) WriteString(val string) (nn int, err error) {
	return w.writer.WriteString(val)
}

func (w *BufferWriter) WriteLine(val string) (nn int, err error) {
	return w.writer.WriteString(val + "\n")
}

func (w *BufferWriter) Close() {
	w.writer.Flush()
	w.file.Close()
}
