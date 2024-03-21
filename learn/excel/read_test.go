package excel

import (
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
