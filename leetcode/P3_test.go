package leetcode

import "testing"

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
