package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// origin url -> target url
var proxyMap = make(map[*url.URL]*url.URL)

// "http://192.168.16.91:32149/(.*)", "http://127.0.0.1:19011/$1"
func registerProxy(proxy map[string]string) {
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
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func findTarget(r *http.Request) *url.URL {
	for k, v := range proxyMap {
		if k.Host == r.Host {
			return v
		}
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	proxyReq := new(http.Request)
	*proxyReq = *r

	targetUrl := findTarget(proxyReq)
	if targetUrl != nil {
		fmt.Printf("proxy: %s -> %s\n", proxyReq.Host, targetUrl.Host)
		proxyReq.Host = targetUrl.Host
		//o.URL.Scheme = targetURL.Scheme
		proxyReq.URL.Host = targetUrl.Host
		proxyReq.URL.Path = singleJoiningSlash(targetUrl.Path, proxyReq.URL.Path)
		//o.URL.RawQuery = targetUrl.RawQuery
	} else {
		fmt.Println("direct:", proxyReq.Host)
	}

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
	if err != nil {
		log.Printf("http: proxy error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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
		io.Copy(w, res.Body)
	}
}

func main() {
	//registerProxy(map[string]string{"http://192.168.16.91:32149/(.*)": "http://127.0.0.1:19011/$1"})
	registerProxy(map[string]string{"http://192.168.16.91:32149": "http://127.0.0.1:19011"})

	log.Println("Start serving on port 1234")

	http.HandleFunc("/", handler)
	http.ListenAndServe(":1234", nil)
	os.Exit(0)
}
