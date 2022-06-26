package weblevel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kuangcp/logger"
	"io/ioutil"
	"net/http"
	"net/url"
)

type (
	WebClient struct {
		host string
		port int
		api  *http.Client
	}
)

func NewClient(host string, port int) *WebClient {
	if port <= 0 {
		port = Port
	}
	return &WebClient{host: host, port: port, api: &http.Client{}}
}

func (w *WebClient) buildPath(path string) string {
	return fmt.Sprintf("http://%v:%v%v", w.host, w.port, path)
}

func (w *WebClient) Del(key string) error {
	fmtKey := url.QueryEscape(key)
	resp, err := w.api.Get(w.buildPath(PathDel) + "?key=" + fmtKey)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (w *WebClient) Get(key string) (string, error) {
	fmtKey := url.QueryEscape(key)
	resp, err := w.api.Get(w.buildPath(PathGet) + "?key=" + fmtKey)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBt), nil
}

func (w *WebClient) Set(key, val string) {
	m := make(map[string]string)
	m[key] = val

	w.Sets(m)
}

func (w *WebClient) Sets(kv map[string]string) {
	if len(kv) == 0 {
		return
	}

	var val []*ValKV
	for k, v := range kv {
		val = append(val, &ValKV{Key: k, Val: v})
	}

	w.sendJsonPost(val, w.buildPath(PathSets))
	//logger.Info(string(rsp))
}

func (w *WebClient) sendJsonPost(value interface{}, url string) []byte {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		logger.Error(err)
		return nil
	}
	//fmt.Println(url, string(jsonBytes))
	reader := bytes.NewReader(jsonBytes)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		logger.Error(err)
		return nil
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := w.api.Do(request)
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer resp.Body.Close()

	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return rspBody
}

func (w *WebClient) Stats() {
	resp, err := http.Get(w.buildPath(PathStat))
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()
}

func (w *WebClient) PrefixSearch(prefix string) map[string]string {
	resp, err := http.Get(w.buildPath(PathSearch) + "?prefix=" + prefix)
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader(all))
	var c map[string]string
	err = decoder.Decode(&c)
	if err != nil {
		return nil
	}
	return c

}
