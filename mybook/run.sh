#!/usr/bin/env bash

red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
blue='\033[0;34m'
purple='\033[0;35m'
cyan='\033[0;36m'
white='\033[0;37m'
end='\033[0m'

BINARY_NAME=mybook.bin

install_app() {
    cd mybook-static &&
        npm run build &&
        cp ../../toolbox/keylogger/static/favicon.ico dist/favicon.ico &&
        cd .. &&
        statik -f -src=mybook-static/dist/ -dest app/common/

    install_server
}

install_server() {
    go build -o bin/${BINARY_NAME}

    if [ ! -d bin/data ]; then
        ln -s data bin/data
    fi

    bin/${BINARY_NAME} -s -p 9090
}

run(){
  bin/${BINARY_NAME} -s -p 9090
}

help() {
    printf "Run：$red sh $0 $green<verb> $yellow<args>$end\n"
    format="  $green%-6s $yellow%-8s$end%-20s\n"
    printf "$format" "-h" "" "帮助"
    printf "$format" "-i" "" "初始化配置"
    printf "$format" "-r" "" "运行已经本地编译前后端的应用"
    printf "$format" "-s" "" "仅编译并运行后端应用"
    printf "$format" "" "" "编译前端，编译并运行后端应用"
}

case $1 in
-h)
    help
    ;;
-i)
    mkdir data && echo -e 'db: \n    file: ./data/main.db\ndebug: true' >> data/mybook.yml &&
     go test -v -test.run TestInit
    ;;
-c)
  rm -f bin/${BINARY_NAME}
  ;;
-r)
    run
    ;;
-s)
    # only backend
    install_server
    ;;
*)
    # contain front and backend
    install_app
    ;;
esac
