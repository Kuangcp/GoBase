BIN=hosts-group.bin

buildExe:
	sh make_exe.sh
install:
	statik -f -src=static/ -dest app && go install
build:
	statik -f -src=static/ -dest app && go build -o ${BIN}
web:
	statik -f -src=static/ -dest app && go build -o ${BIN} && ./${BIN} -d