package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/kuangcp/logger"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/kuangcp/gobase/pkg/ctool"
)

var (
	//go:embed up.html
	uploadHtml string
	//go:embed home.html
	homeHtml string
	//go:embed favicon.ico
	faviconIco string
)

var (
	port       int
	folderPair ctool.ArrayFlags
	syncMode   bool

	help         bool
	defaultHome  bool
	buildVersion string
	internalIP   string

	homePath      = "/h"
	imgFilePath   = "/g"
	videoFilePath = "/v"
	fileSys       = http.FileServer(http.Dir("./"))
	pathDirMap    = make(map[string]string)
	usedPath      = ctool.NewSet[string]("f", "g", "h", "up", "e", "d")
)

var info = ctool.HelpInfo{
	Description:   "Start static file web server on current path",
	Version:       "1.1.0",
	BuildVersion:  buildVersion,
	SingleFlagLen: -2,
	ValueLen:      -6,
	Flags: []ctool.ParamVO{
		{Short: "-h", BoolVar: &help, Comment: "help"},
		{Short: "-g", BoolVar: &defaultHome, Comment: "default home page"},
		{Short: "-s", BoolVar: &syncMode, Comment: "sync file or msg mode"},
	},
	Options: []ctool.ParamVO{
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
		log.Printf("%vWARN: [1-1024] need run by root user.%v", ctool.Red, ctool.End)
	}
	internalIP = ctool.GetInternalIP()

	registerAllFolder()
	printStartUpLog()

	http.Handle("/", appendLink("/", http.StripPrefix("/", fileSys)))
	bindPathAndStatic("/h", homeHtml)

	bindPathAndStatic("/favicon.ico", faviconIco)

	bindPathAndStatic("/up", uploadHtml)

	http.HandleFunc("/f", uploadReadHandler)
	http.HandleFunc("/e", echoHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func printStartUpLog() {
	innerURL := fmt.Sprintf("http://%v:%v", internalIP, port)
	log.Printf("%v/up%v  %v/up\n", ctool.Purple, ctool.End, innerURL)
	log.Printf("%v/f%v   curl -X POST -H 'Content-Type: multipart/form-data' %v/f -F file=@index.html\n",
		ctool.Purple, ctool.End, innerURL)
	log.Printf("%v/e%v   curl %v/e -d 'echo hi'\n", ctool.Purple, ctool.End, innerURL)
	log.Printf("%v/h%v   home: %v\n", ctool.Purple, ctool.End, fmt.Sprintf("http://%v:%v%v", internalIP, port, homePath))

	// sort and print
	var keys []string
	for k := range pathDirMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		printFileAndImgGroup("127.0.0.1", k, pathDirMap[k])
	}
}

func printFileAndImgGroup(host, path, filePath string) {
	local := fmt.Sprintf("http://%v:%v/%v", host, port, path)
	internal := fmt.Sprintf("http://%v:%v/%v", internalIP, port, path)
	internal = strings.TrimRight(internal, "/")
	local = strings.TrimRight(local, "/")

	lineBuf := fmt.Sprintf("%v%-27v", ctool.Green, local)
	lineBuf += fmt.Sprintf("%-29v", fmt.Sprintf("%v", internal+imgFilePath))

	logger.Info("%v %v %v", lineBuf, ctool.End, filePath)
}

func bindPathAndStatic(pattern, binContent string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(binContent))
		if err != nil {
			logger.Error(err)
		}
	})
}

func registerAllFolder() {
	pathDirMap["/"] = "./"

	// current dir
	http.HandleFunc(imgFilePath, buildImgFunc("/"))
	http.HandleFunc(videoFilePath, buildVideoFunc("/"))

	// new pair dir from param
	for _, s := range folderPair {
		if !strings.Contains(s, "=") {
			logger.Info("%v is invalid format. must like a=b", s)
			continue
		}

		pair := strings.Split(s, "=")
		path := pair[0]
		if usedPath.Contains(path) {
			logger.Info("path /%v already bind.", path)
			continue
		}
		pathDirMap[path] = pair[1]

		// 动态生成静态文件页面
		http.Handle("/"+path+"/", appendLink("/"+path+"/",
			http.StripPrefix("/"+path, http.FileServer(http.Dir(pair[1])))))

		// 动态生成图片和视频 页面
		http.HandleFunc("/"+path+imgFilePath, buildImgFunc(path))
		http.HandleFunc("/"+path+videoFilePath, buildVideoFunc(path))
	}
}
