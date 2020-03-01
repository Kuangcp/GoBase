package util

import (
	"encoding/json"
)

func Copy(data interface{}, target interface{}) interface{}{
	jsonStr, e := json.Marshal(data)
	if e != nil {
		return nil
	}
	e = json.Unmarshal(jsonStr, target)
	if e != nil {
		return nil
	}
	return target
}
