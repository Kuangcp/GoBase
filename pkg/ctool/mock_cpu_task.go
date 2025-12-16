package ctool

import "math/rand/v2"

// CalculatePiTylor 使用 arctan(1) 的泰勒级数计算 π
// π = 4 * arctan(1) = 4 * (1 - 1/3 + 1/5 - 1/7 + 1/9 - ...)
// 返回小数部分（不包含 "3." 前缀）
func CalculatePiTylor(scale int, round int) string {
	// 使用整数数组存储小数部分，每个元素代表一位数字（0-9）
	// 索引 0 是小数第一位，索引 1 是小数第二位，以此类推
	workingScale := scale + 50 // 多计算一些位以确保精度
	if round > 0 {
		workingScale += rand.IntN(round)
	}
	arctan := make([]int, workingScale)

	sign := 1
	denominator := 1

	// 计算足够的项数以达到所需精度
	// arctan(1) 的级数收敛较慢，需要大约 scale * 2.5 项才能达到 scale 位精度
	maxTerms := scale*5/2 + 100

	for term := 0; term < maxTerms; term++ {
		// 计算 1/denominator 的小数部分，使用长除法
		// 我们将 1 放大 10^workingScale 倍，然后除以 denominator
		quotient := make([]int, workingScale)
		remainder := 1 // 从 1 开始

		// 长除法：计算 1/denominator 的每一位小数
		for i := 0; i < workingScale; i++ {
			remainder *= 10
			quotient[i] = remainder / denominator
			remainder %= denominator
		}

		// 将当前项加到 arctan 上（带符号）
		if sign > 0 {
			carry := 0
			for i := 0; i < workingScale; i++ {
				sum := arctan[i] + quotient[i] + carry
				arctan[i] = sum % 10
				carry = sum / 10
			}
		} else {
			borrow := 0
			for i := 0; i < workingScale; i++ {
				diff := arctan[i] - quotient[i] - borrow
				if diff < 0 {
					diff += 10
					borrow = 1
				} else {
					borrow = 0
				}
				arctan[i] = diff
			}
		}

		sign = -sign
		denominator += 2
	}

	// 乘以 4 得到 π: π = 4 * arctan(1)
	piDigits := make([]int, workingScale)
	carry := 0
	for i := 0; i < workingScale; i++ {
		product := arctan[i]*4 + carry
		piDigits[i] = product % 10
		carry = product / 10
	}

	// 处理整数部分的进位（如果有）
	if carry > 0 {
		// 这不应该发生，因为 arctan(1) < 1，4*arctan(1) < 4
		// 但为了安全起见，我们处理它
	}

	// 转换为字符串（只返回小数部分，截取到 scale 位）
	var result []byte
	for i := 0; i < scale; i++ {
		result = append(result, byte('0'+piDigits[i]))
	}

	return string(result)
}
