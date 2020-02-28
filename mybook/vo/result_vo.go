package vo

import "github.com/kuangcp/gobase/myth-bookkeeping/constant"

type (
	ResultVO struct {
		Data interface{}
		Code int8
		Msg  string
	}
)

var success = ResultVO{Code: constant.SUCCESS}
var failed = ResultVO{Code: constant.FAILED}

func SuccessWith(data interface{}) ResultVO {
	return ResultVO{Data: data, Code: constant.SUCCESS}
}

func Success() ResultVO {
	return success
}
func Failed() ResultVO {
	return failed
}

func FailedWithMsg(msg string) ResultVO {
	return ResultVO{Msg: msg, Code: constant.FAILED}
}

func (this ResultVO) IsSuccess() bool {
	return this.Code == constant.SUCCESS
}
func (this ResultVO) IsFailed() bool {
	return !this.IsSuccess()
}
