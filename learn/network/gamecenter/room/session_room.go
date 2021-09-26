package room

import (
	"gamecenter/ws"
	"time"
)

type Player struct {
	id      string
	name    string
	session *ws.ServerSession
}

type MultiRoom struct {
	id         string
	name       string
	createTime time.Time
	players    []Player
}
