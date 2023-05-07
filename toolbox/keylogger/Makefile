BIN=keylogger.bin

down:
	sh down_js.sh

heap:
	go tool pprof -inuse_space http://localhost:8891/debug/pprof/heap

heapSvg:
	go tool pprof -inuse_space -cum -svg http://localhost:8891/debug/pprof/heap > heap_inuse.svg

