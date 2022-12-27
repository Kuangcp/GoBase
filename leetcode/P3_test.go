package leetcode

import (
	"strings"
	"testing"
)

/*
 * @lc app=leetcode.cn id=3 lang=golang
 *
 * [3] 无重复字符的最长子串
 */
func lengthOfLongestSubstring(s string) int {
	// 哈希集合，记录每个字符是否出现过
	m := map[byte]int{}
	n := len(s)
	// 右指针，初始值为 -1，相当于我们在字符串的左边界的左侧，还没有开始移动
	rk, ans := -1, 0
	for i := 0; i < n; i++ {
		if i != 0 {
			// 左指针向右移动一格，移除一个字符
			delete(m, s[i-1])
		}
		for rk+1 < n && m[s[rk+1]] == 0 {
			// 不断地移动右指针
			m[s[rk+1]]++
			rk++
		}
		// 第 i 到 rk 个字符是一个极长的无重复字符子串
		ans = max(ans, rk-i+1)
	}
	return ans
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// 左右指针
func lengthOfLongestSubstring2(s string) int {
	length, left, right := 0, 0, 0
	current := s[left:right]

	for ; right < len(s); right++ {
		if index := strings.IndexByte(current, s[right]); index != -1 {
			left += index + 1
		}
		current = s[left : right+1]
		if len(current) > length {
			length = len(current)
		}
	}

	return length
}

func Test_lengthOfLongestSubstring(t *testing.T) {
	testFunc(t, lengthOfLongestSubstring)
}

func Test_lengthOfLongestSubstring2(t *testing.T) {
	testFunc(t, lengthOfLongestSubstring2)
}

func testFunc(t *testing.T, method func(string) int) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{args: args{s: "longlongago"}, want: 5},
		{args: args{s: "uiu"}, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := method(tt.args.s); got != tt.want {
				t.Errorf("lengthOfLongestSubstring_2() = %v, want %v", got, tt.want)
			}
		})
	}
}
