# cron

cron是一个用于设置定时任务的包，使用方式如下：

```go
c := cron.New()
c.AddFunc("30 * * * *", func() { fmt.Println("Every hour on the half hour") })
c.AddFunc("30 3-6,20-23 * * *", func() { fmt.Println(".. in the range 3-6am, 8-11pm") })
c.AddFunc("CRON_TZ=Asia/Tokyo 30 04 * * *", func() { fmt.Println("Runs at 04:30 Tokyo time every day") })
c.AddFunc("@hourly",      func() { fmt.Println("Every hour, starting an hour from now") })
c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty, starting an hour thirty from now") })
c.Start()
```



cron包用于设置定时任务

{% embed url="https://github.com/robfig/cron" %}

文档位置

{% embed url="https://godoc.org/github.com/robfig/cron" %}

cron表达式

{% embed url="https://en.wikipedia.org/wiki/Cron" %}





## 教程

{% embed url="https://eddycjy.gitbook.io/golang/di-3-ke-gin/cron" %}



