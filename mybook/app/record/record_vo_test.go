package record

import (
	"fmt"
	"mybook/app/common/util"
	"testing"
	"time"
)

func TestCopyCompare(t *testing.T) {
	start := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		entity := RecordEntity{
			BookId:     33,
			AccountId:  444,
			TransferId: 33333,
			Amount:     93734,
			Comment:    "test some data like this",
			CategoryId: 128476,
			RecordTime: time.Now(),
		}
		var target RecordEntity
		util.Copy(entity, &target)
	}
	fmt.Println((time.Now().UnixNano()-start)/1000_000, "ms")
}
