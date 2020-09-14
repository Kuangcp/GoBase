down:
	mkdir -p static/js/lib && wget https://cdn.bootcdn.net/ajax/libs/jquery/3.5.1/jquery.min.js -O static/js/lib/jquery.min.js && wget https://cdn.bootcdn.net/ajax/libs/echarts/4.8.0/echarts.min.js -O static/js/lib/echarts.min.js

build:
	statik -f -src=static/ -dest app/ && go build

web:
	statik -f -src=static/ -dest app/ && go build && ./keylogger -ws

install:
	statik -f -src=static/ -dest app/ && go install
