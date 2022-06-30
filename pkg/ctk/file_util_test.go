package ctk

import (
	"fmt"
	"testing"
)

func TestWriteFile(t *testing.T) {
	writer, err := NewWriter("c.log", false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer writer.Close()

	for i := 0; i < 5; i++ {
		writer.WriteString(fmt.Sprint(i) + "  ")
	}
}

func TestIgnoreError(t *testing.T) {
	writer := NewWriterIgnoreError("c.log", true)
	defer writer.Close()

	writer.WriteLine("xxxxxxx")
}
