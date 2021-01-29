# unsafe

> 转载自：[https://github.com/polaris1119/The-Golang-Standard-Library-by-Example/blob/master/chapter15/15.02.md](https://github.com/polaris1119/The-Golang-Standard-Library-by-Example/blob/master/chapter15/15.02.md)

_unsafe_库徘徊在“类型安全”边缘，由于它们绕过了 Golang 的内存安全原则，一般被认为使用该库是不安全的。但是，在许多情况下，_unsafe_库的作用又是不可替代的，灵活地使用它们可以实现对内存的直接读写操作。在_reflect_库、_syscall_库以及其他许多需要操作内存的开源项目中都有对它的引用。

## Pointer

这个类型比较重要，它是实现定位欲读写的内存的基础。官方文档对该类型有四个重要描述：

* （1）任何类型的指针都可以被转化为 Pointer
* （2）Pointer 可以被转化为任何类型的指针
* （3）uintptr 可以被转化为 Pointer
* （4）Pointer 可以被转化为 uintptr

举例来说，该类型可以这样使用：

```go
    func main() {
        i := 100
        fmt.Println(i)  // 100
        p := (*int)(unsafe.Pointer(&i))
        fmt.Println(*p) // 100
        *p = 0
        fmt.Println(i)  // 0
        fmt.Println(*p) // 0
    }
```

## 方法

### Sizeof

该函数的定义如下：

```text
func Sizeof(v ArbitraryType) uintptr
```

Sizeof 函数返回变量 v 占用的内存空间的字节数，该字节数不是按照变量 v 实际占用的内存计算，而是按照 v 的“ top level ”内存计算。

* 比如，在 64 位系统中，如果变量 v 是 int 类型，会返回 16，因为 v 的“ top level ”内存就是它的值使用的内存；
* 如果变量 v 是 string 类型，会返回 16，因为 v 的“ top level ”内存不是存放着实际的字符串，而是该字符串的地址；
* 如果变量 v 是 slice 类型，会返回 24，这是因为 slice 的描述符就占了 24 个字节。

### Offsetof

该函数的定义如下：

```go
func Offsetof(v ArbitraryType) uintptr
```

该函数返回由 v 所指示的某结构体中的字段在该结构体中的位置偏移字节数，注意，v 的表达方式必须是“ struct.filed ”形式。 举例说明，在 64 为系统中运行以下代码：

```go
    type Datas struct{
        c0 byte
        c1 int
        c2 string
        c3 int
    }
    func main(){
        var d Datas
        fmt.Println(unsafe.Offset(d.c0))    // 0
        fmt.Println(unsafe.Offset(d.c1))    // 8
        fmt.Println(unsafe.Offset(d.c2))    // 16
        fmt.Println(unsafe.Offset(d.c3))    // 32
    }
```

如果知道的结构体的起始地址和字段的偏移值，就可以直接读写内存：

```go
d.c3 = 13
p := unsafe.Pointer(&d)
offset := unsafe.Offsetof(d.c3)
q := (*int)(unsafe.Pointer(uintptr(p) + offset))
fmt.Println(*q) // 13
*p = 1013
fmt.Println(d.c3)   // 1013
```

## 注意事项

> 转载自：[https://eddycjy.com/posts/go/pkg/2018-12-15-unsafe/\#%E9%94%99%E8%AF%AF%E7%A4%BA%E4%BE%8B-1](https://eddycjy.com/posts/go/pkg/2018-12-15-unsafe/#%E9%94%99%E8%AF%AF%E7%A4%BA%E4%BE%8B-1)

```go
func main(){
	n := Num{i: "EDDYCJY", j: 1}
	nPointer := unsafe.Pointer(&n)
    ...

	ptr := uintptr(nPointer)
	njPointer := (*int64)(unsafe.Pointer(ptr + unsafe.Offsetof(n.j)))
	...
}
```

这里存在一个问题，`uintptr` 类型是不能存储在临时变量中的。因为从 GC 的角度来看，`uintptr` 类型的临时变量只是一个无符号整数，并不知道它是一个指针地址

因此当满足一定条件后，`ptr` 这个临时变量是可能被垃圾回收掉的，那么接下来的内存操作就会产生不可预期的问题。

## 推荐阅读

{% embed url="https://github.com/polaris1119/The-Golang-Standard-Library-by-Example/blob/master/chapter15/15.02.md" %}

{% embed url="https://www.cnblogs.com/qcrao-2018/p/10964692.html" %}

{% embed url="https://www.flysnow.org/2017/07/06/go-in-action-unsafe-pointer.html" %}

{% embed url="https://time.geekbang.org/column/article/18042" %}

{% embed url="https://eddycjy.com/posts/go/pkg/2018-12-15-unsafe/" %}





