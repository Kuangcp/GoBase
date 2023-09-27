package core

import (
	"bytes"
	"encoding/json"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"net/http"
	"os/exec"
)

// avoid & => \u0026
func toJSONBuffer(val any) *bytes.Buffer {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(val)
	return buffer
}

func writeJsonParamError(writer http.ResponseWriter, msg string) {
	resp := ctool.ResultVO[string]{}
	resp.Code = 400
	resp.Msg = msg
	writeJsonRsp(writer, resp)
}

func writeJsonError(writer http.ResponseWriter, code int, msg string) {
	resp := ctool.ResultVO[string]{}
	resp.Code = code
	resp.Msg = msg
	writeJsonRsp(writer, resp)
}

func writeJsonRsp(writer http.ResponseWriter, val any) {
	writer.Header().Set("Content-Type", "application/json")
	buffer := toJSONBuffer(val)
	writer.Write(buffer.Bytes())
}

func JSONFunc[T any](serviceFunc func(request *http.Request) T) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		buffer := toJSONBuffer(serviceFunc(request))
		writer.Write(buffer.Bytes())
	}
}

func copyObj[T any, R any](src T) *R {
	jsonStr := toJSONBuffer(src).String()
	var r R
	rObj := &r
	err := json.Unmarshal([]byte(jsonStr), rObj)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return rObj
}

func execCommand(command string) (string, bool) {
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
