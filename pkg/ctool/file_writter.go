package ctool

import (
	"bufio"
	"fmt"
	"os"
)

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
		fileMod = os.O_WRONLY | os.O_TRUNC | os.O_CREATE
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
