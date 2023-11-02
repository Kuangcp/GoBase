package graph

import (
	"fmt"
	"testing"
)

func equationsPossible(equations []string) bool {
	// 26 个英文字母
	uf := NewUF(26)
	// 先让相等的字母形成连通分量
	for _, eq := range equations {
		if eq[1] == '=' {
			x := eq[0]
			y := eq[3]
			uf.Union(int(x-'a'), int(y-'a'))
		}
	}
	// 检查不等关系是否打破相等关系的连通性
	for _, eq := range equations {
		if eq[1] == '!' {
			x := eq[0]
			y := eq[3]
			// 如果相等关系成立，就是逻辑冲突
			if uf.Find(int(x-'a')) == uf.Find(int(y-'a')) {
				return false
			}
		}
	}
	return true
}

// https://leetcode.cn/problems/satisfiability-of-equality-equations/
func TestSameExpression(t *testing.T) {
	fmt.Println(equationsPossible([]string{"a==b", "b==c", "c==a"}))
	fmt.Println(equationsPossible([]string{"a==b", "b==c", "c!=a"}))
}
