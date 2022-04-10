package record

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/kuangcp/gobase/pkg/cuibase"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestReadExport(t *testing.T) {
	// 交易号 商家订单号 交易创建时间 付款时间 最近修改时间
	// 交易来源地 类型 交易对方 商品名称 金额（元）
	// 收/支 交易状态 服务费（元） 成功退款（元） 备注
	// 资金状态

	db, _ := gorm.Open("sqlite3", "./ali_tmp.db")

	// 11
	filterT := "2021-10"
	printIdx := []int{2, 3, 5, 6, 7, 8, 9, 10, 11, 13, 15}
	//printIdx := []int{12}

	var sum float64
	table := cuibase.ReadCsvFile("/home/kcp/Documents/2021支付宝 handle.csv")
	for _, row := range table {
		inOrOut := strings.TrimSpace(row[10])
		//if inOrOut != "支出" && inOrOut != "收入" {
		if inOrOut != "其他" {
			continue
		}
		if !strings.HasPrefix(row[2], filterT) {
			continue
		}
		if !strings.Contains(row[8], "买入") {
			//if !strings.Contains(row[8], "卖出") {
			continue
		}

		ss, err := strconv.ParseFloat(strings.TrimSpace(row[9]), 64)
		if err != nil {
			continue
		}
		sum += ss
		//fmt.Println(row[9])
		for _, idx := range printIdx {
			fmt.Print(row[idx] + "|")
		}
		fmt.Println()

		db.Exec("select 1")
		//db.Exec("insert into ali_re values(?, ?,?,?)", row[3], row[7], row[10], row[9])
		//db.Exec("insert into ali_france values(?, ?,?,?)", row[2], row[7], row[10], row[9])

		//logger.Info(row[2], row[6], row[8], row[10])
		//fmt.Println()
		// ct time, name varchar, t varchar, amount number
		//fmt.Printf("insert into a (%v ,%v ,%v, %v)", row[2], row[7], row[10], row[9])
		//fmt.Println()
	}

	fmt.Println(sum)

}

func TestFrance(t *testing.T) {
	type AmountVO struct {
		month  string
		amount float64
	}
	cache := make(map[string]*AmountVO)
	//printIdx := []int{2, 3, 5, 6, 7, 8, 9, 10, 11, 13, 15}
	table := cuibase.ReadCsvFile("/home/kcp/Documents/2021支付宝 handle.csv")
	for _, row := range table {
		//inOrOut := strings.TrimSpace(row[10])

		if !strings.Contains(row[8], "买入") && !strings.Contains(row[8], "卖出") {
			continue
		}
		in := "out"
		if strings.Contains(row[8], "买入") {
			in = "in"
		}

		//for _, idx := range printIdx {
		//	fmt.Print(row[idx] + "|")
		//}
		//fmt.Println()
		m := row[2][:7]

		ss, err := strconv.ParseFloat(strings.TrimSpace(row[9]), 64)
		if err != nil {
			continue
		}

		ex, ok := cache[m+in]
		if ok {
			ex.amount += ss
		} else {
			cache[m+in] = &AmountVO{month: m, amount: ss}
		}

	}

	var ks []string
	for s, _ := range cache {
		ks = append(ks, s)
	}

	sort.Strings(ks)
	for i := range ks {
		//fmt.Println(ks[i])
		fmt.Println(ks[i], cache[ks[i]].amount)
	}

	//for k, v := range cache {
	//	fmt.Println(k, v.amount)
	//}
}
