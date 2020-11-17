red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
blue='\033[0;34m'
purple='\033[0;35m'
cyan='\033[0;36m'
white='\033[0;37m'
end='\033[0m'

run_app(){
    type notify-send >/dev/null 2>&1 || { echo >&2 "notify-send not installed.  Aborting."; exit 1; }

    make buildSingle 

    if [ ! -d bin/data ]; then 
        ln -s data bin/data
    fi 

    content='<span color="#57dafd" font="26px"> <a href="http://localhost:9090/static/">Enter</a></span>'

    notify-send -t 3000 "MyBook" "$content"

    bin/mybook -s -p 9090
}

down(){
    mkdir -p conf/static/js/lib 
    wget https://cdn.bootcdn.net/ajax/libs/jquery/3.5.1/jquery.min.js -O conf/static/js/lib/jquery.min.js 
    wget https://cdn.bootcdn.net/ajax/libs/echarts/4.8.0/echarts.min.js -O conf/static/js/lib/echarts.min.js
    wget https://cdn.jsdelivr.net/npm/vue/dist/vue.js -O conf/static/js/lib/vue.js
}

help(){
    printf "Run：$red sh $0 $green<verb> $yellow<args>$end\n"
    format="  $green%-6s $yellow%-8s$end%-20s\n"
    printf "$format" "-h" "" "帮助"
}

case $1 in 
    -h)
        help ;;
    -d) down ;;
    *)
       run_app ;;
esac