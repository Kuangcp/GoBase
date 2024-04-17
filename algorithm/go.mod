module algorithm

go 1.21

toolchain go1.22.0

require (
	github.com/kuangcp/gobase/pkg/ctool v1.1.9
	github.com/stretchr/testify v1.8.4
	github.com/tidwall/pretty v1.2.1
	golang.org/x/exp v0.0.0-20240103183307-be819d1f06fc
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/kuangcp/gobase/pkg/ctool => ../pkg/ctool
