# string的底层实现

字符串底层结构如下所示

```text
type StringHeader struct {
	Data uintptr
	Len  int
}
```

其中Data使用变长编码来表示要表示的字符。

## 推荐资源

{% embed url="https://www.bilibili.com/video/BV1ff4y1m72A" %}



