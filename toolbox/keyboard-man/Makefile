build:
	statik -f -src=static/ -dest app/ && go build

web:
	statik -f -src=static/ -dest app/ && go build && ./keyboard-man -ws

install:
	statik -f -src=static/ -dest app/ && go install