package excel

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"github.com/xuri/excelize/v2"
	"github.com/yeka/zip"
	"io"
	"os"
	"strings"
)

var (
	debug    = false
	debugLog = false
)

func excelizeBuild(path, dstDir string) {
	file, err := excelize.OpenFile(path)
	if err != nil {
		logger.Info(err)
		return
	}
	//TODO 通用处理时 Sheet名如何获取
	rows, err := file.Rows("数据源")
	if err != nil {
		logger.Error(err)
		return
	}
	writer, _ := ctool.NewWriter(dstDir+"/"+extractFileName(path)+".lize.csv", true)
	defer writer.Close()

	for rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			logger.Error(err)
			continue
		}
		//logger.Info(cols)
		writer.WriteLine(strings.Join(cols, ","))
	}
}

type (
	Rows struct {
		R []Row `xml:"row"`
	}
	Row struct {
		Cell []Cell `xml:"c"`
	}
	Cell struct {
		Val string `xml:"v"`  // number
		Is  Text   `xml:"is"` // text
	}
	Text struct {
		T string `xml:"t"`
	}
)

func (c Cell) String() string {
	if c.Val == "" {
		return c.Is.T
	}
	return c.Val
}

// https://stackoverflow.com/questions/38208137/processing-large-xlsx-file-in-python
// https://blog.51cto.com/yuzhou1su/5165546
func zipBuild(path, dstDir string) {
	bufferSize := 8192
	x, err := os.Open(path)
	if err != nil {
		logger.Error(err)
		return
	}
	all, err := io.ReadAll(x)
	if err != nil {
		logger.Error(err)
		return
	}
	zipFile, err := zip.NewReader(bytes.NewReader(all), int64(len(all)))
	for _, f := range zipFile.File {
		if !strings.Contains(f.Name, "xl/worksheets") {
			continue
		}
		logger.Debug("Handle", f)

		r, err := f.Open()
		if err != nil {
			logger.Error(err)
			panic(err)
		}

		writer, _ := ctool.NewWriter(dstDir+"/"+extractFileName(path)+".csv", true)
		defer writer.Close()

		rowBuf := ""
		buffer := bufio.NewReaderSize(r, bufferSize)
		x := 0
		for {
			// ReadLine is a low-level line-reading primitive.
			// Most callers should use ReadBytes('\n') or ReadString('\n') instead or use a Scanner.
			bts, _, err := buffer.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Error(err)
				break
			}

			//logger.Debug(len(bts))
			//fmt.Println(">>>>>>", string(bts), "<<<<<<")
			rowBuf += string(bts)
			if len(rowBuf) > bufferSize {
				rowBuf = handleRowCells(rowBuf, writer)
			}

			x++
			if debug && x > 20 {
				break
			}
		}
		handleRowCells(rowBuf, writer)
	}
}

// TODO 字符串列，日期列, 处理多行
//  1. sharedStrings.xml 还要去取回被压缩的字典值，字符串列在sheet中只是存了索引号
//     格式： <si><t xml:space="preserve">备注</t></si>
//  2. inlineStr 则是直接读
func handleRowCells(buf string, writer *ctool.BufferWriter) string {
	start := strings.Index(buf, "<row")
	end := strings.LastIndex(buf, "</row>")
	if start == -1 || end == -1 {
		return buf
	}
	row := buf[start : end+6]
	r := Rows{}
	if debugLog {
		fmt.Println("RRRRR", row)
	}
	err := xml.Unmarshal([]byte("<x>"+row+"</x>"), &r)
	if err != nil {
		logger.Error(err)
		return buf
	}
	if debugLog {
		logger.Info(r)
	}
	//logger.Info(len(r.R))
	for _, row := range r.R {
		tmp := ""
		for i, c := range row.Cell {
			if i != 0 {
				tmp += ","
			}
			tmp += c.String()
		}
		writer.WriteString(tmp + "\n")
	}

	return buf[end+6:]
}

func extractFileName(path string) string {
	if path == "" {
		return ""
	}

	start := strings.LastIndex(path, "/")
	end := strings.LastIndex(path, ".")
	if end == -1 {
		return path[start+1:]
	}

	return path[start+1 : end]
}
