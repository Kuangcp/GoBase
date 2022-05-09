package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/logger"
	"golang.org/x/net/webdav"
)

var (
	port     int
	username string
	pwd      string
	dirPair  cuibase.ArrayFlags
)

func init() {
	flag.IntVar(&port, "p", 33311, "port")
	flag.StringVar(&username, "user", "gin", "username")
	flag.StringVar(&pwd, "pwd", "jiushi", "username")
	flag.Var(&dirPair, "d", "dir eg: x=/path/to")
}

func main() {
	flag.Parse()

	logger.Info("starting: ", port)

	MultiHandle(fmt.Sprintf(":%v", port))
}

func MultiHandle(port string) {
	mux := http.NewServeMux()
	if len(dirPair) == 0 {
		bindSingleHandler(mux)
	} else {
		bindMultiHandler(mux)
	}

	err := http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Println("dav server run error:", err)
	}
}

func bindSingleHandler(mux *http.ServeMux) {
	loHandler := &webdav.Handler{
		FileSystem: webdav.Dir("."),
		LockSystem: webdav.NewMemLS(),
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if !authUser(w, req) {
			return
		}

		loHandler.ServeHTTP(w, req)
	})
}

func bindMultiHandler(mux *http.ServeMux) {
	var l []*webdav.Handler
	for _, s := range dirPair {
		pair := strings.Split(s, "=")
		if len(pair) != 2 {
			continue
		}
		l = append(l, &webdav.Handler{
			Prefix:     "/" + pair[0] + "/",
			FileSystem: webdav.Dir(pair[1]),
			LockSystem: webdav.NewMemLS(),
		})
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if !authUser(w, req) {
			return
		}

		//switch req.Method {
		//case "PUT", "DELETE", "PROPPATCH", "MKCOL", "COPY", "MOVE":
		//	http.Error(w, "WebDAV: Read Only!!!", http.StatusForbidden)
		//	return
		//}

		if l != nil && len(l) != 0 {
			for _, handler := range l {
				if strings.HasPrefix(req.RequestURI, handler.Prefix) {
					handler.ServeHTTP(w, req)
					return
				}
			}
		}

		w.WriteHeader(404)
	})
}

func authUser(w http.ResponseWriter, req *http.Request) bool {
	// 获取用户名/密码
	username, password, ok := req.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	// 验证用户名/密码
	if username != "gin" || password != "jiushi" {
		http.Error(w, "WebDAV: need authorized!", http.StatusUnauthorized)
		return false
	}
	return true
}
