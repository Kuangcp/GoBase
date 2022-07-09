package ctool

import "testing"

func TestTsvFile(t *testing.T) {
	writer, _ := NewWriter("test.tsv", true)
	writer.WriteLine("tt\t1\t9")
	writer.WriteLine("tt\t1\t9\t44\t899")
	defer writer.Close()
}
