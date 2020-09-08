# 按键记录

> 柱状图 每天按键数据
![](https://img-blog.csdnimg.cn/20200908173215731.png)

> 热力图 一周每小时数据
![](https://img-blog.csdnimg.cn/20200908173215775.png)

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
