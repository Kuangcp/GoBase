#!/usr/bin/env bash

red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
blue='\033[0;34m'
purple='\033[0;35m'
cyan='\033[0;36m'
white='\033[0;37m'
end='\033[0m'

run_app(){
    cd mybook-static \
    && npm run build \
    && sed -i 's/="\/css/="css/g;s/="\/js/="js/g' dist/index.html \
    && cp ../conf/static/favicon.ico dist/favicon.ico \
    && cd .. \
    && statik -f -src=mybook-static/dist/ -dest app/common/ \
    && go build -o bin/${BINARY_NAME}

    if [ ! -d bin/data ]; then 
        ln -s data bin/data
    fi

    bin/mybook -s -p 9090
}

help(){
    printf "Run：$red sh $0 $green<verb> $yellow<args>$end\n"
    format="  $green%-6s $yellow%-8s$end%-20s\n"
    printf "$format" "-h" "" "帮助"
}

case $1 in 
    -h)
        help
    ;;
    -d)
        down
    ;;
    *)
        run_app
    ;;
esac