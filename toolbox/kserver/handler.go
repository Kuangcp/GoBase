package main

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/cuibase"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var imgSuffixSet = cuibase.NewSet(".jpg", ".jpeg", ".png", ".svg", ".webp", ".bmp", ".gif", ".ico")
var videoSuffixSet = cuibase.NewSet(".mp4")

type MediaParam struct {
	rawSize bool
	count   int
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
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

func buildVideoFunc(parentPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		param := resolveMediaParam(r)

		dir, err := os.ReadDir(pathDirMap[parentPath])
		if err != nil {
			w.Write([]byte("read dir " + pathDirMap[parentPath] + " error"))
			return
		}

		w.Write([]byte(`<!DOCTYPE html><html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>` + pathDirMap[parentPath] + `</title>
			</head>
			<body>`))

		w.Write([]byte(buildMediaList(dir, param.count, videoSuffixSet, videoTag) + `</body></html>`))
	}
}

func buildImgFunc(parentPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		param := resolveMediaParam(r)
		dir, err := os.ReadDir(pathDirMap[parentPath])
		if err != nil {
			w.Write([]byte("read dir " + pathDirMap[parentPath] + " error"))
			return
		}

		w.Write([]byte(`<!DOCTYPE html><html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>` + pathDirMap[parentPath] + `</title>
			<style>`))

		if !param.rawSize {
			w.Write([]byte(`
			img {
				width: 210px;
				max-height: 100px;
			}`))
		}

		w.Write([]byte(`</style></head><body>` +
			buildMediaList(dir, param.count, imgSuffixSet, imgTag) +
			`</body></html>`))
	}
}

func resolveMediaParam(r *http.Request) MediaParam {
	query := r.URL.Query()
	rawSize := query.Get("raw")
	count := query.Get("count")

	countInt := 5
	countTmp, err := strconv.Atoi(count)
	if err == nil && countTmp > 0 {
		countInt = countTmp
	}
	return MediaParam{rawSize: rawSize != "", count: countInt}
}

func sortByModTime(dir []os.DirEntry) {
	sort.Slice(dir, func(i, j int) bool {
		iInfo, _ := dir[i].Info()
		jInfo, _ := dir[j].Info()
		return iInfo.ModTime().After(jInfo.ModTime())
	})
}

func matchSuffix(set *cuibase.Set, fileName string) bool {
	idx := strings.LastIndex(fileName, ".")
	if idx == -1 {
		return true
	}

	suffixType := fileName[idx:]
	return set.Contains(suffixType)
}

func buildMediaList(dir []os.DirEntry, count int, set *cuibase.Set, tagFunc func(string) string) string {
	sortByModTime(dir)
	mediaBody := ""
	mediaCount := 0

	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		if mediaCount == count {
			break
		}

		fileName := entry.Name()
		if matchSuffix(set, fileName) {
			mediaBody += tagFunc(fileName)
			mediaCount++
			continue
		}
	}

	if mediaCount == 0 {
		return "<h2>No Media</h2>"
	}
	return mediaBody
}

func videoTag(fileName string) string {
	return "<video src=\"" + url.PathEscape(fileName) + "\" controls=\"controls\"></video>"
}

func imgTag(fileName string) string {
	return "<img  src=\"" + url.PathEscape(fileName) + "\" alt=\"" + fileName + "\">"
}
