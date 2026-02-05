package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// JavaFile 表示一个 Java 源文件解析结果
type JavaFile struct {
	Path      string   // 完整路径
	Package   string   // 包名
	Class     string   // 当前文件主类名（与文件名一致）
	FullClass string   // 全限定名 package.Class
	Imports   []string // import 的全限定类名（不含 import static 的 *）
}

var (
	rePackage = regexp.MustCompile(`(?m)^\s*package\s+([\w.]+)\s*;`)
	reImport  = regexp.MustCompile(`(?m)^\s*import\s+(?:static\s+)?([\w.]+)(?:\.\*)?\s*;`)
	reWord    = regexp.MustCompile(`\b([A-Z][a-zA-Z0-9]*)\b`) // 简单识别首字母大写的标识符（类名）
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

// parseJavaFile 解析一个 Java 文件，得到包名、类名和 import 列表
func parseJavaFile(path string) (*JavaFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
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

	return &JavaFile{
		Path:      path,
		Package:   pkg,
		Class:     class,
		FullClass: fqcn,
		Imports:   imports,
	}, nil
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

	// 所有 .java 文件
	var allJava []string
	for _, dir := range modDirs {
		files, _ := collectJavaFiles(dir)
		allJava = append(allJava, files...)
	}

	// 解析每个文件
	defined := make(map[string]JavaFile) // fqcn -> JavaFile
	var files []*JavaFile
	for _, p := range allJava {
		jf, err := parseJavaFile(p)
		if err != nil {
			continue
		}
		defined[jf.FullClass] = *jf
		files = append(files, jf)
	}

	// 被引用的类：出现在任意文件的 import 中，或同包下被当作标识符使用
	referenced := make(map[string]struct{})

	for _, jf := range files {
		for _, imp := range jf.Imports {
			referenced[imp] = struct{}{}
		}
	}

	// 同包引用：在源码中出现同包类的简单类名（作为词）
	for _, jf := range files {
		content, err := os.ReadFile(jf.Path)
		if err != nil {
			continue
		}
		text := string(content)
		for fqcn, _ := range defined {
			if fqcn == jf.FullClass {
				continue
			}
			pkg := jf.Package
			if pkg == "" {
				continue
			}
			if !strings.HasPrefix(fqcn, pkg+".") {
				continue
			}
			simple := fqcn[strings.LastIndex(fqcn, ".")+1:]
			if simple == jf.Class {
				continue
			}
			// 在源码中作为词出现
			for _, match := range reWord.FindAllString(text, -1) {
				if match == simple {
					referenced[fqcn] = struct{}{}
					break
				}
			}
		}
	}

	var noDep []JavaFile
	for fqcn, jf := range defined {
		if _, ref := referenced[fqcn]; !ref {
			noDep = append(noDep, jf)
		}
	}
	return noDep, nil
}
