package ctool

import (
	"bytes"
	"encoding/json"
	"log"
)

// avoid & => \u0026
func ToJSONBuffer(val any) *bytes.Buffer {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(val)
	return buffer
}

func CopyObj[T any, R any](src T) *R {
	jsonStr := ToJSONBuffer(src).String()
	var r R
	rObj := &r
	err := json.Unmarshal([]byte(jsonStr), rObj)
	if err != nil {
		log.Println(err)
		return nil
	}
	return rObj
}
