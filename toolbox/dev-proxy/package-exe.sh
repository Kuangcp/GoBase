#!/bin/sh

# doc: https://github.com/getlantern/systray/tree/master/example
# 1-4 complete convert /static/favicon.svg to app/icon for tray icon use
# svgFile=static/favicon.svg
svgFile=core/static/favicon.svg
imgFile=core/static/favicon.ico
OUTPUT=gui/iconunix.go

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
#isCmdExist upx upx

#############################
# start build
#############################

echo 1/6 Start build

if [ ! -f static/favicon.ico ]; then
  # convert -density 300 $svgFile -background none -colors 256 -define icon:auto-resize static/favicon.ico
  convert -background none -density 300 $svgFile -define icon:auto-resize $imgFile
  echo 2/6 Finish convert svg to ico
else
  echo 2/6 Finish convert svg to ico by cache
fi

if [ -z "$GOPATH" ]; then
  echo GOPATH environment variable not set
  exit
fi

if [ -z "$imgFile" ]; then
  echo Please specify a PNG file
  exit
fi

if [ ! -f "$imgFile" ]; then
  echo $imgFile is not a valid file
  exit
fi

# mkdir -p app/icon/
rm $OUTPUT
cat "$imgFile" | $GOPATH/bin/2goarray Data main >>$OUTPUT
if [ $? -ne 0 ]; then
  echo Failure generating $OUTPUT
  exit
fi
echo 3/6 Finish Generating $OUTPUT

go mod tidy

echo 4/6 Finish package exe icon, tray icon

##########
# Windows
##########

#if [ ! -f hosts.syso ]; then
#  rsrc -manifest hosts-group.exe.mainfest -ico $imgFile -o hosts.syso
#  echo 5/6 Windows Finish build icon syso
#else
#  echo 5/6 Windows Finish build icon syso by cache
#fi
#
#version=$(go run . -v)
#
#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H windowsgui" -o hosts-group.$version.exe
#
#if [ -f hosts-group.$version.exe ]; then
#  echo
#  upx hosts-group.$version.exe
#  echo
#
#  echo 6/6 Windows Finish hosts-group.$version.exe
#fi
