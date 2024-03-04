package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/kuangcp/gobase/pkg/ctk"
	"github.com/kuangcp/logger"
	"golang.org/x/net/webdav"
)

var (
	port    int
	user    string
	pwd     string
	dirPair ctk.ArrayFlags
)

func init() {
	flag.IntVar(&port, "p", 33311, "port")
	flag.StringVar(&user, "user", "gin", "username")
	flag.StringVar(&pwd, "pwd", "jiushi", "password")
	flag.Var(&dirPair, "d", "webdav dir(default current dir). eg: x=/path/to")
}

func main() {
	flag.Parse()

	logger.Info("starting: ", port)

	MultiHandle(fmt.Sprintf(":%v", port))
}

func MultiHandle(port string) {
	mux := http.NewServeMux()
	var list []*webdav.Handler

	if len(dirPair) == 0 {
		list = append(list, &webdav.Handler{
			Prefix:     "/",
			FileSystem: webdav.Dir("."),
			LockSystem: webdav.NewMemLS(),
		})
	} else {
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
	}

	bindMultiHandler(mux, list)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Println("dav server run error:", err)
	}
}

func bindMultiHandler(mux *http.ServeMux, l []*webdav.Handler) {
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
		logger.Warn("no auth")
		return false
	}

	// 验证用户名/密码
	if username != user || password != pwd {
		http.Error(w, "WebDAV: need authorized!", http.StatusUnauthorized)
		logger.Warn("no match user:", username, password)
		return false
	}
	return true
}
