#/bin/sh

# doc: https://github.com/getlantern/systray/tree/master/example

if [ -z "$GOPATH" ]; then
  echo GOPATH environment variable not set
  exit
fi

if [ ! -e "$GOPATH/bin/2goarray" ]; then
  echo "Installing 2goarray..."
  go get github.com/cratonica/2goarray
  if [ $? -ne 0 ]; then
    echo Failure executing go get github.com/cratonica/2goarray
    exit
  fi
fi

if [ -z "$1" ]; then
  echo Please specify a PNG file
  exit
fi

if [ ! -f "$1" ]; then
  echo $1 is not a valid file
  exit
fi

mkdir -p app/icon/

OUTPUT=app/icon/iconunix.go

rm $OUTPUT

echo Generating $OUTPUT

cat "$1" | $GOPATH/bin/2goarray Data icon >>$OUTPUT

if [ $? -ne 0 ]; then
  echo Failure generating $OUTPUT
  exit
fi

echo Finished
