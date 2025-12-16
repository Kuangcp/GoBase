package main

import (
	"bytes"
	"encoding/json"
	"github.com/kuangcp/logger"
	"io"
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

func copyStream(src io.ReadCloser) ([]byte, io.ReadCloser) {
	bodyBt, err := io.ReadAll(src)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	return bodyBt, io.NopCloser(bytes.NewBuffer(bodyBt))
}

func convertList[T any, R any](src []T, mapFun func(T) R, filterFun func(T) bool) []R {
	var result []R
	for _, d := range src {
		if filterFun == nil || filterFun(d) {
			result = append(result, mapFun(d))
		}
	}
	return result
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
