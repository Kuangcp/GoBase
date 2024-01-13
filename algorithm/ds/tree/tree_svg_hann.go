package main

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	cycle_r = 25
)

// https://blog.csdn.net/boysoft2002/article/details/126908846
type btNode struct {
	Data   any
	Lchild *btNode
	Rchild *btNode
}

type biTree struct {
	Root *btNode
	Info *biTreeInfo
}

type biTreeInfo struct {
	Data                []any
	DataLevel           [][]any
	L, R                []bool
	X, Y, W             []int
	Index, Nodes        int
	Width, Height       int
	MarginX, MarginY    int
	SpaceX, SpaceY      int
	SvgWidth, SvgHeight int
	SvgXml              string
}

func build(Data ...any) *biTree {
	if len(Data) == 0 || Data[0] == nil {
		return &biTree{}
	}
	node := &btNode{Data: Data[0]}
	Queue := []*btNode{node}
	for lst := Data[1:]; len(lst) > 0 && len(Queue) > 0; {
		cur, val := Queue[0], lst[0]
		Queue, lst = Queue[1:], lst[1:]
		if val != nil {
			cur.Lchild = &btNode{Data: val}
			Queue = append(Queue, cur.Lchild)
		}
		if len(lst) > 0 {
			val, lst = lst[0], lst[1:]
			if val != nil {
				cur.Rchild = &btNode{Data: val}
				Queue = append(Queue, cur.Rchild)
			}
		}
	}
	return &biTree{Root: node}
}

func buildByTree[T ctool.Ordered](root *BsNode[T]) *biTree {
	if root == nil {
		return &biTree{}
	}

	return &biTree{Root: fillNode(root)}
}

func fillNode[T ctool.Ordered](root *BsNode[T]) *btNode {
	if root == nil {
		return nil
	}
	node := &btNode{Data: root.Data}
	node.Lchild = fillNode(root.Left)
	node.Rchild = fillNode(root.Right)
	return node
}

func BuildFromList(List []any) *biTree {
	return build(List...)
}

func AinArray(sub int, array []int) int {
	for idx, arr := range array {
		if sub == arr {
			return idx
		}
	}
	return -1
}

func Pow2(x int) int { //x>=0
	return int(math.Pow(2, float64(x)))
}

func Max(L, R int) int {
	if L > R {
		return L
	} else {
		return R
	}
}

func (bt *btNode) MaxDepth() int {
	if bt == nil {
		return 0
	}
	return 1 + Max(bt.Lchild.MaxDepth(), bt.Rchild.MaxDepth())
}

func (bt *btNode) Coordinate(x, y, w int) []any {
	var res []any
	if bt != nil {
		L, R := bt.Lchild != nil, bt.Rchild != nil
		res = append(res, []any{bt.Data, L, R, x, y, w})
		res = append(res, bt.Lchild.Coordinate(x-w, y+1, w/2)...)
		res = append(res, bt.Rchild.Coordinate(x+w, y+1, w/2)...)
	}
	return res
}

func (bt *biTree) NodeInfo() []any {
	return bt.Root.Coordinate(0, 0, Pow2(bt.Root.MaxDepth()-2))
}

func (bt *biTree) buildTreeInfo() {
	height := bt.Root.MaxDepth()
	width := Pow2(height - 1)
	lsInfo := bt.NodeInfo()
	btInfo := &biTreeInfo{
		Height: height,
		Width:  width,
		Nodes:  len(lsInfo),
	}
	for _, data := range lsInfo {
		for i, info := range data.([]any) {
			switch i {
			case 0:
				btInfo.Data = append(btInfo.Data, info.(any))
			case 1:
				btInfo.L = append(btInfo.L, info.(bool))
			case 2:
				btInfo.R = append(btInfo.R, info.(bool))
			case 3:
				btInfo.X = append(btInfo.X, info.(int))
			case 4:
				btInfo.Y = append(btInfo.Y, info.(int))
			case 5:
				btInfo.W = append(btInfo.W, info.(int))
			}
		}
	}
	for j, k := 0, width*2; j < height; j++ {
		DLevel := []any{}
		for i := k / 2; i < width*2; i += k {
			index := AinArray(i-width, btInfo.X)
			if index > -1 {
				DLevel = append(DLevel, btInfo.Data[index])
			} else {
				DLevel = append(DLevel, nil)
			}
			DLevel = append(DLevel, []int{i, j})
			if k/4 == 0 {
				DLevel = append(DLevel, []int{0, 0})
				DLevel = append(DLevel, []int{0, 0})
			} else {
				DLevel = append(DLevel, []int{i - k/4, j + 1})
				DLevel = append(DLevel, []int{i + k/4, j + 1})
			}
		}
		k /= 2
		btInfo.DataLevel = append(btInfo.DataLevel, DLevel)
	}
	bt.Info = btInfo
}

