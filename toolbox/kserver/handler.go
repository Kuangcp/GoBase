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
	"strings"
	"time"
)

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

func buildImgFunc(parentPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		rawSize := query.Get("rawSize")

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
			<style>`))

		if rawSize == "" {
			w.Write([]byte(`
			img {
				width: 210px;
				max-height: 100px;
			}`))
		}

		w.Write([]byte(`</style></head><body>`))

		sort.Slice(dir, func(i, j int) bool {
			iInfo, _ := dir[i].Info()
			jInfo, _ := dir[j].Info()
			return iInfo.ModTime().After(jInfo.ModTime())
		})

		for _, entry := range dir {
			if entry.IsDir() {
				continue
			}

			fileName := entry.Name()
			idx := strings.LastIndex(fileName, ".")
			if idx == -1 {
				writeImgTag(w, fileName)
				continue
			}
			suffixType := fileName[idx:]
			if suffixType == ".jpg" || suffixType == ".png" || suffixType == ".svg" || suffixType == ".webp" ||
				suffixType == ".bmp" || suffixType == ".gif" || suffixType == ".ico" {
				writeImgTag(w, fileName)
			}
		}
		w.Write([]byte(`</body></html>`))
	}
}

func writeImgTag(w http.ResponseWriter, fileName string) (int, error) {
	return w.Write([]byte("<img  src=\"" + fileName + "\" alt=\"" + fileName + "\">"))
}
