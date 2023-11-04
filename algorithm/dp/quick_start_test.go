package dp

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

// https://labuladong.github.io/algo/di-er-zhan-a01c6/dong-tai-g-a223e/dong-tai-g-1e688/

func Fib(n int) int {
	if n == 1 || n == 2 {
		return 1
	}
	return Fib(n-1) + Fib(n-2)
}

// 缓存计算结果，降低重复运算
func FibWithMemo(n int) int {
	memo := make([]int, n+1)
	return fibWithMemo(memo, n)
}

func fibWithMemo(memo []int, n int) int {
	if n == 0 || n == 1 {
		return n
	}
	// 已经计算过，不用再计算了
	if memo[n] != 0 {
		return memo[n]
	}
	memo[n] = fibWithMemo(memo, n-1) + fibWithMemo(memo, n-2)
	return memo[n]
}

func FibWithDp(n int) int {
	if n == 0 {
		return 0
	}
	// dp table
	dp := make([]int, n+1)
	dp[0], dp[1] = 0, 1
	// 状态转移
	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}

	return dp[n]
}

func FibWithDpOpt(n int) int {
	if n == 0 {
		return 0
	}
	// 分别代表 dp[i - 1] 和 dp[i - 2]
	dpI1, dpI2 := 1, 0
	for i := 2; i <= n; i++ {
		// dp[i] = dp[i - 1] + dp[i - 2];
		tmp := dpI1 + dpI2
		// 滚动更新
		dpI2 = dpI1
		dpI1 = tmp
	}
	return dpI1
}

func TestFib(t *testing.T) {
	a := assert.New(t)
	a.Equal(3, Fib(4))
}

func TestFibWithMemos(t *testing.T) {
	a := assert.New(t)
	a.Equal(8, FibWithMemo(6))
}

func TestFibWithDp(t *testing.T) {
	a := assert.New(t)
	a.Equal(8, FibWithDp(6))
}

func TestFibWithDpOpt(t *testing.T) {
	a := assert.New(t)
	a.Equal(8, FibWithDpOpt(6))
}

// 定义：要凑出金额n，至少要dp(coins,n)个硬币
func coinChange(coins []int, amount int) int {
	//base case
	if amount == 0 {
		return 0
	}
	if amount < 0 {
		return -1
	}
	res := math.MaxInt
	for _, coin := range coins {
		//计算子问题的结果
		subProblem := coinChange(coins, amount-coin)
		//子问题无解则跳过
		if subProblem == -1 {
			continue
		}
		//在子问题中选择最优解，然后加一
		res = minInt(res, subProblem+1)
	}
	if res == math.MaxInt {
		return -1
	}
	return res
}
func coinChangeWithArray(coins []int, amount int) int {
	dp := make([]int, amount+1)
	// 数组大小为 amount + 1，初始值也为 amount + 1
	for i := 0; i < len(dp); i++ {
		dp[i] = amount + 1
	}

	// base case
	dp[0] = 0
	// 外层 for 循环在遍历所有状态的所有取值
	for i := 0; i < len(dp); i++ {
		// 内层 for 循环在求所有选择的最小值
		for _, coin := range coins {
			// 子问题无解，跳过
			if i-coin < 0 {
				continue
			}
			dp[i] = minInt(dp[i], dp[i-coin]+1)

		}
	}

	if dp[amount] == amount+1 {
		return -1
	}
	return dp[amount]
}
func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func TestCoinChange(t *testing.T) {
	a := assert.New(t)
	a.Equal(3, coinChange([]int{1, 2, 5}, 11))
	a.Equal(3, coinChangeWithArray([]int{1, 2, 5}, 11))
}
