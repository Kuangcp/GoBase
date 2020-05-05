# 按键记录

## Redis
> 全局
- **string** `keyboard:last-event`
- **zset** `keyboard:total`
	- value: 年:月:日 score: 当日按键总数

> 每日滚动数据
- **zset** `keyboard:年:月:日:detail`
	- value: 时间戳(微秒) score: keyCode
- **zset** `keyboard:年:月:日:rank`
	- value: keyCode score: 当日按键数
