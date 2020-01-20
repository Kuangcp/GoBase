package main

import "github.com/kuangcp/gobase/myth-bookkeeping/service"

func main() {
	// 建立数据库结构
	service.AutoMigrateAll()
}
