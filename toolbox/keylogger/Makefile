build:
	statik -f -src=static/ -dest app/ && go build

web:
	statik -f -src=static/ -dest app/ && go build && ./keylogger -ws

install:
	statik -f -src=static/ -dest app/ && go install