module mybook

go 1.13

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/jinzhu/gorm v1.9.12
	github.com/kuangcp/gobase/pkg/cuibase v0.0.0-20201103041857-ea5c95ff0199
	github.com/kuangcp/gobase/pkg/ginhelper v0.0.0-20201103041857-ea5c95ff0199
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/rakyll/statik v0.1.7
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/wonderivan/logger v1.0.0
)

replace github.com/kuangcp/gobase/pkg/cuibase => ./../pkg/cuibase
