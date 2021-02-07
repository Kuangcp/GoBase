BIN=hosts-group.bin

install:
	statik -f -src=static/ -dest app && go install
web:
	statik -f -src=static/ -dest app && go build -o ${BIN} && ./${BIN} -d