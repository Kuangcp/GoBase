module github.com/kuangcp/gobase/mybook

go 1.13

require (
	github.com/gin-gonic/gin v1.5.0
	github.com/jinzhu/gorm v1.9.12
	github.com/kuangcp/gobase/cuibase v0.0.0-20200120172943-d8144d065aaf
	github.com/mattn/go-sqlite3 v2.0.2+incompatible
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.4.0
	github.com/wonderivan/logger v1.0.0
)

replace github.com/kuangcp/gobase/cuibase => ./../cuibase
