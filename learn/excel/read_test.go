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
	"testing"
)

func TestReadParse(t *testing.T) {
	writer, err := ctool.NewWriter("crm.sql", true)
	if err != nil {
		return
	}
	defer writer.Close()

	f, err := excelize.OpenFile("/home/zk/old/Downloads/crm_export_20220920_all.xlsx")
	if err != nil {
		println(err.Error())
		return
	}

	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("客户数据")

	sql := ""
	for _, row := range rows {
		rowSql := "insert into crm_master_company_export values("
		for _, colCell := range row {
			//print(colCell, "\t")
			if strings.Contains(colCell, "'") {
				rowSql += fmt.Sprintf("'%v',", strings.Replace(colCell, "'", "''", -1))
			} else {
				rowSql += fmt.Sprintf("'%v',", colCell)
			}
		}

		if len(row) == 41 {
			rowSql += "'','',"
		}

		rowSql = rowSql[:len(rowSql)-1] + ");\n"

		sql += rowSql
		//if i > 5 {
		//	break
		//}
	}
	writer.WriteString(sql)
	//fmt.Println(sql)
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

func TestFileName(t *testing.T) {
	logger.Info(extractFileName("/home/kcp/test/tocsv/short-parts-price.xlsx"))
}

func TestToCsvSpeed(t *testing.T) {
	path := "/home/kcp/test/tocsv/short-parts-price.xlsx"
	//excelizeBuild(path)
	zipBuild(path)
}

func excelizeBuild(path string) {
	file, err := excelize.OpenFile(path)
	if err != nil {
		logger.Info(err)
		return
	}
	rows, err := file.Rows("Sheet1")
	if err != nil {
		logger.Error(err)
		return
	}
	writer, _ := ctool.NewWriter(extractFileName(path)+".csv", true)
	defer writer.Close()

	for rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			logger.Error(err)
			continue
		}
		logger.Info(cols)
		writer.WriteLine(strings.Join(cols, ","))
	}
}

type (
	Row struct {
		Cell []Val `xml:"c"`
	}
	Val struct {
		Val string `xml:"v"`
	}
)

func TestParseXml(t *testing.T) {
	row := "<row r=\"1\" spans=\"1:6\"><c r=\"A1\"><v>1</v></c><c r=\"B1\"><v>2</v></c><c r=\"C1\"><v>3</v></c><c r=\"D1\"><v>4</v></c><c r=\"E1\"><v>5</v></c><c r=\"F1\"><v>6</v></c></row>"
	r := Row{}
	err := xml.Unmarshal([]byte(row), &r)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(r)
}

// https://stackoverflow.com/questions/38208137/processing-large-xlsx-file-in-python
// https://blog.51cto.com/yuzhou1su/5165546
func zipBuild(path string) {
	x, err := os.Open(path)
	if err != nil {
		return
	}
	all, err := io.ReadAll(x)
	if err != nil {
		return
	}
	r, err := zip.NewReader(bytes.NewReader(all), int64(len(all)))
	for _, f := range r.File {
		if !strings.Contains(f.Name, "xl/worksheets") {
			continue
		}
		fmt.Println(f)

		r, err := f.Open()
		if err != nil {
			logger.Error(err)
			panic(err)
		}
		//
		//buf, err := io.ReadAll(r)
		//if err != nil {
		//	logger.Error(err)
		//	panic(err)
		//}
		//
		//fmt.Println(string(buf))

		writer, _ := ctool.NewWriter(extractFileName(f.Name)+".csv", true)
		defer writer.Close()

		rowBuf := ""
		buffer := bufio.NewReader(r)
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

			fmt.Println(len(bts))
			//fmt.Println(">>>>>>", string(bts), "<<<<<<")
			rowBuf += string(bts)
			rowBuf = handleRowCells(rowBuf, writer)

			x++
			if x > 5 {
				break
			}
		}
	}
}

// TODO 字符串列，日期列, 处理多行
// sharedStrings.xml 还要去取回被压缩的字典值，字符串列在sheet中只是存了索引号
// 格式： <si><t xml:space="preserve">备注</t></si>
func handleRowCells(buf string, writer *ctool.BufferWriter) string {
	start := strings.Index(buf, "<row ")
	end := strings.Index(buf, "</row>")
	if start == -1 || end == -1 {
		return buf
	}
	row := buf[start : end+6]
	r := Row{}
	fmt.Println(row)
	err := xml.Unmarshal([]byte(row), &r)
	if err != nil {
		logger.Error(err)
		return buf
	}
	logger.Info(r)
	return buf[end+6:]
}
