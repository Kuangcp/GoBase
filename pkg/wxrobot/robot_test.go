package wxrobot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)
func TestImageBase64(t *testing.T) {
	dir, err := os.ReadDir("/home/kcp/Pictures/")
	if err != nil {
		return
	}
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		file, err := ioutil.ReadFile("/home/kcp/Pictures/" + entry.Name())
		if err != nil {
			continue
		}
		base64 := imgToBase64(file)
		fmt.Println(len(file), len(base64), float32(len(base64))/float32(len(file)))
	}
}
