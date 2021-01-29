# \[\]byte和string互转

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

