package main

import (
	"fmt"
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

type ResourceParam struct {
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
		param := resolveImgParam(r)

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

		w.Write([]byte(buildVideoList(dir, param.count)))
		w.Write([]byte(`</body></html>`))
	}
}

func buildImgFunc(parentPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		param := resolveImgParam(r)
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

		w.Write([]byte(`</style></head><body>`))
		imgBody := buildImgListArea(dir, param.count)
		w.Write([]byte(imgBody + `</body></html>`))
	}
}

func resolveImgParam(r *http.Request) ResourceParam {
	query := r.URL.Query()
	rawSize := query.Get("raw")
	count := query.Get("count")

	countInt := 5
	countTmp, err := strconv.Atoi(count)
	if err == nil && countTmp > 0 {
		countInt = countTmp
	}
	return ResourceParam{rawSize: rawSize != "", count: countInt}
}

func buildVideoList(dir []os.DirEntry, count int) string {
	sort.Slice(dir, func(i, j int) bool {
		iInfo, _ := dir[i].Info()
		jInfo, _ := dir[j].Info()
		return iInfo.ModTime().After(jInfo.ModTime())
	})

	videoBody := ""
	videoCount := 0
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		if videoCount == count {
			break
		}

		fileName := entry.Name()
		idx := strings.LastIndex(fileName, ".")
		if idx == -1 {
			videoBody += buildVideoTag(fileName)
			videoCount++
			continue
		}

		suffixType := fileName[idx:]
		if suffixType == ".mp4" {
			videoBody += buildVideoTag(fileName)
			videoCount++
		}
	}

	if videoBody == "" {
		return "<h2>No Video</h2>"
	}
	return videoBody
}

func buildImgListArea(dir []os.DirEntry, countInt int) string {
	sort.Slice(dir, func(i, j int) bool {
		iInfo, _ := dir[i].Info()
		jInfo, _ := dir[j].Info()
		return iInfo.ModTime().After(jInfo.ModTime())
	})

	imgCount := 0
	imgBodyH5 := ""
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		if imgCount == countInt {
			break
		}

		fileName := entry.Name()
		idx := strings.LastIndex(fileName, ".")
		if idx == -1 {
			imgBodyH5 += buildImgTag(fileName)
			imgCount++
			continue
		}

		suffixType := fileName[idx:]
		if suffixType == ".jpg" || suffixType == ".png" || suffixType == ".svg" || suffixType == ".webp" ||
			suffixType == ".bmp" || suffixType == ".gif" || suffixType == ".ico" {
			imgBodyH5 += buildImgTag(fileName)
			imgCount++
		}
	}
	if imgBodyH5 == "" {
		return "<h2>No Image</h2>"
	}
	return imgBodyH5
}

func buildVideoTag(fileName string) string {
	return "<video src=\"" + url.PathEscape(fileName) + "\" controls=\"controls\"></video>"
}

func buildImgTag(fileName string) string {
	return "<img  src=\"" + fileName + "\" alt=\"" + fileName + "\">"
}
