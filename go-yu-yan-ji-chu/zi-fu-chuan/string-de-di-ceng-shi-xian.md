# string的底层实现

字符串底层结构如下所示

```text
type stringStruct struct {
	str unsafe.Pointer
	len int
}
```

其中str指向存储字符编码的底层数组（使用变长编码来表示要表示的字符）。且str与len各占8个字节。

## 推荐资源

{% embed url="https://www.bilibili.com/video/BV1ff4y1m72A" %}



