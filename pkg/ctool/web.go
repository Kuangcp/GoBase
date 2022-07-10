package ctool

const (
	SUCCESS = 0
	FAILED  = 1
)

type (
	ResultVO[T any] struct {
		Data T      `json:"data"`
		Code int8   `json:"code"`
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

func FailedWithMsg[T any](msg string) ResultVO[T] {
	return ResultVO[T]{Msg: msg, Code: FAILED}
}

func (result ResultVO[any]) IsSuccess() bool {
	return result.Code == SUCCESS
}
func (result ResultVO[any]) IsFailed() bool {
	return !result.IsSuccess()
}
