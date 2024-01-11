package algo

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"math"
)

const (
	header = "<?xml version=\"1.0\" standalone=\"no\"?><!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">"
	wpad   = 2
	hpad   = 8
)

type ()

// TODO https://blog.csdn.net/boysoft2002/article/details/126908846 树绘制成标准svg图 或者直接2D绘制窗口
// https://www.w3schools.com/graphics/svg_line.asp
// https://blog.csdn.net/boysoft2002/article/details/126908846
// 层次遍历 标记节点在一个满二叉树里的序号，依据序号计算出 x y 坐标
func BuildSvg(tree IBinTree) string {
	height, weight := 0, 0
	//writer.WriteString("<svg width=\"400\" height=\"180\">\n " +
	//	" <rect x=\"50\" y=\"40\" rx=\"20\" ry=\"40\" width=\"50\" height=\"50\" style=\"fill:red;stroke:black;opacity:0.5\" />\n" +
	//	"  Sorry, your browser does not support inline SVG.\n " +
	//	" <line x1=\"110\" y1=\"0\" x2=\"60\" y2=\"40\" style=\"stroke:rgb(0,0,0);stroke-width:2\" />\n  " +
	//	"<line x1=\"110\" y1=\"0\" x2=\"160\" y2=\"40\" style=\"stroke:rgb(0,0,0);stroke-width:2\" />\n    " +
	//	"<rect x=\"120\" y=\"40\" rx=\"20\" ry=\"40\" width=\"50\" height=\"50\" style=\"fill:red;stroke:black;opacity:0.5\" />\n" +
	//	"</svg>")

	maxH := Height(tree)

	var cur []IBinTree
	cur = appendIfAbsent(cur, tree.GetLeft())
	cur = appendIfAbsent(cur, tree.GetRight())
	layer := 1
	for {
		layer++
		if len(cur) == 0 {
			break
		}

		var nextLayer []IBinTree
		for _, node := range cur {
			wp := pow2(max(maxH-layer-1, 0)) * wpad
			fmt.Println(wp)

			nextLayer = appendIfAbsent(nextLayer, node.GetLeft())
			nextLayer = appendIfAbsent(nextLayer, node.GetRight())
		}

		if len(nextLayer) == 0 {
			break
		}

		cur = nextLayer
	}

	start := fmt.Sprintf("<svg width=\"%v\" height=\"%v\" version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\">", weight, height)
	return header + start + "" + "</svg>"
}

func appendIfAbsent(cur []IBinTree, node IBinTree) []IBinTree {
	if ctool.IsNil(node) {
		return cur
	}
	return append(cur, node)
}

func buildRect(x, y, rx, ry int) string {
	return fmt.Sprintf("<rect x=\"%v\" y=\"%v\" rx=\"%v\" ry=\"%v\" width=\"50\" height=\"50\""+
		" style=\"fill:red;stroke:black;opacity:0.5\" />\n", x, y, rx, ry)
}
func pow2(pow int) int {
	return int(math.Pow(2, float64(pow)))
}
func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
