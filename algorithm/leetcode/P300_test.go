package leetcode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// https://labuladong.github.io/algo/di-er-zhan-a01c6/dong-tai-g-a223e/dong-tai-g-6ea57/
// https://leetcode.cn/problems/longest-increasing-subsequence/
func lengthOfLISWithSimple(nums []int) int {
	// 定义：dp[i] 表示以 nums[i] 这个数结尾的最长递增子序列的长度
	dp := make([]int, len(nums))
	// base case：dp 数组全都初始化为 1
	for i := range dp {
		dp[i] = 1
	}
	for i := 0; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				dp[i] = max(dp[i], dp[j]+1)

			}
		}
	}

	res := 0
	for i := range dp {
		res = max(res, dp[i])
	}
	return res
}

func lengthOfLISWithBinSearch(nums []int) int {
	top := make([]int, len(nums))
	// 牌组数量初始化为 0
	piles := 0
	for i := 0; i < len(nums); i++ {
		// 要处理的扑克牌
		poker := nums[i]

		/***** 搜索左侧边界的二分查找 *****/
		left, right := 0, piles
		for left < right {
			mid := (left + right) / 2
			if top[mid] > poker {
				right = mid
			} else if top[mid] < poker {
				left = mid + 1
			} else {
				right = mid
			}
		}
		/*********************************/

		// 没找到合适的牌组，新建一组
		if left == piles {
			piles++
		}
		// 把这张牌放到牌组顶
		top[left] = poker
	}
	// 牌组数就是 LIS 长度
	return piles
}

func TestLengthOfLIS(t *testing.T) {
	a := assert.New(t)
	a.Equal(4, lengthOfLISWithSimple([]int{10, 9, 2, 5, 3, 7, 101, 18}))
	a.Equal(4, lengthOfLISWithSimple([]int{0, 1, 0, 3, 2, 3}))

	a.Equal(4, lengthOfLISWithBinSearch([]int{10, 9, 2, 5, 3, 7, 101, 18}))
	a.Equal(4, lengthOfLISWithBinSearch([]int{0, 1, 0, 3, 2, 3}))
}
