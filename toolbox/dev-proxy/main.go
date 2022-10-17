package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var targetURL *url.URL

// origin url -> target url
var proxyMap map[*url.URL]*url.URL

// "http://192.168.16.91:32149/(.*)", "http://127.0.0.1:19011/$1"
func registerProxy(proxy map[string]string) {
	targetServer := "http://localhost:19011"
	tmpUrl, err := url.Parse(targetServer)
	if err != nil {
		log.Println("Bad target URL")
		return
	}
	targetURL = tmpUrl

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

func handler(w http.ResponseWriter, r *http.Request) {
	// copy object
	o := new(http.Request)
	*o = *r

	o.Host = targetURL.Host
	o.URL.Scheme = targetURL.Scheme
	o.URL.Host = targetURL.Host
	o.URL.Path = singleJoiningSlash(targetURL.Path, o.URL.Path)
	if q := o.URL.RawQuery; q != "" {
		o.URL.RawPath = o.URL.Path + "?" + q
	} else {
		o.URL.RawPath = o.URL.Path
	}

	o.URL.RawQuery = targetURL.RawQuery

	o.Proto = "HTTP/1.1"
	o.ProtoMajor = 1
	o.ProtoMinor = 1
	o.Close = false

	transport := http.DefaultTransport

	res, err := transport.RoundTrip(o)
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
	registerProxy(map[string]string{"http://192.168.16.91:32149/(.*)": "http://127.0.0.1:19011/$1"})

	log.Println("Start serving on port 1234")

	http.HandleFunc("/", handler)
	http.ListenAndServe(":1234", nil)
	os.Exit(0)
}
