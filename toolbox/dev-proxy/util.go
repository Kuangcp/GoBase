package main

import (
	"bytes"
	"encoding/json"
	"github.com/kuangcp/logger"
	"io"
)

// avoid & => \u0026
func toJSONBuffer(val any) *bytes.Buffer {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(val)
	return buffer
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
