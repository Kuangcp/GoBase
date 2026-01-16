package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	
	"regexp"
	"strings"
)

//go:embed doc
var javaFS embed.FS

// MarkdownSection 表示一个 markdown 标题及其内容
type MarkdownSection struct {
	Level   int      // 标题级别 (1-6)
	Title   string   // 标题文本
	Content []string // 该标题下的所有行（包括子标题）
}

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	// 检查是否是列出标题模式
	if os.Args[1] == "-l" {
		listAllHeadings()
		return
	}

	keyword := strings.ToLower(os.Args[1])

	// 搜索所有 markdown 文件
	found := false
	err := fs.WalkDir(javaFS, "doc", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只处理 .md 文件
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		// 读取文件内容
		content, err := javaFS.ReadFile(path)
		if err != nil {
			return err
		}

		// 查找匹配的标题并提取内容
		sections := extractMatchingSections(string(content), keyword)
		if len(sections) > 0 {
			if found {
				fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
			}
			found = true
			fmt.Printf("📄 文件: %s\n\n", path)
			for _, section := range sections {
				fmt.Println(strings.Join(section, "\n"))
				fmt.Println()
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	if !found {
		fmt.Printf("未找到关于 '%s' 的文档\n", os.Args[1])
		os.Exit(1)
	}
}

// extractMatchingSections 提取所有匹配关键词的标题及其内容
func extractMatchingSections(content string, keyword string) [][]string {
	lines := strings.Split(content, "\n")
	var results [][]string

	// 正则匹配 markdown 标题
	headerRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		matches := headerRegex.FindStringSubmatch(line)

		if matches != nil {
			level := len(matches[1])
			title := strings.TrimSpace(matches[2])

			// 检查标题是否包含关键词（不区分大小写）
			if strings.Contains(strings.ToLower(title), keyword) {
				// 提取该标题及其所有子内容
				section := extractSection(lines, i, level)
				results = append(results, section)
			}
		}
	}

	return results
}

// extractSection 提取从指定位置开始的一个标题节及其所有子标题内容
func extractSection(lines []string, startIdx int, baseLevel int) []string {
	var section []string
	headerRegex := regexp.MustCompile(`^(#{1,6})\s+`)

	// 添加当前标题行
	section = append(section, lines[startIdx])

	// 遍历后续行，直到遇到同级或更高级的标题
	for i := startIdx + 1; i < len(lines); i++ {
		line := lines[i]
		matches := headerRegex.FindStringSubmatch(line)

		if matches != nil {
			level := len(matches[1])
			// 如果遇到同级或更高级的标题，停止
			if level <= baseLevel {
				break
			}
		}

		section = append(section, line)
	}

	return section
}

func showUsage() {
	fmt.Println("用法: jdoc <关键词>")
	fmt.Println("      jdoc -l")
	fmt.Println("\n选项:")
	fmt.Println("  -l          列出所有文件的标题")
	fmt.Println("\n说明:")
	fmt.Println("  在 doc 目录的所有 markdown 文件中搜索匹配的标题")
	fmt.Println("  并显示该标题下的所有内容（包括子标题）")
	fmt.Println("\n示例:")
	fmt.Println("  jdoc jcmd      # 搜索并显示 jcmd 相关内容")
	fmt.Println("  jdoc jstack    # 搜索并显示 jstack 相关内容")
	fmt.Println("  jdoc -l        # 列出所有标题")
}

// listAllHeadings 列出所有文件的标题
func listAllHeadings() {
	headerRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

	err := fs.WalkDir(javaFS, "doc", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只处理 .md 文件
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		// 读取文件内容
		content, err := javaFS.ReadFile(path)
		if err != nil {
			return err
		}

		lines := strings.Split(string(content), "\n")
		var headings []string

		for _, line := range lines {
			matches := headerRegex.FindStringSubmatch(line)
			if matches != nil {
				level := len(matches[1])
				title := strings.TrimSpace(matches[2])
				// 使用缩进表示层级
				indent := strings.Repeat("  ", level-1)
				headings = append(headings, fmt.Sprintf("%s%s %s", indent, matches[1], title))
			}
		}

		if len(headings) > 0 {
			fmt.Printf("\n📄 %s\n", path)
			for _, heading := range headings {
				fmt.Println(heading)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
}
