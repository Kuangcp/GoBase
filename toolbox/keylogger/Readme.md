# 按键记录

> `柱状图` 每天按键数据

![](https://img-blog.csdnimg.cn/20200908173215731.png)

> `热力图` 时间段内分布数据

![](https://img-blog.csdnimg.cn/20200908173215775.png)

> `热力图` 多个星期横向对比

![](https://img-blog.csdnimg.cn/20200912222920568.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2tjcDYwNg==,size_16,color_FFFFFF,t_70#pic_center)

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

## Version 
- 1.0.3 import thread pool 
- 1.0.2 GA
