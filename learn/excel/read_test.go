package excel

import (
	"encoding/xml"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"github.com/xuri/excelize/v2"
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

func TestToCsvSpeed(t *testing.T) {
	path := "data/以岭-感冒药_238.xlsx"
	// 8782ms
	excelizeBuild(path, "data")
	//zipBuild(path, "data")
}

func TestZipBuild(t *testing.T) {
	// 5300ms Python pandas 需要 17993ms
	zipBuild("./data/以岭-感冒药_238.xlsx", "data")
	// 2min Python pandas 需要7分钟
	//zipBuild("./data/英诺珐-月度-三品类-自定义品名_362.xlsx", "data")
}

func TestFileName(t *testing.T) {
	logger.Info(extractFileName("/home/kcp/test/tocsv/short-parts-price.xlsx"))
}

func TestParseXml(t *testing.T) {
	row := "<x><row r=\"1\" spans=\"1:6\"><c r=\"A1\"><v>1</v></c><c r=\"B1\"><v>2</v></c><c r=\"C1\"><v>3</v></c><c r=\"D1\"><v>4</v></c><c r=\"E1\"><v>5</v></c><c r=\"F1\"><v>6</v></c></row>" +
		"<row r=\"1\" spans=\"1:6\"><c r=\"A1\"><v>5</v></c><c r=\"B1\"><v>2</v></c><c r=\"C1\"><v>3</v></c><c r=\"D1\"><v>4</v></c><c r=\"E1\"><v>5</v></c><c r=\"F1\"><v>6</v></c></row></x>"
	r := Rows{}
	err := xml.Unmarshal([]byte(row), &r)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(r)
}
