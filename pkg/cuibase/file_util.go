package cuibase

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// ReadFileLines
func ReadFileLines(filename string, filterFunc func(string) bool, mapFunc func(string) interface{}) []interface{} {
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
		if filterFunc == nil || filterFunc(line) {
			if mapFunc != nil {
				result = append(result, mapFunc(line))
			} else {
				result = append(result, line)
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Println("Read file error!", err)
				return nil
			}
		}
	}
	return result
}

func ReadFileLinesNoFilter(filename string, mapFunc func(string) interface{}) []interface{} {
	return ReadFileLines(filename, nil, mapFunc)
}

type BufferWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func NewWriter(filePath string) (*BufferWriter, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
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

func (w *BufferWriter) Close() {
	w.writer.Flush()
	w.file.Close()
}
