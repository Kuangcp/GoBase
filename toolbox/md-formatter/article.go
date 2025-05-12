package main

import (
	"fmt"
	"github.com/kuangcp/logger"
	"os"
	"strings"
	"time"
)

type Article struct {
	filename string
	tag      []string
	catalog  []string
	content  []string
}

func BuildArticle(filename string) *Article {
	lines := readFileLines(filename)
	if len(lines) == 0 {
		return nil
	}

	var tag []string
	header := false
	tagEnd := false
	catalogMatch := false
	catalogIdx := 0
	contentIdx := 0
	for i, line := range lines {
		line = strings.Replace(line, "**目录 ", splitTag, 1)
		line = strings.Replace(line, "\r\n", "\n", 1)
		if line == headerLast {
			header = false
			contentIdx = i + 1
			break
		}

		// 兼容脏数据
		if strings.Contains(line, splitTag) && catalogMatch {
			header = false
			contentIdx = i + 2
			break
		}
		// 找到第一个 catalog 定位行
		if strings.Contains(line, splitTag) && !catalogMatch {
			catalogIdx = i
			catalogMatch = true
		}

		if line == headerFirst && header {
			tagEnd = true
			tag = append(tag, line)
			continue
		}
		if line == headerFirst {
			header = true
		}
		if header {
			if !tagEnd {
				tag = append(tag, line)
			}
			continue
		}
	}
	if catalogIdx > contentIdx {
		return nil
	}
	return &Article{filename: filename, tag: tag, catalog: lines[catalogIdx:contentIdx], content: lines[contentIdx:]}
}

func (a *Article) writeToDisk(hiddenCatalog bool) {
	content := ""
	if len(a.tag) > 0 {
		content += strings.Join(a.tag, "")
	}
	if len(a.catalog) > 0 && !hiddenCatalog {
		content += strings.Join(a.catalog, "")
	}
	content += strings.Join(a.content, "")

	if os.WriteFile(a.filename, []byte(content), 0644) != nil {
		logger.Error("write error", a.filename)
	}
}

func (a *Article) generateCatalog() {
	var pPath []int
	var catalog []string
	code := false
	catalog = append(catalog, "\n"+splitTag+"\n\n")
	for _, line := range a.content {

		if strings.HasPrefix(line, codeBlock) {
			code = !code
		}
		if code {
			continue
		}
		if !strings.HasPrefix(line, "#") {
			continue
		}

		if len(pPath) == 0 {
			pPath = []int{0}
		}
		level := strings.Count(line, "#")
		for len(pPath) < level {
			pPath = append(pPath, 0)
		}
		pPath[level-1] += 1
		if level < len(pPath) {
			for i := level; i < len(pPath); i++ {
				pPath[i] = 0
			}
		}

		title := strings.TrimSpace(strings.Replace(line, "#", "", -1))
		strings.Count(line, "#")
		temps := strings.Split(line, "# ")
		levelStr := strings.Replace(temps[0], "#", "    ", -1)
		row := fmt.Sprintf("%s- %s. [%s](#%s)\n", levelStr, pathToString(pPath[:level]), title, normalizeForTitle(title))
		catalog = append(catalog, row)
	}
	catalog = append(catalog, "\n"+splitTag+" "+time.Now().Format("2006-01-02 15:04:05")+"\n", headerLast)
	a.catalog = catalog
}

// 创建/保留 tag 删除/创建 catalog 保留content
func (a *Article) Refresh() {
	// 处理 tag
	if len(a.tag) == 0 {
		var tag = fmt.Sprintf(tagTemplate, buildTitle(a.filename), time.Now().Format("2006-01-02 15:04:05"))
		rows := strings.Split(tag, "\n")
		for _, r := range rows {
			a.tag = append(a.tag, r+"\n")
		}
	}

	// 处理 catalog
	a.generateCatalog()
}
