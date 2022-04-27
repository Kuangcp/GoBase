module keylogger-gui

go 1.16

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gotk3/gotk3 v0.6.1
	github.com/kuangcp/gobase/pkg/cuibase v1.0.6
	github.com/kuangcp/gobase/toolbox/keylogger v0.0.0-20220417170000-8b486ea221d8
	github.com/kuangcp/logger v1.0.8
)

replace github.com/kuangcp/gobase/toolbox/keylogger v0.0.0-20220417170000-8b486ea221d8 => ../keylogger
