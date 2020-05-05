# 按键记录

- string keyboard:last-event
- `keyboard:total` **zset**
	- value:年:月:日 score: 当日按键总数

- `keyboard:年:月:日:detail` **zset**
	- value: 时间戳(微秒) score: keyCode
- `keyboard:年:月:日:rank` **zset**
	- value: keyCode score: 当日按键数
