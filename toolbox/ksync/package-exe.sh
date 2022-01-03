#/bin/sh

# doc: https://github.com/getlantern/systray/tree/master/example

imgFile=favicon.ico
originImg=sync.svg
OUTPUT=icon/iconunix.go

down_dep() {
  cmd=$1
  url=$2

  if [ ! -e "$GOPATH/bin/"$1 ]; then
    echo "Installing rsrc..."
    go get $2
    if [ $? -ne 0 ]; then
      echo Failure executing go get github.com/akavel/rsrc
      exit
    fi
  fi
}

function isCmdExist() {
  which "$1" >/dev/null 2>&1
  if [ $? -eq 0 ]; then
    return 0
  fi

  echo "need install $2 for command $1"
  exit 0
}

#############################
# install dep
#############################
down_dep rsrc github.com/akavel/rsrc
down_dep 2goarray github.com/cratonica/2goarray

isCmdExist convert imagemagick
isCmdExist upx upx

#############################
# start build
#############################

echo 1/6 Start build

if [ ! -f favicon.ico ]; then
  convert -density 300 $originImg -background transparent -colors 256 -define icon:auto-resize $imgFile
  echo 2/6 Finish convert svg to ico
else
  echo 2/6 Finish convert svg to ico by cache
fi

if [ -z "$GOPATH" ]; then
  echo GOPATH environment variable not set
  exit
fi

if [ ! -f "$imgFile" ]; then
  echo $imgFile is not a valid file
  exit
fi

if [ ! -f hosts.syso ]; then
  rsrc -manifest ksync.exe.mainfest -ico $imgFile -o hosts.syso
  echo 3/6 Finish build icon syso
else
  echo 3/6 Finish build icon syso by cache
fi

mkdir -p icon/
rm $OUTPUT
cat "$imgFile" | $GOPATH/bin/2goarray Data icon >>$OUTPUT
if [ $? -ne 0 ]; then
  echo Failure generating $OUTPUT
  exit
fi
echo 4/6 Finish Generating $OUTPUT

go mod tidy

echo 5/6 Finish package exe icon, tray icon

version=$(go run . -v)

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H windowsgui" -o ksync.$version.exe

if [ -f ksync.$version.exe ]; then
  echo
  upx ksync.$version.exe
  echo

  echo 6/6 Finish ksync.$version.exe
fi
