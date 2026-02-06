package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// JavaFile 表示一个 Java 源文件解析结果
type JavaFile struct {
	Path        string   // 完整路径
	Package     string   // 包名
	Class       string   // 当前文件主类名（与文件名一致）
	FullClass   string   // 全限定名 package.Class
	Imports     []string // import 的全限定类名（不含 import static 的 *）
	IsController bool   // 是否为 Controller 层（含 @RestController / @Controller），此类不被依赖视为正常
}

// 预编译正则，避免在热路径重复编译（regexp.Regexp 并发只读安全）
var (
	rePackage    = regexp.MustCompile(`(?m)^\s*package\s+([\w.]+)\s*;`)
	reImport     = regexp.MustCompile(`(?m)^\s*import\s+(?:static\s+)?([\w.]+)(?:\.\*)?\s*;`)
	reWord       = regexp.MustCompile(`\b([A-Z][a-zA-Z0-9]*)\b`) // 简单识别首字母大写的标识符（类名）
	reController = regexp.MustCompile(`@(?:Rest)?Controller\b`) // @Controller 或 @RestController
)

// findMavenProjectRoot 从 dir 向上查找包含 pom.xml 的目录，若 dir 自身有 pom.xml 则返回 dir
func findMavenProjectRoot(dir string) (string, bool) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", false
	}
	for {
		if _, err := os.Stat(filepath.Join(abs, "pom.xml")); err == nil {
			return abs, true
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return "", false
		}
		abs = parent
	}
}

// collectPomDirs 收集 root 下所有包含 pom.xml 的目录（含 root 自身及子模块）
func collectPomDirs(root string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == "pom.xml" {
			dirs = append(dirs, filepath.Dir(path))
			return nil
		}
		return nil
	})
	return dirs, err
}

// collectJavaFiles 在 moduleRoot 下收集 src/main/java 和 src/test/java 中的 .java 文件
func collectJavaFiles(moduleRoot string) ([]string, error) {
	var list []string
	//rootDirs := []string{"src/main/java", "src/test/java"}
	rootDirs := []string{"src/main/java"}
	for _, sub := range rootDirs {
		dir := filepath.Join(moduleRoot, sub)
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				return err
			}
			if info.IsDir() || !strings.HasSuffix(path, ".java") {
				return nil
			}
			list = append(list, path)
			return nil
		})
	}
	return list, nil
}

// parseJavaFile 解析一个 Java 文件，得到包名、类名和 import 列表；content 为文件原文，供后续“类名出现”阶段复用，避免二次读盘。
func parseJavaFile(path string) (*JavaFile, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	content := string(data)

	pkg := ""
	if m := rePackage.FindStringSubmatch(content); len(m) > 1 {
		pkg = strings.TrimSpace(m[1])
	}

	class := strings.TrimSuffix(filepath.Base(path), ".java")
	fqcn := class
	if pkg != "" {
		fqcn = pkg + "." + class
	}

	var imports []string
	for _, m := range reImport.FindAllStringSubmatch(content, -1) {
		if len(m) > 1 && !strings.HasSuffix(m[1], "*") {
			imports = append(imports, strings.TrimSpace(m[1]))
		}
	}

	isController := reController.MatchString(content)

	return &JavaFile{
		Path:         path,
		Package:      pkg,
		Class:        class,
		FullClass:    fqcn,
		Imports:      imports,
		IsController: isController,
	}, data, nil
}

// FindNoDepClasses 扫描 Maven 项目，返回未被任何其他类的 import 或源码引用的类（仅限 import + 同包名引用）
func FindNoDepClasses(projectRoot string) ([]JavaFile, error) {
	root, ok := findMavenProjectRoot(projectRoot)
	if !ok {
		return nil, nil
	}

	modDirs, err := collectPomDirs(root)
	if err != nil {
		return nil, err
	}

	// 所有 .java 文件（各模块并发收集）
	var allJava []string
	var listMu sync.Mutex
	var wg sync.WaitGroup
	for _, dir := range modDirs {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			files, _ := collectJavaFiles(d)
			if len(files) > 0 {
				listMu.Lock()
				allJava = append(allJava, files...)
				listMu.Unlock()
			}
		}(dir)
	}
	wg.Wait()

	// 解析每个文件（并发读文件+解析），并缓存文件内容供阶段二复用，避免二次读盘
	defined := make(map[string]JavaFile)
	contentByPath := make(map[string][]byte)
	var definedMu sync.Mutex
	var files []*JavaFile
	var filesMu sync.Mutex
	var contentMu sync.Mutex
	for _, p := range allJava {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			jf, content, err := parseJavaFile(path)
			if err != nil {
				return
			}
			definedMu.Lock()
			defined[jf.FullClass] = *jf
			definedMu.Unlock()
			filesMu.Lock()
			files = append(files, jf)
			filesMu.Unlock()
			contentMu.Lock()
			contentByPath[path] = content
			contentMu.Unlock()
		}(p)
	}
	wg.Wait()

	// 被引用的类：1) 出现在任意文件的显式 import 中（不处理 import *）
	// 2) 在任意类文件源码中作为「词」出现（含跨包、含 import * 引入的类）
	referenced := make(map[string]struct{})

	for _, jf := range files {
		for _, imp := range jf.Imports {
			referenced[imp] = struct{}{}
		}
	}

	// 源码中出现类名（词）：复用已读内容，每文件只做一次 reWord 扫描，用 set 做 O(1) 查找；用 FindAll([]byte) 避免整文件转 string
	var refMu sync.Mutex
	for _, jf := range files {
		wg.Add(1)
		go func(j *JavaFile) {
			defer wg.Done()
			content := contentByPath[j.Path]
			if len(content) == 0 {
				return
			}
			wordBytes := reWord.FindAll(content, -1)
			wordSet := make(map[string]struct{}, len(wordBytes))
			for _, b := range wordBytes {
				wordSet[string(b)] = struct{}{}
			}
			for fqcn := range defined {
				if fqcn == j.FullClass {
					continue
				}
				simple := fqcn[strings.LastIndex(fqcn, ".")+1:]
				if simple == j.Class {
					continue
				}
				if _, ok := wordSet[simple]; ok {
					refMu.Lock()
					referenced[fqcn] = struct{}{}
					refMu.Unlock()
				}
			}
		}(jf)
	}
	wg.Wait()

	var noDep []JavaFile
	for fqcn, jf := range defined {
		if _, ref := referenced[fqcn]; !ref && !jf.IsController {
			noDep = append(noDep, jf)
		}
	}
	return noDep, nil
}