func (bt *biTree) info2SVG(Margin ...int) string {
	var res, Line, Color string
	info := bt.Info
	MarginX, MarginY := 10, 10
	SpaceX, SpaceY := 28, 90
	switch len(Margin) {
	case 0:
		break
	case 1:
		MarginX = Margin[0]
	case 2:
		MarginX, MarginY = Margin[0], Margin[1]
	case 3:
		MarginX, MarginY, SpaceX = Margin[0], Margin[1], Margin[2]
	default:
		MarginX, MarginY = Margin[0], Margin[1]
		SpaceX, SpaceY = Margin[2], Margin[3]
	}
	info.MarginX, info.MarginY = MarginX, MarginY
	info.SpaceX, info.SpaceY = SpaceX, SpaceY
	info.SvgWidth = Pow2(info.Height)*info.SpaceX + info.SpaceX
	info.SvgHeight = info.Height * info.SpaceY
	for i, Data := range info.Data {
		Node := "\n\t<g id=\"INDEX,M,N\">\n\t<CIRCLE/>\n\t<TEXT/>\n\t<LEAF/>\n\t</g>"
		DataStr := ""
		switch Data.(type) {
		case int:
			DataStr = strconv.Itoa(Data.(int))
		case float64:
			DataStr = strconv.FormatFloat(Data.(float64), 'g', -1, 64)
		case string:
			DataStr = Data.(string)
		default:
			DataStr = "Error Type"
		}
		Node = strings.Replace(Node, "INDEX", strconv.Itoa(info.Index), 1)
		Node = strings.Replace(Node, "M", strconv.Itoa(info.X[i]), 1)
		Node = strings.Replace(Node, "N", strconv.Itoa(info.Y[i]), 1)
		x0, y0 := (info.X[i]+info.Width)*SpaceX+MarginX, 50+info.Y[i]*SpaceY+MarginY
		x1, y1 := x0-info.W[i]*SpaceX, y0+SpaceY-cycle_r
		x2, y2 := x0+info.W[i]*SpaceX, y0+SpaceY-cycle_r
		Color = "lightgreen"
		offset := 18
		if info.L[i] && info.R[i] {
			Line = XmlLine(x0-offset, y0+offset, x1, y1) + "\n\t" + XmlLine(x0+offset, y0+offset, x2, y2)
		} else if info.L[i] && !info.R[i] {
			Line = XmlLine(x0-offset, y0+offset, x1, y1)
		} else if !info.L[i] && info.R[i] {
			Line = XmlLine(x0+offset, y0+offset, x2, y2)
		} else {
			Color = "lightgreen"
		}
		Node = strings.Replace(Node, "<CIRCLE/>", XmlCircle(x0, y0, Color), 1)
		Node = strings.Replace(Node, "<TEXT/>", XmlText(x0, y0, DataStr), 1)
		if info.L[i] || info.R[i] {
			Node = strings.Replace(Node, "<LEAF/>", Line, 1)
		}
		res += Node
	}
	info.SvgXml = res
	return res
}

func XmlCircle(X, Y int, Color string) string {
	Circle := "<circle cx=\"" + strconv.Itoa(X) + "\" cy=\"" + strconv.Itoa(Y) +
		"\" r=\"" + fmt.Sprint(cycle_r) + "\" stroke=\"black\" stroke-width=" +
		"\"2\" style=\"fill:" + Color + ";stroke:black;opacity:0.5\"/>"
	return Circle
}

func XmlText(X, Y int, DATA string) string {
	iFontSize, tColor := 25, "black"
	Text := "<text x=\"" + strconv.Itoa(X) + "\" y=\"" + strconv.Itoa(Y) +
		"\" fill=\"" + tColor + "\" font-size=\"" + strconv.Itoa(iFontSize) +
		"\" text-anchor=\"middle\" dominant-baseline=\"middle\">" + DATA + "</text>"
	return Text
}

func XmlLine(X1, Y1, X2, Y2 int) string {
	Line := "<line x1=\"" + strconv.Itoa(X1) + "\" y1=\"" + strconv.Itoa(Y1) +
		"\" x2=\"" + strconv.Itoa(X2) + "\" y2=\"" + strconv.Itoa(Y2) +
		"\" style=\"stroke:black;stroke-width:2;opacity:0.7\" />"
	return Line
}

func (bt *biTree) showSVG(FileName ...string) {
	var file *os.File
	var err1 error
	Head := "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink" +
		"=\"http://www.w3.org/1999/xlink\" version=\"1.1\" width=" +
		"\"Width\" height=\"Height\">\nLINKCONTENT\n</svg>"
	Link := `<a xlink:href="https://blog.csdn.net/boysoft2002" target="_blank">
	<text x="5" y="20" fill="blue">Hann's CSDN Homepage</text></a>`
	Xml := strings.Replace(Head, "LINK", Link, 1)
	Xml = strings.Replace(Xml, "Width", strconv.Itoa(bt.Info.SvgWidth), 1)
	Xml = strings.Replace(Xml, "Height", strconv.Itoa(bt.Info.SvgHeight), 1)
	Xml = strings.Replace(Xml, "CONTENT", bt.Info.SvgXml, 1)
	svgFile := "biTree.svg"
	if len(FileName) > 0 {
		svgFile = FileName[0] + ".svg"
	}
	file, err1 = os.Create(svgFile)
	if err1 != nil {
		panic(err1)
	}
	_, err1 = io.WriteString(file, Xml)
	if err1 != nil {
		panic(err1)
	}
	file.Close()
	exec.Command("cmd", "/c", "start", svgFile).Start()
	//Linux 代码：
	//exec.Command("xdg-open", svgFile).Start()
	//Mac 代码：
	//exec.Command("open", svgFile).Start()
}
