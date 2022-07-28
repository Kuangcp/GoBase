package ctool

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func PostJson(url string, value interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return PostJsons(url, jsonBytes)
}

func PostJsons(url string, value []byte) ([]byte, error) {
	reader := bytes.NewReader(value)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return rspBody, nil
}
