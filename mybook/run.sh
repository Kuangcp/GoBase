type notify-send >/dev/null 2>&1 || { echo >&2 "notify-send not installed.  Aborting."; exit 1; }

content='<span color="#57dafd" font="26px"> <a href="http://localhost:10006/static/">Enter</a></span>'

notify-send -t 3000 "MyBook" "$content"

make build 

./bin/mybook -s
