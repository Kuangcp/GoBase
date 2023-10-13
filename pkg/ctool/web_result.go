package ctool

import (
	"encoding/json"
	"fmt"
)

const (
	SUCCESS = 0
	FAILED  = 1
)

type (
	PageVO[T any] struct {
		Total int `json:"total"`
		Page  int `json:"page"`
		Data  []T `json:"data"`
	}

	ResultVO[T any] struct {
		Data T      `json:"data"`
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
)

func SuccessWith[T any](data T) ResultVO[T] {
	return ResultVO[T]{Data: data, Code: SUCCESS}
}

func Success[T any]() ResultVO[T] {
	return ResultVO[T]{Code: SUCCESS}
}

func Failed[T any]() ResultVO[T] {
	return ResultVO[T]{Code: FAILED}
}
func Fail[T any](code int, msg string) ResultVO[T] {
	return ResultVO[T]{Msg: msg, Code: code}
}

func FailedWithMsg[T any](msg string) ResultVO[T] {
	return ResultVO[T]{Msg: msg, Code: FAILED}
}

func (result ResultVO[any]) IsSuccess() bool {
	return result.Code == SUCCESS
}
func (result ResultVO[any]) IsFailed() bool {
	return !result.IsSuccess()
}

func (result ResultVO[any]) JSON() []byte {
	marshal, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return marshal
}

func (result ResultVO[any]) JSONStr() string {
	return string(result.JSON())
}
