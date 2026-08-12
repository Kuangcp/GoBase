package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestReplace(t *testing.T) {
	title := "pre（）测试 【使用】"

	prepareContext()
	result := normalizeForTitle(title)
	fmt.Println(result)

	if result != "pre测试-使用" {
		t.Fail()
	}
}

func TestNormalizeForTitle(t *testing.T) {
	prepareContext()
	cases := []struct {
		in   string
		want string
	}{
		{"pre（）测试 【使用】", "pre测试-使用"},
		{"一、LLM 原理入门", "一llm-原理入门"},
		{"C++ 语法", "c-语法"},
		{"Hello World", "hello-world"},
		{"标题：带冒号", "标题带冒号"},
		{"A/B 测试", "ab-测试"},
		{"问号？与斜杠/", "问号与斜杠"},
	}
	for _, c := range cases {
		if got := normalizeForTitle(c.in); got != c.want {
			t.Errorf("normalizeForTitle(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestPathToString(t *testing.T) {
	cases := []struct {
		in   []int
		want string
	}{
		{[]int{1}, "1"},
		{[]int{1, 2}, "1.2"},
		{[]int{1, 2, 3}, "1.2.3"},
		{[]int{}, ""},
	}
	for _, c := range cases {
		if got := pathToString(c.in); got != c.want {
			t.Errorf("pathToString(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestBuildTitle(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"LLM.md", "LLM"},
		{"/a/b/LLM.md", "LLM"},
		{"/a/b/LLM", "LLM"},
		{"/a/b/c/note.markdown", "note.markdown"},
		{"relative/path/foo.md", "foo"},
		{"noext", "noext"},
	}
	for _, c := range cases {
		if got := buildTitle(c.in); got != c.want {
			t.Errorf("buildTitle(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestIsNeedHandleFile(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"foo.md", true},
		{"foo.markdown", true},
		{"foo.txt", true},
		{"foo.go", false},
		{"foo.pdf", false},
		{"README.md", false},
		{"dir/readme.md", false},
		{"dir/SUMMARY.md", false},
		{"dir/LICENSE.md", false},
		{"dir/Process.md", false},
	}
	for _, c := range cases {
		if got := isNeedHandleFile(c.in); got != c.want {
			t.Errorf("isNeedHandleFile(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestRefreshTitle(t *testing.T) {
	prepareContext()
	//RefreshTagAndCatalog("test.md")
}

func TestBuildArticleUnformattedDash(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "LLM.md")
	content := "# Title\n\nbody\n\n---\n\nmore body\n"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	a := BuildArticle(file)
	if a == nil {
		t.Fatal("article should not be nil")
	}
	if len(a.tag) != 0 {
		t.Fatalf("tag should be empty for unformatted file with '---' in body, got: %q", a.tag)
	}
	if len(a.catalog) != 0 {
		t.Fatalf("catalog should be empty, got: %q", a.catalog)
	}
	if len(a.content) == 0 {
		t.Fatal("content should hold the whole file")
	}
}

func TestBuildArticleUnformattedHeaderLast(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "X.md")
	content := "# Title\n\nbody\n\n" + headerLast + "\n# Section2\n"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	a := BuildArticle(file)
	if a == nil {
		t.Fatal("article should not be nil")
	}
	if len(a.tag) != 0 {
		t.Fatalf("tag should be empty, got: %q", a.tag)
	}
	if len(a.catalog) != 0 {
		t.Fatalf("catalog should be empty, got: %q", a.catalog)
	}
	if len(a.content) == 0 {
		t.Fatal("content should hold the whole file")
	}
}
