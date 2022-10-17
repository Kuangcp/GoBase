package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var port int

// origin url -> target url
var proxyMap = make(map[*url.URL]*url.URL)

//	/api/a -> /a
//
// registerReplace(map[string]string{"http://host1:port1/api": "http://host2:port2"})
//
//	/api/a -> /api2/a
//
// registerReplace(map[string]string{"http://host1:port1/api": "http://host2:port2/api2"})
func registerReplace(proxy map[string]string) {
	for k, v := range proxy {
		kUrl, err := url.Parse(k)
		if err != nil {
			continue
		}
		vUrl, err := url.Parse(v)
		if err != nil {
			continue
		}
		proxyMap[kUrl] = vUrl
		logger.Info("register:", k, "=>", v)
	}
}

// TODO 当前按Host维度替换，需要实现按路径维度替换
func findTargetReplace(r *http.Request) (*url.URL, *url.URL) {
	for k, v := range proxyMap {
		if k.Host == r.Host {
			return k, v
		}
	}
	return nil, nil
}

func concatIgnoreSlash(left, right string) string {
	aslash := strings.HasSuffix(left, "/")
	bslash := strings.HasPrefix(right, "/")
	switch {
	case aslash && bslash:
		return left + right[1:]
	case !aslash && !bslash:
		return left + "/" + right
	}
	return left + right
}

func handlePath(origin, target *url.URL, path string) string {
	// 原始路径去前缀
	if origin.Path != "" && strings.HasPrefix(path, origin.Path) {
		path = path[len(origin.Path):]
	}
	// 原始路径加前缀
	if target.Path != "" {
		path = concatIgnoreSlash(target.Path, path)
	}
	return path
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	proxyReq := new(http.Request)
	*proxyReq = *r

	log := proxyReplace(proxyReq)

	startMs := time.Now().UnixMilli()
	if q := proxyReq.URL.RawQuery; q != "" {
		proxyReq.URL.RawPath = proxyReq.URL.Path + "?" + q
	} else {
		proxyReq.URL.RawPath = proxyReq.URL.Path
	}
	proxyReq.Proto = "HTTP/1.1"
	proxyReq.ProtoMajor = 1
	proxyReq.ProtoMinor = 1
	proxyReq.Close = false

	transport := http.DefaultTransport
	res, err := transport.RoundTrip(proxyReq)
	endMs := time.Now().UnixMilli()
	if err != nil {
		logger.Error("%vms %v proxy error", endMs-startMs, log, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if log != "" {
		logger.Debug("%vms %v", endMs-startMs, log)
	}

	hdr := w.Header()
	for k, vv := range res.Header {
		for _, v := range vv {
			hdr.Add(k, v)
		}
	}
	for _, c := range res.Cookies() {
		w.Header().Add("Set-Cookie", c.Raw)
	}

	w.WriteHeader(res.StatusCode)
	if res.Body != nil {
		written, err := io.Copy(w, res.Body)
		if err != nil {
			logger.Error(written, err)
		}
	}
}

func proxyReplace(proxyReq *http.Request) string {
	originUrl, targetUrl := findTargetReplace(proxyReq)
	if targetUrl != nil {
		path := handlePath(originUrl, targetUrl, proxyReq.URL.Path)
		log := ""
		if path == proxyReq.URL.Path {
			log = fmt.Sprintf("proxy: %s => %s", proxyReq.Host+proxyReq.URL.Path, targetUrl.Host+" .")
		} else {
			log = fmt.Sprintf("proxy: %s => %s", proxyReq.Host+proxyReq.URL.Path, targetUrl.Host+path)
		}

		proxyReq.Host = targetUrl.Host
		//o.URL.Scheme = targetURL.Scheme
		proxyReq.URL.Host = targetUrl.Host
		proxyReq.URL.Path = path
		//o.URL.RawQuery = targetUrl.RawQuery
		return log
	} else {
		//logger.Debug("direct:", proxyReq.Host)
		return ""
	}
}

func main() {
	flag.IntVar(&port, "p", 1234, "port")
	flag.Parse()

	home, err := ctool.Home()
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.ReadFile(home + "/.dev-proxy.json")
	if err == nil {
		var configMap map[string]string
		err := json.Unmarshal(file, &configMap)
		if err != nil {
			logger.Error(err)
			return
		}
		registerReplace(configMap)
	}

	logger.Info("Start serving on 127.0.0.1:%d", port)
	http.HandleFunc("/", proxyHandler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		logger.Error(err)
	}
	os.Exit(0)
}