// FindDuplicateClasses 扫描 Maven 项目，返回简单类名重复的类：类名 -> 多个 JavaFile（含路径）
func FindDuplicateClasses(projectRoot string) (map[string][]JavaFile, error) {
	root, ok := findMavenProjectRoot(projectRoot)
	if !ok {
		return nil, nil
	}
	modDirs, err := collectPomDirs(root)
	if err != nil {
		return nil, err
	}

	var allJava []string
	var listMu sync.Mutex
	var wg sync.WaitGroup
	for _, dir := range modDirs {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			files, _ := collectJavaFiles(d)
			if len(files) > 0 {
				listMu.Lock()
				allJava = append(allJava, files...)
				listMu.Unlock()
			}
		}(dir)
	}
	wg.Wait()

	var allFiles []JavaFile
	var allMu sync.Mutex
	for _, p := range allJava {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			jf, _, err := parseJavaFile(path)
			if err != nil {
				return
			}
			allMu.Lock()
			allFiles = append(allFiles, *jf)
			allMu.Unlock()
		}(p)
	}
	wg.Wait()

	// 按简单类名分组，只保留出现多于一次的
	bySimple := make(map[string][]JavaFile)
	for _, jf := range allFiles {
		bySimple[jf.Class] = append(bySimple[jf.Class], jf)
	}
	dup := make(map[string][]JavaFile)
	for name, list := range bySimple {
		if len(list) > 1 {
			dup[name] = list
		}
	}
	return dup, nil
}

// normalizeJavaContent 做基础规范化便于比较：去掉注释、规整空白（忽略字符串内注释等边界情况）
func normalizeJavaContent(data []byte) string {
	s := string(data)
	// 块注释 /* ... */
	for {
		start := strings.Index(s, "/*")
		if start == -1 {
			break
		}
		end := strings.Index(s[start:], "*/")
		if end == -1 {
			break
		}
		s = s[:start] + " " + s[start+end+2:]
	}
	// 行注释 // ...
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		idx := strings.Index(line, "//")
		if idx != -1 {
			lines[i] = line[:idx]
		}
	}
	// 对比时忽略 package 行（同名类可能在不同包下）
	var filtered []string
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "package ") && strings.HasSuffix(t, ";") {
			continue
		}
		filtered = append(filtered, line)
	}
	s = strings.Join(filtered, "\n")
	// 空白规整：连续空白变为单空格并 trim
	var b strings.Builder
	lastSpace := true
	for _, r := range strings.TrimSpace(s) {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if !lastSpace {
				b.WriteRune(' ')
				lastSpace = true
			}
			continue
		}
		lastSpace = false
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

// levenshteinDistance 编辑距离（rune 为单位）
func levenshteinDistance(a, b []rune) int {
	na, nb := len(a), len(b)
	if na == 0 {
		return nb
	}
	if nb == 0 {
		return na
	}
	dp := make([]int, nb+1)
	for j := 0; j <= nb; j++ {
		dp[j] = j
	}
	for i := 1; i <= na; i++ {
		prev := dp[0]
		dp[0] = i
		for j := 1; j <= nb; j++ {
			curr := dp[j]
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			dp[j] = min(min(dp[j]+1, dp[j-1]+1), prev+cost)
			prev = curr
		}
	}
	return dp[nb]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// contentDiffPct 返回 0~100 的差异百分比（基于规范化后的内容，编辑距离 / max(len)）
func contentDiffPct(normA, normB string) float64 {
	ra, rb := []rune(normA), []rune(normB)
	maxLen := len(ra)
	if len(rb) > maxLen {
		maxLen = len(rb)
	}
	if maxLen == 0 {
		return 0
	}
	dist := levenshteinDistance(ra, rb)
	return float64(dist) / float64(maxLen) * 100
}

// CompareDuplicateContents 对同名类的多个文件读入并比较：是否完全一致，以及各文件与首文件的差异百分比
func CompareDuplicateContents(list []JavaFile) (identical bool, diffPcts []float64) {
	if len(list) <= 1 {
		return true, nil
	}
	norm := make([]string, len(list))
	for i := range list {
		data, err := os.ReadFile(list[i].Path)
		if err != nil {
			norm[i] = ""
			continue
		}
		norm[i] = normalizeJavaContent(data)
	}
	first := norm[0]
	diffPcts = make([]float64, len(list))
	diffPcts[0] = 0
	identical = true
	for i := 1; i < len(list); i++ {
		pct := contentDiffPct(first, norm[i])
		diffPcts[i] = pct
		if pct > 0 {
			identical = false
		}
	}
	return identical, diffPcts
}
