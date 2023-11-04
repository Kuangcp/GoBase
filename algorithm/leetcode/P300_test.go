package leetcode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// https://leetcode.cn/problems/longest-increasing-subsequence/
func lengthOfLISWithSimple(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	maxL := 1
	for i := range nums {
		idx := nums[i]

		ldx := idx
		rdx := idx
		maxTmp := 1
		for j := i + 1; j < len(nums); j++ {
			val := nums[j]
			if val > rdx {
				rdx = val
				maxTmp++
			}
		}
		for j := i - 1; j > 0; j-- {
			val := nums[j]
			if val < ldx {
				ldx = val
				maxTmp++
			}
		}
		if maxTmp > maxL {
			maxL = maxTmp
		}
	}
	return maxL
}

func TestLengthOfLIS(t *testing.T) {
	a := assert.New(t)
	//a.Equal(4, lengthOfLISWithSimple([]int{10, 9, 2, 5, 3, 7, 101, 18}))
	a.Equal(4, lengthOfLISWithSimple([]int{0, 1, 0, 3, 2, 3}))
}
