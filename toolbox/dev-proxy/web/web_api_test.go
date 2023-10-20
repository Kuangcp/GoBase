package web

import (
	"fmt"
	"github.com/kuangcp/logger"
	"net/url"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {

	var s = []string{"ss", "a", "a", "a", "a", "a", "b"}
	array := splitArray(s, 5)
	logger.Info(array)
}

func TestUrl(t *testing.T) {
	parse, err := url.Parse("https://yunshu.sinohealth.com/xxl-job-admin/joblog/getJobsByGroup?jobGroup=-1")
	if err != nil {
		logger.Error(err)
	}
	path := parse.Path
	parts := strings.Split(path, "/")
	depth := 5
	depth += 1
	if len(parts) < depth {
		depth = len(parts)
	}
	newPath := strings.Join(parts[:depth], "/")
	fmt.Println(parse.Scheme + "://" + parse.Host + newPath)
}

func TestFindDir(t *testing.T) {
	s := "/home/zk/Apps/dev-proxy-cd/dev-proxy.json"
	dirs := strings.Split(s, "/")
	d := strings.Join(dirs[:len(dirs)-1], "/") + "/bak"
	fmt.Println(d)
}
