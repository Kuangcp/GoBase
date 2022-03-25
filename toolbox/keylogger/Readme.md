# 按键记录
> [gradient.shapefactory.co](https://gradient.shapefactory.co)

![](https://img-blog.csdnimg.cn/20201012105207695.png)

************************

> `bar chart` every day

![](https://img-blog.csdnimg.cn/20200908173215731.png)

************************

> `heatmap chart` every hour in weeks

************************

![](https://img-blog.csdnimg.cn/20200908173215775.png)

> `heatmap chart` comparison of several weeks

![](https://img-blog.csdnimg.cn/20200912222920568.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2tjcDYwNg==,size_16,color_FFFFFF,t_70#pic_center)


> `柱状图` 每天按键数据

![](https://img-blog.csdnimg.cn/20200908173215731.png)

> `热力图` 时间段内分布数据

![](https://img-blog.csdnimg.cn/20200908173215775.png)

> `热力图` 多个星期横向对比

![](https://img-blog.csdnimg.cn/20200912222920568.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2tjcDYwNg==,size_16,color_FFFFFF,t_70#pic_center)

## Install 
1. make down 
1. make statik
1. make install 
1. make web

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
- 1.1.0 fix memory leak
- 1.0.9 fix 00：00 cache problem
- 1.0.8 transfer the task of calculating KPM to the listening input device process
- 1.0.6 add interactive select device
- 1.0.4 remove thread pool(memory leak)
- 1.0.3 import thread pool 
- 1.0.2 GA

## Debug
> `go tool pprof -inuse_space http://localhost:8891/debug/pprof/heap`
> `go tool pprof -inuse_space -cum -svg http://localhost:8891/debug/pprof/heap > heap_inuse.svg`

## TODO
- [x] gtk 窗口有内存泄漏的问题，随着刷新次数的增多，内存也随之增长
    - 使用timeoutAdd 解决
1. [webview vs electron](https://www.zhihu.com/question/396199869)
