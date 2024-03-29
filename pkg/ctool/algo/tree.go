package algo

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"math"
	"strings"
)

// IBinTree 二叉树接口
type IBinTree interface {
	GetLeft() IBinTree
	GetRight() IBinTree
	ToString() string
}

// INTree 多叉树接口
type INTree interface {
	ToString() string
	GetChild() []INTree
}

func Height(tree IBinTree) int {
	if ctool.IsNil(tree) {
		return 0
	}
	maxVal := math.Max(float64(Height(tree.GetLeft())), float64(Height(tree.GetRight())))
	return int(maxVal) + 1
}

// PrintBiMindMap https://plantuml.com/mindmap-diagram
func PrintBiMindMap(tree IBinTree) string {
	return printBiMindMap(tree, 1, "")
}

// Plugin: PlantUML Integration. file suffix: .pu .puml
func WriteBiMindMap(tree IBinTree, file string) {
	mindMap := printBiMindMap(tree, 1, "")
	writer, err := ctool.NewWriter(file, true)
	if err != nil {
		return
	}
	defer writer.Close()
	writer.WriteLine("@startmindmap")
	writer.WriteString(mindMap)
	writer.WriteLine("@endmindmap")
}

func printBiMindMap(tree IBinTree, level int, ct string) string {
	if ctool.IsNil(tree) {
		return ct
	}

	ct += fmt.Sprintln(strings.Repeat("*", level), tree.ToString())
	ct = printBiMindMap(tree.GetLeft(), level+1, ct)
	ct = printBiMindMap(tree.GetRight(), level+1, ct)
	return ct
}

func PrintNMindMap(tree INTree) string {
	return printNMindMap(tree, 1, "")
}

func WriteNMindMap(tree INTree, file string) {
	mindMap := printNMindMap(tree, 1, "")
	writer, err := ctool.NewWriter(file, true)
	if err != nil {
		return
	}
	defer writer.Close()
	writer.WriteLine("@startmindmap")
	writer.WriteString(mindMap)
	writer.WriteLine("@endmindmap")
}

func printNMindMap(tree INTree, level int, ct string) string {
	if ctool.IsNil(tree) {
		return ct
	}

	ct += fmt.Sprintln(strings.Repeat("*", level), tree.ToString())
	if len(tree.GetChild()) > 0 {
		for _, c := range tree.GetChild() {
			if !ctool.IsNil(c) {
				ct = printNMindMap(c, level+1, ct)
			}
		}
	}
	return ct
}
