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

其中AddFunc中的参数解释如下

![](../.gitbook/assets/image%20%2834%29.png)

在cron中可以设置的值如下

![](../.gitbook/assets/image%20%2836%29.png)

对于可设置的值解释如下

* `*:`The asterisk indicates that the cron expression will match for all values of the field; e.g., using an asterisk in the 5th field \(month\) would indicate every month.
* `/:`Slashes are used to describe increments of ranges. For example 3-59/15 in the 1st field \(minutes\) would indicate the 3rd minute of the hour and every 15 minutes thereafter. The form "\*\/..." is equivalent to the form "first-last/...", that is, an increment over the largest possible range of the field. The form "N/..." is accepted as meaning "N-MAX/...", that is, starting at N, use the increment until the end of that specific range. It does not wrap around.
* `,:`Commas are used to separate items of a list. For example, using "MON,WED,FRI" in the 5th field \(day of week\) would mean Mondays, Wednesdays and Fridays.
* `-:`Hyphens are used to define ranges. For example, 9-17 would indicate every hour between 9am and 5pm inclusive.
* `?:`Question mark may be used instead of '\*' for leaving either day-of-month or day-of-week blank.

## 预定义的值

![](../.gitbook/assets/image%20%2835%29.png)

![](../.gitbook/assets/image%20%2833%29.png)

## 设置时区

![](../.gitbook/assets/image%20%2837%29.png)

## 推荐资源

{% embed url="https://github.com/robfig/cron" %}

{% embed url="https://godoc.org/github.com/robfig/cron" %}

{% embed url="https://en.wikipedia.org/wiki/Cron" %}

## 教程

{% embed url="https://eddycjy.gitbook.io/golang/di-3-ke-gin/cron" %}



