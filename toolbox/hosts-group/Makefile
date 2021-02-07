BIN=hosts-group.bin

buildExe:
	statik -f -src=static/ -dest app \
	&& convert -density 300 static/favicon.svg -background transparent -colors 256 -define icon:auto-resize static/favicon.ico \
	&& sh make_icon.sh static/favicon.ico \
	&& CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H windowsgui"
install:
	statik -f -src=static/ -dest app && go install
build:
	statik -f -src=static/ -dest app && go build -o ${BIN}
web:
	statik -f -src=static/ -dest app && go build -o ${BIN} && ./${BIN} -d