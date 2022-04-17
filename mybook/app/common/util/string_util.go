package util

import (
	"encoding/json"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"strconv"
	"strings"
)

// 中文字符占用 1.5 英文字符宽度
func BuildCHCharFormat(expectLen int, str string) string {
	return "%" + strconv.Itoa(expectLen-len(str)/3*2) + "s"
}

func Json(data interface{}) string {
	bytes, e := json.Marshal(data)
	if e != nil {
		return "ERROR"
	}
	return string(bytes)
}

func ParseMultiPrice(amount string) ghelp.ResultVO {
	amount = strings.Replace(amount, "，", ",", -1)
	amountList := strings.Split(amount, ",")
	var totalAmount = 0
	for _, one := range amountList {
		parseResult := ParsePrice(one)
		if parseResult.IsFailed() {
			return parseResult
		}
		totalAmount += parseResult.Data.(int)
	}
	return ghelp.SuccessWith(totalAmount)
}

func ParsePrice(amount string) ghelp.ResultVO {
	floatAmount, e := strconv.ParseFloat(amount, 64)
	if e != nil || floatAmount <= 0 {
		return ghelp.FailedWithMsg("金额错误" + amount)
	}

	nums := strings.Split(amount, ".")
	if len(nums) == 1 {
		return ghelp.SuccessWith(int(floatAmount) * 100)
	}

	v := nums[1]
	pInt, _ := strconv.Atoi(nums[0])
	vInt, _ := strconv.Atoi(v)

	vLen := len(v)
	if vLen > 2 {
		return ghelp.FailedWithMsg("金额仅保留两位小数")
	}

	// vLen only 1 Or 2
	if vLen == 1 {
		vInt *= 10
	}
	return ghelp.SuccessWith(pInt*100 + vInt)
}
