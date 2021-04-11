module mybook

go 1.13

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kuangcp/gobase/pkg/cuibase v0.0.0-20201103041857-ea5c95ff0199
	github.com/kuangcp/gobase/pkg/ghelp v0.0.0-20201103041857-ea5c95ff0199
	github.com/kuangcp/gobase/pkg/stopwatch v0.0.0-20210409094425-3f724b872d91 // indirect
	github.com/kuangcp/logger v1.0.3
	github.com/mattn/go-sqlite3 v1.14.4
	github.com/rakyll/statik v0.1.7
	github.com/stretchr/testify v1.6.1
	github.com/ugorji/go v1.1.13 // indirect
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/sys v0.0.0-20201101102859-da207088b7d1 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace github.com/kuangcp/gobase/pkg/cuibase => ./../pkg/cuibase

replace github.com/kuangcp/gobase/pkg/ghelp => ./../pkg/ghelp
