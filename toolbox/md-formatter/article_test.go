package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return file
}

func readTempFile(t *testing.T, file string) string {
	t.Helper()
	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

const frontMatter = "---\ntitle: x\ndate: 2025-01-01\ntags: \ncategories: \n---\n"

// ---------- BuildArticle ----------

func TestBuildArticleFormatted(t *testing.T) {
	content := frontMatter + "\n💠\n\n- 1. [T](#t)\n\n💠 2025-01-01 00:00:00\n" + headerLast + "正文内容\n"
	file := writeTempFile(t, "a.md", content)

	a := BuildArticle(file)
	if a == nil {
		t.Fatal("article should not be nil")
	}
	if len(a.tag) != 6 {
		t.Fatalf("tag should have 6 lines, got %d: %q", len(a.tag), a.tag)
	}
	if len(a.catalog) == 0 {
		t.Fatal("catalog should not be empty")
	}
	if a.catalog[len(a.catalog)-1] != headerLast {
		t.Fatalf("catalog last line should be headerLast, got %q", a.catalog[len(a.catalog)-1])
	}
	if got := strings.Join(a.content, ""); !strings.HasPrefix(got, "正文内容") {
		t.Fatalf("content should start with 正文内容, got %q", a.content)
	}
}

func TestBuildArticleUnformatted(t *testing.T) {
	content := "# Title\n\n正文\n\n---\n\n" + headerLast + "\n# Section2\n"
	file := writeTempFile(t, "b.md", content)

	a := BuildArticle(file)
	if a == nil {
		t.Fatal("article should not be nil")
	}
	if len(a.tag) != 0 {
		t.Fatalf("tag should be empty, got %q", a.tag)
	}
	if len(a.catalog) != 0 {
		t.Fatalf("catalog should be empty, got %q", a.catalog)
	}
	if len(a.content) == 0 {
		t.Fatal("content should hold the whole file")
	}
}

func TestBuildArticleFrontMatterOnly(t *testing.T) {
	content := frontMatter + "正文内容\n"
	file := writeTempFile(t, "c.md", content)

	a := BuildArticle(file)
	if a == nil {
		t.Fatal("article should not be nil")
	}
	if len(a.tag) != 6 {
		t.Fatalf("tag should have 6 lines, got %d: %q", len(a.tag), a.tag)
	}
	if len(a.catalog) != 0 {
		t.Fatalf("catalog should be empty, got %q", a.catalog)
	}
	joined := strings.Join(a.content, "")
	if !strings.HasPrefix(joined, "正文内容") {
		t.Fatalf("content should start with 正文内容, got %q", a.content)
	}
	if strings.Contains(joined, "title: x") {
		t.Fatalf("content should not contain front matter, got %q", a.content)
	}
}

func TestBuildArticleDirtyData(t *testing.T) {
	// 目录以第二个 💠 结束，其后无 headerLast 分隔行
	content := frontMatter + "\n💠\n\n- 1. [T](#t)\n\n💠 2025-01-01 00:00:00\n\n正文内容\n"
	file := writeTempFile(t, "d.md", content)

	a := BuildArticle(file)
	if a == nil {
		t.Fatal("article should not be nil")
	}
	if len(a.tag) != 6 {
		t.Fatalf("tag should have 6 lines, got %d: %q", len(a.tag), a.tag)
	}
	if len(a.catalog) == 0 {
		t.Fatal("catalog should not be empty")
	}
	if got := strings.Join(a.content, ""); !strings.HasPrefix(got, "正文内容") {
		t.Fatalf("content should start with 正文内容, got %q", a.content)
	}
}

func TestBuildArticleEmpty(t *testing.T) {
	file := writeTempFile(t, "e.md", "")

	if a := BuildArticle(file); a != nil {
		t.Fatalf("empty file should return nil article, got %v", a)
	}
}

// ---------- generateCatalog ----------

func TestGenerateCatalog(t *testing.T) {
	prepareContext()
	a := &Article{
		content: []string{
			"# 一级\n",
			"正文\n",
			"## 二级\n",
			"### 三级\n",
			"```go\n",
			"# 代码块里的标题\n",
			"```\n",
			"## 另一个二级\n",
		},
	}
	a.generateCatalog()

	if len(a.catalog) < 7 {
		t.Fatalf("catalog should have at least 7 lines, got %d: %q", len(a.catalog), a.catalog)
	}
	wantHead := []string{
		"\n💠\n\n",
		"- 1. [一级](#一级)\n",
		"    - 1.1. [二级](#二级)\n",
		"        - 1.1.1. [三级](#三级)\n",
		"    - 1.2. [另一个二级](#另一个二级)\n",
	}
	for i, want := range wantHead {
		if a.catalog[i] != want {
			t.Errorf("catalog[%d] = %q, want %q", i, a.catalog[i], want)
		}
	}
	if a.catalog[len(a.catalog)-1] != headerLast {
		t.Errorf("catalog last should be headerLast, got %q", a.catalog[len(a.catalog)-1])
	}
	ts := a.catalog[len(a.catalog)-2]
	if !strings.HasPrefix(ts, "\n💠 ") {
		t.Errorf("catalog second-to-last should be timestamp line, got %q", ts)
	}
}

func TestGenerateCatalogNormalizeTitle(t *testing.T) {
	prepareContext()
	a := &Article{
		content: []string{
			"# 一、LLM 原理入门\n",
		},
	}
	a.generateCatalog()

	row := a.catalog[1]
	if !strings.Contains(row, "[一、LLM 原理入门](#一llm-原理入门)") {
		t.Errorf("catalog row should normalize title, got %q", row)
	}
}

// ---------- Refresh ----------

func TestRefreshCreateTag(t *testing.T) {
	prepareContext()
	a := &Article{filename: "/a/b/LLM.md", content: []string{"# 标题\n"}}
	a.Refresh()

	if len(a.tag) != 7 {
		t.Fatalf("tag should have 7 lines (6 front matter + trailing blank), got %d: %q", len(a.tag), a.tag)
	}
	if a.tag[0] != "---\n" {
		t.Errorf("tag[0] should be '---\\n', got %q", a.tag[0])
	}
	if a.tag[1] != "title: LLM\n" {
		t.Errorf("tag[1] should be title line, got %q", a.tag[1])
	}
	if len(a.catalog) == 0 {
		t.Fatal("catalog should be generated")
	}
}

func TestRefreshKeepTag(t *testing.T) {
	prepareContext()
	a := &Article{
		filename: "/a/b/LLM.md",
		tag:      []string{"---\n", "title: keep\n", "date: 2020-01-01\n", "tags: \n", "categories: \n", "---\n", "\n"},
		content:  []string{"# 标题\n"},
	}
	a.Refresh()

	if got := strings.Join(a.tag, ""); !strings.Contains(got, "title: keep") {
		t.Errorf("existing tag should be preserved, got %q", got)
	}
	if len(a.catalog) == 0 {
		t.Fatal("catalog should be generated")
	}
}

// ---------- writeToDisk ----------

func TestWriteToDiskWithCatalog(t *testing.T) {
	file := writeTempFile(t, "w1.md", "")
	a := &Article{
		filename: file,
		tag:      []string{"---\n", "title: x\n", "---\n", "\n"},
		catalog:  []string{"\n💠\n\n", "- 1. [T](#t)\n", "\n💠 t\n", headerLast},
		content:  []string{"正文\n"},
	}
	a.writeToDisk(false)

	got := readTempFile(t, file)
	if !strings.Contains(got, "title: x") {
		t.Errorf("should write tag, got %q", got)
	}
	if !strings.Contains(got, "- 1. [T](#t)") {
		t.Errorf("should write catalog, got %q", got)
	}
	if !strings.Contains(got, "正文") {
		t.Errorf("should write content, got %q", got)
	}
}

func TestWriteToDiskHiddenCatalog(t *testing.T) {
	file := writeTempFile(t, "w2.md", "")
	a := &Article{
		filename: file,
		tag:      []string{"---\n", "title: x\n", "---\n", "\n"},
		catalog:  []string{"\n💠\n\n", "- 1. [T](#t)\n", "\n💠 t\n", headerLast},
		content:  []string{"正文\n"},
	}
	a.writeToDisk(true)

	got := readTempFile(t, file)
	if !strings.Contains(got, "title: x") {
		t.Errorf("should write tag, got %q", got)
	}
	if strings.Contains(got, "💠") || strings.Contains(got, "- 1. [T](#t)") {
		t.Errorf("catalog should be hidden, got %q", got)
	}
	if !strings.Contains(got, "正文") {
		t.Errorf("should write content, got %q", got)
	}
}

// ---------- 端到端 ----------

func TestRefreshTagAndCatalog(t *testing.T) {
	prepareContext()
	file := writeTempFile(t, "LLM.md", "# 一级\n\n正文\n\n## 二级\n\n```sh\n# 代码\n```\n")

	RefreshTagAndCatalog(file)

	got := readTempFile(t, file)
	if !strings.Contains(got, "title: LLM") {
		t.Errorf("should contain front matter title, got:\n%s", got)
	}
	if !strings.Contains(got, "- 1. [一级](#一级)") {
		t.Errorf("should contain level-1 catalog entry, got:\n%s", got)
	}
	if !strings.Contains(got, "    - 1.1. [二级](#二级)") {
		t.Errorf("should contain level-2 catalog entry, got:\n%s", got)
	}
	if strings.Contains(got, "[代码](#代码)") {
		t.Errorf("code block heading should be skipped, got:\n%s", got)
	}
	if !strings.Contains(got, "正文") {
		t.Errorf("should preserve content, got:\n%s", got)
	}
}

func TestRemoveCatalog(t *testing.T) {
	prepareContext()
	file := writeTempFile(t, "LLM.md", "# 一级\n\n正文\n")

	RefreshTagAndCatalog(file)
	RemoveCatalog(file)

	got := readTempFile(t, file)
	if strings.Contains(got, "💠") {
		t.Errorf("catalog should be removed, got:\n%s", got)
	}
	if !strings.Contains(got, "title: LLM") {
		t.Errorf("front matter should be kept, got:\n%s", got)
	}
	if !strings.Contains(got, "正文") {
		t.Errorf("content should be kept, got:\n%s", got)
	}
}
