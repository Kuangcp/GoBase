BIN=keylogger.bin

down:
	mkdir -p static/js/lib && wget https://cdn.bootcdn.net/ajax/libs/jquery/3.5.1/jquery.min.js -O static/js/lib/jquery.min.js && wget https://cdn.bootcdn.net/ajax/libs/echarts/4.8.0/echarts.min.js -O static/js/lib/echarts.min.js

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

