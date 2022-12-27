package leetcode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// 两个指针游走
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	len1 := len(nums1)
	len2 := len(nums2)

	if len1 == 0 {
		if len2%2 == 0 {
			return float64(nums2[len2/2]+nums2[len2/2-1]) / 2.0
		}
		return float64(nums2[len2/2])
	}

	if len2 == 0 {
		if len1%2 == 0 {
			return float64(nums1[len1/2]+nums1[len1/2-1]) / 2.0
		}
		return float64(nums1[len1/2])
	}

	idx1 := 0
	idx2 := 0
	totalLen := len1 + len2
	m := 0

	cache := [2]int{0, 0}
	cacheIdx := 1

	mid := totalLen / 2
	for idx1+idx2 <= mid {
		if idx1 >= len1 {
			m = nums2[idx2]
			cacheIdx++
			cache[cacheIdx%2] = m
			idx2++
			continue
		} else if idx2 >= len2 {
			m = nums1[idx1]
			cacheIdx++
			cache[cacheIdx%2] = m
			idx1++
			continue
		}

		if nums1[idx1] <= nums2[idx2] {
			m = nums1[idx1]
			cacheIdx++
			cache[cacheIdx%2] = m
			idx1++
		} else if nums2[idx2] < nums1[idx1] {
			m = nums2[idx2]
			cacheIdx++
			cache[cacheIdx%2] = m
			idx2++
		}
	}
	if totalLen%2 == 0 {
		return float64(cache[0]+cache[1]) / 2.0
	}
	return float64(m)
}

func TestFindMedianSortedArrays(t *testing.T) {
	assert.Equal(t, float64(6), findMedianSortedArrays([]int{3, 6}, []int{44}))
	assert.Equal(t, float64(5), findMedianSortedArrays([]int{3, 6}, []int{1, 4, 7, 44}))

	assert.Equal(t, float64(2.5), findMedianSortedArrays([]int{1, 2}, []int{3, 4}))
	assert.Equal(t, float64(0), findMedianSortedArrays([]int{0, 0}, []int{0, 0}))
}
