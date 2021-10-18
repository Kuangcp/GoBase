package cuibase

import (
	"fmt"
	"testing"
)

func TestWriteFile(t *testing.T) {
	writer, err := NewWriter("b.log")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer writer.Close()

	for i := 0; i < 10000000; i++ {
		writer.writer.Write([]byte(fmt.Sprint(i) + "\n"))
	}
}
