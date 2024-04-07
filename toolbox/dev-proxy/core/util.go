package core

import (
	"bytes"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"net/http"
	"os/exec"
)

func WriteJsonParamError(writer http.ResponseWriter, msg string) {
	resp := ctool.ResultVO[string]{}
	resp.Code = 400
	resp.Msg = msg
	WriteJsonRsp(writer, resp)
}

func WriteJsonError(writer http.ResponseWriter, code int, msg string) {
	resp := ctool.ResultVO[string]{}
	resp.Code = code
	resp.Msg = msg
	WriteJsonRsp(writer, resp)
}

func WriteJsonRsp(writer http.ResponseWriter, val any) {
	writer.Header().Set("Content-Type", "application/json")
	buffer := ctool.ToJSONBuffer(val)
	writer.Write(buffer.Bytes())
}

func Json[T any](serviceFunc func(request *http.Request) T) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		buffer := ctool.ToJSONBuffer(serviceFunc(request))
		writer.Write(buffer.Bytes())
	}
}

func ExecCommand(command string) (string, bool) {
	cmd := exec.Command("/usr/bin/bash", "-c", command)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.Error(err)
		return "", false
	}

	result := out.String()
	return result, true
}

func RspStr(writer http.ResponseWriter, val string) {
	_, err := writer.Write([]byte(val))
	if err != nil {
		logger.Error(err)
	}
}

func Go(act func()) {
	go func() {
		defer func() {
			//捕获抛出的panic
			if err := recover(); err != nil {
				logger.Warn(err)
			}
		}()
		act()
	}()
}
