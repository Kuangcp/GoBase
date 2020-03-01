package vo

import "github.com/kuangcp/gobase/mybook/constant"

type (
	ResultVO struct {
		Data    interface{}
		Code    int8
		Success bool
		Msg     string
	}
)

var success = ResultVO{Code: constant.SUCCESS, Success: true}
var failed = ResultVO{Code: constant.FAILED, Success: false}

func SuccessWith(data interface{}) ResultVO {
	return ResultVO{Data: data, Code: constant.SUCCESS, Success: true}
}

func Success() ResultVO {
	return success
}
func Failed() ResultVO {
	return failed
}

func FailedWithMsg(msg string) ResultVO {
	return ResultVO{Msg: msg, Code: constant.FAILED, Success: false}
}

func (result ResultVO) IsSuccess() bool {
	return result.Code == constant.SUCCESS
}
func (result ResultVO) IsFailed() bool {
	return !result.IsSuccess()
}
