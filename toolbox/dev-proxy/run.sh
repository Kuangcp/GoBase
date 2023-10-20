red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
blue='\033[0;34m'
purple='\033[0;35m'
cyan='\033[0;36m'
white='\033[0;37m'
end='\033[0m'
grey='\033[3;37;40m'

log_info(){
    i=$(($i + 1))
    ts=$(date "+%F %T.%N" | cut -b 1-23)
    printf "$green$ts $i/$total $blue $1 $end\n"
}

build_image(){
    total=4
    i=0
    last_commit=$(git log --oneline | head -n 1 | awk '{print $1}')
    dt=$(date +%m%d%H%M)

    log_info 'Compile'
    go build -o dev-proxy.bin
    upx dev-proxy.bin

    log_info 'Stop exist container'
    docker stop dev-proxy

    log_info 'Delete container'
    docker rm dev-proxy

    log_info 'Build image'
    # yay docker-buildx
    docker buildx build -t dev-proxy:1.0-$dt-$last_commit .

    log_info 'Start new container'
    docker run --init -d --name dev-proxy --network host -v $HOME/Apps/dev-proxy-cd:/root/.dev-proxy/ dev-proxy:1.0-$dt-$last_commit

    log_info 'Finish rebuild'
}

case $1 in 
clean)

;;

*)
   build_image
;;
esac
