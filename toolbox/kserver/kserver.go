package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	_ "embed"

	"github.com/kuangcp/gobase/cuibase"
)

//go:embed up.html
var uploadStaticPage string

//go:embed home.html
var homeStaticPage string

//go:embed favicon.ico
var faviconIco string

var (
	help        bool
	defaultHome bool

	port         int
	buildVersion string
)

type Value interface {
	String() string
	Set(string) error
}

type arrayFlags []string

// Value ...
func (i *arrayFlags) String() string {
	return fmt.Sprint(*i)
}

// Set 方法是flag.Value接口, 设置flag Value的方法.
// 通过多个flag指定的值， 所以我们追加到最终的数组上.
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var folderPair arrayFlags
var pathDirMap = make(map[string]string)

func getInternalIP() string {
	address, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, addr := range address {
		if ipNet, ok := addr.(*net.IPNet); ok &&
			!ipNet.IP.IsLoopback() &&
			ipNet.IP.To4() != nil &&
			strings.HasPrefix(ipNet.IP.String(), "192") {
			return ipNet.IP.String()
		}
	}
	return ""
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var maxMib int64 = 10
	err := r.ParseMultipartForm(maxMib << 20)
	if err != nil {
		log.Println(err)
	}

	for _, headers := range r.MultipartForm.File {
		for _, header := range headers {
			if err := receiveFile(header); err != nil {
				return
			}
		}
	}

	http.Redirect(w, r, "/up", http.StatusMovedPermanently)
}

func receiveFile(header *multipart.FileHeader) error {
	filename := header.Filename
	log.Printf("upload: %s", filename)

	exist := isFileExist(filename)
	if exist {
		filename = fmt.Sprint(time.Now().Nanosecond()) + "-" + filename
	}

	dst, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return err
	}

	open, err := header.Open()
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		err := dst.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	defer func() {
		err := open.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = io.Copy(dst, open)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func echoHandler(_ http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	content := string(body)

	decode, _ := url.QueryUnescape(content)
	if strings.HasPrefix(decode, "content") {
		decode = decode[8:]
	}
	decode = "Content: \n" + decode
	log.Print(decode)
}

func printStartUpLog(port int, internalIP string) {
	innerURL := fmt.Sprintf("http://%v:%v", internalIP, port)
	log.Printf("static file web server has started.\n")

	// sort and print
	var keys []string
	for k := range pathDirMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		printFileAndImgGroup("127.0.0.1", internalIP, k, pathDirMap[k], port)
	}

	log.Printf("%v/up%v  %v/up\n", cuibase.Purple, cuibase.End, innerURL)
	log.Printf("%v/f%v   curl -X POST -H 'Content-Type: multipart/form-data' %v/f -F file=@index.html\n",
		cuibase.Purple, cuibase.End, innerURL)
	log.Printf("%v/e%v   curl %v/e -d 'echo hi'\n", cuibase.Purple, cuibase.End, innerURL)
}

func printFileAndImgGroup(host, internalHost, path, filePath string, port int) {
	if path == "/" {
		path = ""
	}
	local := fmt.Sprintf("http://%v:%v/%v", host, port, path)
	internal := fmt.Sprintf("http://%v:%v/%v", internalHost, port, path)

	log.Printf("%v%-30v %-30v %-35v %v %v", cuibase.Green, local, internal,
		fmt.Sprintf("%v/img", internal), cuibase.End, filePath)
}

func bindPathAndStatic(pattern, binContent string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(binContent))
		if err != nil {
			log.Println(err)
		}
	})
}

func buildImgFunc(parentPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dir, err := os.ReadDir(pathDirMap[parentPath])
		if err != nil {
			w.Write([]byte("error"))
			return
		}
		w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Img</title>
</head>
<body>`))

		sort.Slice(dir, func(i, j int) bool {
			iInfo, _ := dir[i].Info()
			jInfo, _ := dir[j].Info()
			return iInfo.ModTime().After(jInfo.ModTime())
		})

		for _, entry := range dir {
			fileInfo, _ := entry.Info()
			fileInfo.ModTime()
			if entry.IsDir() {
				continue
			}
			w.Write([]byte("<img width=\"300px\" src=\"" + entry.Name() + "\">"))
		}
		w.Write([]byte(`</body></html>`))
	}
}

var info = cuibase.HelpInfo{
	Description:   "Start static file web server on current path",
	Version:       "1.0.9",
	BuildVersion:  buildVersion,
	SingleFlagLen: -2,
	ValueLen:      -6,
	Flags: []cuibase.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help"},
		{Short: "-g", BoolVar: &defaultHome, Comment: "default home page"},
	},
	Options: []cuibase.ParamVO{
		{Short: "-p", Value: "port", Comment: "web server port"},
		{Short: "-d", Value: "folder", Comment: "folder pair. like -d x=y "},
	}}

func init() {
	flag.IntVar(&port, "p", 8989, "")
	flag.Var(&folderPair, "d", "")
}

func main() {
	info.Parse()
	if help {
		info.PrintHelp()
		return
	}

	if port > 65535 || port == 0 {
		log.Fatalf("Please input correct port [1, 65535]. now: %v", port)
	}
	if port < 1024 {
		log.Printf("%vWARN: [1-1024] need run by root user.%v", cuibase.Red, cuibase.End)
	}

	pathDirMap["/"] = "./"
	for _, s := range folderPair {
		if !strings.Contains(s, "=") {
			log.Printf("%vWARN %v is invalid format. must like a=b %v", cuibase.Red, s, cuibase.End)
			continue
		}

		pair := strings.Split(s, "=")
		path := pair[0]
		if path == "f" || path == "img" || path == "h" || path == "up" || path == "e" || path == "d" {
			log.Printf("%vWARN path /%v already bind. %v", cuibase.Red, path, cuibase.End)
			continue
		}
		pathDirMap[path] = pair[1]

		http.Handle("/"+path+"/", http.StripPrefix("/"+path, http.FileServer(http.Dir(pair[1]))))
		http.HandleFunc("/"+path+"/img", buildImgFunc(path))
	}

	printStartUpLog(port, getInternalIP())

	// TODO 优化设计 绑定路由 与 当前目录
	fs := http.FileServer(http.Dir("./"))
	if defaultHome {
		http.Handle("/d/", http.StripPrefix("/d", fs))
		bindPathAndStatic("/", homeStaticPage)
	} else {
		http.Handle("/", http.StripPrefix("/", fs))
		bindPathAndStatic("/h", homeStaticPage)
	}
	http.HandleFunc("/img", buildImgFunc("/"))

	bindPathAndStatic("/favicon.ico", faviconIco)
	bindPathAndStatic("/up", uploadStaticPage)
	http.HandleFunc("/f", uploadHandler)
	http.HandleFunc("/e", echoHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}
