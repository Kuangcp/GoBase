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
