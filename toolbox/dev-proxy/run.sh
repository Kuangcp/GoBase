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
    CName=dev-proxy

    total=6
    i=0
    last_commit=$(git log --oneline | head -n 1 | awk '{print $1}')
    dt=$(date +%m%d%H%M)

    log_info 'Compile'
    go build -o dev-proxy.bin
    upx dev-proxy.bin

    log_info 'Stop exist container '$CName
    docker stop $CName

    log_info 'Delete container '$CName
    docker rm $CName

    log_info 'Build image'
    # yay docker-buildx
    docker buildx build -t dev-proxy:1.0-$dt-$last_commit .

    log_info 'Start new container '$CName
    docker run --init -d --name $CName --network host -v $HOME/.dev-proxy/container:/root/.dev-proxy/ dev-proxy:1.0-$dt-$last_commit
#    docker run --init -d --name $CName --network host dev-proxy:1.0-$dt-$last_commit

    log_info 'Finish rebuild '$CName
}

case $1 in 

clean)

;;
mem.alloc)
  go tool pprof -alloc_space http://localhost:1255/debug/pprof/heap
;;
mem.use)
  go tool pprof -inuse_space http://localhost:1255/debug/pprof/heap
;;

*)
   build_image
;;
esac
