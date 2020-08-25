type notify-send >/dev/null 2>&1 || { echo >&2 "notify-send not installed.  Aborting."; exit 1; }

make buildSingle 

if [ ! -d bin/data ]; then 
    ln -s data bin/data
fi 

content='<span color="#57dafd" font="26px"> <a href="http://localhost:9090/static/">Enter</a></span>'

notify-send -t 3000 "MyBook" "$content"

bin/mybook -s -p 9090
