module github.com/kuangcp/gobase/keyboard-man

go 1.14

require (
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/gvalkov/golang-evdev v0.0.0-20191114124502-287e62b94bcb
	github.com/kuangcp/gobase/cuibase v0.0.0-20200409163938-76c18e7b9704
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
)

replace github.com/kuangcp/gobase/cuibase => ./../../cuibase
