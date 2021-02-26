BIN=keylogger.bin

down:
	sh down_js.sh

statik:
	go get github.com/rakyll/statik

build:
	statik -f -src=static/ -dest app/ && go build -o ${BIN}

web:
	statik -f -src=static/ -dest app/ && go build -o ${BIN} && ./${BIN} -S

install:
	statik -f -src=static/ -dest app/ && go install

heap:
	go tool pprof -inuse_space http://localhost:8891/debug/pprof/heap

heapSvg:
	go tool pprof -inuse_space -cum -svg http://localhost:8891/debug/pprof/heap > heap_inuse.svg

