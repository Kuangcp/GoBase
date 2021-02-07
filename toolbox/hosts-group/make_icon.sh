#/bin/sh

# doc: https://github.com/getlantern/systray/tree/master/example

imgFile=$1

if [ -z "$GOPATH" ]; then
  echo GOPATH environment variable not set
  exit
fi

if [ ! -e "$GOPATH/bin/rsrc" ]; then
  echo "Installing rsrc..."
  go get github.com/akavel/rsrc
  if [ $? -ne 0 ]; then
    echo Failure executing go get github.com/akavel/rsrc
    exit
  fi
fi

rsrc -manifest hosts-group.exe.mainfest -ico $imgFile -o hosts.syso


if [ ! -e "$GOPATH/bin/2goarray" ]; then
  echo "Installing 2goarray..."
  go get github.com/cratonica/2goarray
  if [ $? -ne 0 ]; then
    echo Failure executing go get github.com/cratonica/2goarray
    exit
  fi
fi

if [ -z "$imgFile" ]; then
  echo Please specify a PNG file
  exit
fi

if [ ! -f "$imgFile" ]; then
  echo $imgFile is not a valid file
  exit
fi

mkdir -p app/icon/

OUTPUT=app/icon/iconunix.go

rm $OUTPUT

echo Generating $OUTPUT

cat "$imgFile" | $GOPATH/bin/2goarray Data icon >>$OUTPUT

if [ $? -ne 0 ]; then
  echo Failure generating $OUTPUT
  exit
fi

echo Finished
