package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
)

var imgSuffixSet = ctool.NewSet(".jpg", ".jpeg", ".png", ".svg", ".webp", ".bmp", ".gif", ".ico")
var videoSuffixSet = ctool.NewSet(".mp4")

type MediaParam struct {
	rawSize bool
	count   int
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func uploadReadHandler(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		logger.Error(err)
		w.Write([]byte("error"))
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		filename := part.FileName()
		fmt.Printf("FileName: %s, FormType: %s\n", filename, part.FormName())
		if filename == "" { // this is FormData
			data, _ := io.ReadAll(part)
			fmt.Printf("FormData=[%s]\n", string(data))
		} else {
			// This is FileData
			exist := isFileExist(filename)
			if exist {
				filename = time.Now().Format(ctool.HH_MM_SS_MS) + "-" + filename
			}

			dst, _ := os.Create("./" + filename)
			defer dst.Close()
			_, err := io.Copy(dst, part)
			if err != nil {
				logger.Error(err)
				w.Write([]byte("copy file error"))
				return
			}
		}
	}
	http.Redirect(w, r, "/up", http.StatusMovedPermanently)
}

func appendLink(rootPath string, origin http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if syncMode && request.URL.Path == rootPath {
			writer.Write([]byte(`<html>
<button onclick='location.href=("/h")'>首页</button>
<button onclick='location.href=("/up")'>上传</button>
<br/>
`))
			origin.ServeHTTP(writer, request)
			writer.Write([]byte("</html>"))
		} else {
			origin.ServeHTTP(writer, request)
		}
	}
}

func echoHandler(_ http.ResponseWriter, request *http.Request) {
	body, _ := io.ReadAll(request.Body)
	content := string(body)

	decode, _ := url.QueryUnescape(content)
	if strings.HasPrefix(decode, "content") {
		decode = decode[8:]
	}
	decode = "💠Content: \n" + decode
	logger.Info(decode)
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

	countInt := 4
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

func matchSuffix(set *ctool.Set[string], fileName string) bool {
	idx := strings.LastIndex(fileName, ".")
	if idx == -1 {
		return true
	}

	suffixType := fileName[idx:]
	return set.Contains(suffixType)
}

func buildMediaList(dir []os.DirEntry, count int, set *ctool.Set[string], tagFunc func(string) string) string {
	sortByModTime(dir)
	mediaBody := ""
	mediaCount := 0

	for _, entry := range dir {
		fileName := entry.Name()
		if entry.IsDir() {
			mediaBody += dirTag(fileName)
			continue
		}
		if mediaCount == count {
			break
		}

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

func dirTag(fileName string) string {
	return "<a href=" + url.PathEscape(fileName) + ">" + fileName + "</a><br/>"
}

func videoTag(fileName string) string {
	return "<video src=\"" + url.PathEscape(fileName) + "\" controls=\"controls\"></video>"
}

func imgTag(fileName string) string {
	return "<img  src=\"" + url.PathEscape(fileName) + "\" alt=\"" + fileName + "\">"
}
