package playground

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"strconv"
	"strings"
	"testing"
)

// 统计 http://localhost:7000/histo/ 堆大小
func TestCountHeapSize(t *testing.T) {
	sum := 0
	lines := ctool.ReadStrLinesNoFilter("histo.log")
	for _, l := range lines {
		fields := strings.Fields(l)
		if len(fields) > 3 {
			//fmt.Println(fields[3])
			atoi, _ := strconv.Atoi(fields[3])
			sum += atoi
		}
	}
	fmt.Println(sum, sum>>10, sum>>20)
}
