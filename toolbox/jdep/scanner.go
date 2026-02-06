package main

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
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

// DiffLine 表示 diff 的一行：Op 为 "-"（仅在第一份）、"+"（仅在第二份）、" "（两边相同）
type DiffLine struct {
	Op   string // "-", "+", " "
	Line string
}

// LineDiff 对两份内容按行做 LCS 对齐，返回 git 风格的 diff 行序列（删除=第一份独有标 "-"，新增=第二份独有标 "+"）
func LineDiff(contentFirst, contentOther string) []DiffLine {
	a := strings.Split(contentFirst, "\n")
	b := strings.Split(contentOther, "\n")
	n, m := len(a), len(b)
	if n == 0 && m == 0 {
		return nil
	}
	// dp[i][j] = LCS length of a[:i] and b[:j]
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] > dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}
	var out []DiffLine
	i, j := n, m
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			out = append(out, DiffLine{" ", a[i-1]})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			out = append(out, DiffLine{"+", b[j-1]})
			j--
		} else {
			out = append(out, DiffLine{"-", a[i-1]})
			i--
		}
	}
	// 反转使顺序为从文件开头到结尾
	for l, r := 0, len(out)-1; l < r; l, r = l+1, r-1 {
		out[l], out[r] = out[r], out[l]
	}
	return out
}

// DiffTwoFiles 读取两个文件并返回相对第一个的 diff 行（第二个相对第一个：- 表示仅第一个有，+ 表示仅第二个有）
func DiffTwoFiles(pathFirst, pathOther string) ([]DiffLine, error) {
	data1, err := os.ReadFile(pathFirst)
	if err != nil {
		return nil, err
	}
	data2, err := os.ReadFile(pathOther)
	if err != nil {
		return nil, err
	}
	return LineDiff(string(data1), string(data2)), nil
}

// getModuleFromPath 从路径中解析出模块名（相对 projectRoot 的第一级目录，即含 pom.xml 的模块目录名）
func getModuleFromPath(path, projectRoot string) string {
	rel, err := filepath.Rel(projectRoot, path)
	if err != nil {
		return ""
	}
	rel = filepath.ToSlash(rel)
	parts := strings.Split(rel, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

// moduleKeepPriority 用于“保留谁”的优先级，数值越小越优先保留（api > dao > service > 其他）
func moduleKeepPriority(moduleName string) int {
	name := strings.ToLower(moduleName)
	switch {
	case strings.Contains(name, "api"):
		return 0
	case strings.Contains(name, "dao"):
		return 1
	case strings.Contains(name, "service"):
		return 2
	default:
		return 99
	}
}

// SortDuplicateListForKeep 对同名类的多个文件按“优先保留”排序：优先保留顶级模块（api > dao > service），同模块取第一个；返回新切片，保留项在 [0]
func SortDuplicateListForKeep(projectRoot string, list []JavaFile) []JavaFile {
	if len(list) <= 1 {
		return list
	}
	root, _ := filepath.Abs(projectRoot)
	copied := make([]JavaFile, len(list))
	copy(copied, list)
	// 稳定排序：先按模块优先级，再按路径
	sort.Slice(copied, func(i, j int) bool {
		mi, mj := getModuleFromPath(copied[i].Path, root), getModuleFromPath(copied[j].Path, root)
		pi, pj := moduleKeepPriority(mi), moduleKeepPriority(mj)
		if pi != pj {
			return pi < pj
		}
		return copied[i].Path < copied[j].Path
	})
	return copied
}

// collectAllJavaPaths 返回项目下所有 .java 文件路径（用于批量替换 import）
func collectAllJavaPaths(projectRoot string) ([]string, error) {
	root, ok := findMavenProjectRoot(projectRoot)
	if !ok {
		return nil, nil
	}
	modDirs, err := collectPomDirs(root)
	if err != nil {
		return nil, err
	}
	var allJava []string
	for _, dir := range modDirs {
		files, _ := collectJavaFiles(dir)
		allJava = append(allJava, files...)
	}
	return allJava, nil
}

// ReplaceImportInProject 将项目中所有 Java 文件里的 import oldFQCN 改为 import newFQCN；
// excludePaths 中的文件不修改（如即将被删除的重复类文件）；返回被修改过的文件路径列表。
func ReplaceImportInProject(projectRoot, oldFQCN, newFQCN string, excludePaths map[string]struct{}) (modified []string, err error) {
	paths, err := collectAllJavaPaths(projectRoot)
	if err != nil || len(paths) == 0 {
		return nil, err
	}
	// 匹配整行：空白 + import + oldFQCN + 空白 + ;
	reImportOld := regexp.MustCompile(`(?m)^(\s*)import\s+` + regexp.QuoteMeta(oldFQCN) + `\s*;\s*$`)
	reHasNew := regexp.MustCompile(`\bimport\s+` + regexp.QuoteMeta(newFQCN) + `\s*;`)

	for _, path := range paths {
		if _, skip := excludePaths[path]; skip {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		content := string(data)
		if !reImportOld.MatchString(content) {
			continue
		}
		hasNew := reHasNew.MatchString(content)
		lines := strings.Split(content, "\n")
		var out []string
		changed := false
		for _, line := range lines {
			if reImportOld.MatchString(line) {
				if hasNew {
					changed = true
					continue
				}
				changed = true
				subs := reImportOld.FindStringSubmatch(line)
				indent := ""
				if len(subs) > 1 {
					indent = subs[1]
				}
				out = append(out, indent+"import "+newFQCN+";")
				continue
			}
			out = append(out, line)
		}
		if !changed {
			continue
		}
		if err := os.WriteFile(path, []byte(strings.Join(out, "\n")), 0644); err != nil {
			return modified, err
		}
		modified = append(modified, path)
	}
	return modified, nil
}
