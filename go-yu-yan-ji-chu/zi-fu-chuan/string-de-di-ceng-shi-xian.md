# string的底层实现

字符串底层结构如下所示

```go
type stringStruct struct {
	str unsafe.Pointer
	len int
}
```

其中str指向存储字符编码的底层数组（使用变长编码来表示要表示的字符）。且str与len各占8个字节。

## string与\[\]byte类型转换 

> 转载自：[https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-string/](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-string/)

当我们使用 Go 语言解析和序列化 JSON 等数据格式时，经常需要将数据在 `string` 和 `[]byte` 之间来回转换，类型转换的开销并没有想象的那么小，我们经常会看到 [`runtime.slicebytetostring`](https://draveness.me/golang/tree/runtime.slicebytetostring) 等函数出现在火焰图[1](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-string/#fn:1)中，成为程序的性能热点。

从字节数组到字符串的转换需要使用 [`runtime.slicebytetostring`](https://draveness.me/golang/tree/runtime.slicebytetostring) 函数，例如：`string(bytes)`，该函数在函数体中会先处理两种比较常见的情况，也就是长度为 0 或者 1 的字节数组，这两个情况处理起来都非常简单：

```go
func slicebytetostring(buf *tmpBuf, b []byte) (str string) {
	l := len(b)
	if l == 0 {
		return ""
	}
	if l == 1 {
		stringStructOf(&str).str = unsafe.Pointer(&staticbytes[b[0]])
		stringStructOf(&str).len = 1
		return
	}
	var p unsafe.Pointer
	if buf != nil && len(b) <= len(buf) {
		p = unsafe.Pointer(buf)
	} else {
		p = mallocgc(uintptr(len(b)), nil, false)
	}
	stringStructOf(&str).str = p
	stringStructOf(&str).len = len(b)
	memmove(p, (*(*slice)(unsafe.Pointer(&b))).array, uintptr(len(b)))
	return
}
```

处理过后会根据传入的缓冲区大小决定是否需要为新字符串分配一片内存空间，[`runtime.stringStructOf`](https://draveness.me/golang/tree/runtime.stringStructOf) 会将传入的字符串指针转换成 [`runtime.stringStruct`](https://draveness.me/golang/tree/runtime.stringStruct) 结构体指针，然后设置结构体持有的字符串指针 `str` 和长度 `len`，最后通过 [`runtime.memmove`](https://draveness.me/golang/tree/runtime.memmove) 将原 `[]byte` 中的字节全部复制到新的内存空间中。

当我们想要将字符串转换成 `[]byte` 类型时，需要使用 [`runtime.stringtoslicebyte`](https://draveness.me/golang/tree/runtime.stringtoslicebyte) 函数，该函数的实现非常容易理解：

```go
func stringtoslicebyte(buf *tmpBuf, s string) []byte {
	var b []byte
	if buf != nil && len(s) <= len(buf) {
		*buf = tmpBuf{}
		b = buf[:len(s)]
	} else {
		b = rawbyteslice(len(s))
	}
	copy(b, s)
	return b
}
```

上述函数会根据是否传入缓冲区做出不同的处理：

* 当传入缓冲区时，它会使用传入的缓冲区存储 `[]byte`；
* 当没有传入缓冲区时，运行时会调用 [`runtime.rawbyteslice`](https://draveness.me/golang/tree/runtime.rawbyteslice) 创建新的字节切片并将字符串中的内容拷贝过去；

字符串和 `[]byte` 中的内容虽然一样，但是字符串的内容是只读的，我们不能通过下标或者其他形式改变其中的数据，而 `[]byte` 中的内容是可以读写的。不过无论从哪种类型转换到另一种都需要拷贝数据，而内存拷贝的性能损耗会随着字符串和 `[]byte` 长度的增长而增长。

> 转载自：[https://www.cnblogs.com/qcrao-2018/p/10964692.html](https://www.cnblogs.com/qcrao-2018/p/10964692.html)

完成这个任务，我们需要了解 slice 和 string 的底层数据结构：

```go
type StringHeader struct {
	Data uintptr
	Len  int
}

type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
```

上面是反射包下的结构体，路径：src/reflect/value.go。只需要共享底层 \[\]byte 数组就可以实现 `zero-copy`。

```go
func string2bytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))

	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}

	return *(*[]byte)(unsafe.Pointer(&bh))
}

func bytes2string(b []byte) string{
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}

	return *(*string)(unsafe.Pointer(&sh))
}
```

## 推荐资源

{% embed url="https://www.bilibili.com/video/BV1ff4y1m72A" %}

{% embed url="https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-string/" %}

