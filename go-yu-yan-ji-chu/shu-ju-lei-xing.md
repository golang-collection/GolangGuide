# 数据类型

在golang中默认定义了多种基础数据类型

## 基本类型类型

| 类型名称 | 长度 | 默认值 | 取值范围 |
| :--- | :--- | :--- | :--- |
| 布尔型 bool | 1字节 | false | true, false |
| 整型 int/uint | 32或64位 | 0 |  |
| 8位整型 int8/uint8 | 1字节 | 0 | -128～127/0～255 |
| 字节型 byte | 注：是uint8的别名 | 0 |  |
| 16位整型 int16/uint16 | 2字节 | 0 | -32768～32767/0～65535 |
| 32位整型 int32/uint32 | 4字节 | 0 | $$-\frac{2^{32}}2$$～$$\frac{2^{32}}{2}{-1}$$/$$0$$~$${2^{32}}{-1}$$ |
| 64位整型 int64/uint64 | 8字节 | 0 | $$-\frac{2^{64}}2$$～$$\frac{2^{64}}{2}{-1}$$/$$0$$~$${2^{64}}{-1}$$ |
| 浮点型 float32/float64 | 4/8字节 | 0.0 | 精确到7/15位小数 |
| 复数 complex64/complex128 | 8/16字节 |  |  |
| rune | 4字节 | 0 | Unicode Code Point |
| uintptr | 4/8字节 | 0 |  |
| string |  | "" |  |
| array |  |  |  |
| struct |  |  |  |

字符串string较为特殊，我们放在后面单独讲解一下[字符串](zi-fu-chuan/)。

```go
var a bool = true
var b int = 1
var c int8 = 127
var d byte = 'a'
var e int16 = 32767
var f int32
var g int64
var h float32 = 2.5
var i float64
var j complex64 = complex(3,4) //32位实数和虚数
real := real(j) //实部 3
imag := imag(j)//虚部 4
var h complex128
```

## 引用类型

| 类型名称 | 默认值 | 说明 |
| :--- | :--- | :--- |
| pointer | nil | 指针 |
| map | nil | 字典 |
| slice | nil | 切片 |
| channel | nil | 通道 |
| func | nil | 函数 |
| interface | nil | 接口 |

在这里判断是否为引用类型是根据这一类型的变量通过参数传递后，是值拷贝还是地址拷贝。

而且go语言可以把一个函数当作一种变量类型，这里需要特别注意下。

## 类型转换

在golang中隐式转换造成的问题会大于它带来的好处。如果过度使用会在程序中埋下很多莫名其妙的错误。

```go
package main

import "fmt"

func main() { 
    var sum int = 17 
    var count int = 5 
    var mean float32
    mean = float32(sum)/float32(count) 
    fmt.Printf("mean 的值为: %f\n",mean) 
}
```

```go
mean 的值为: 3.400000
```

## new与make

转载自：[https://www.flysnow.org/2017/10/23/go-new-vs-make.html](https://www.flysnow.org/2017/10/23/go-new-vs-make.html)

### new

在go语言中，new可以为变量申请一块内存，分配好内存后，返回一个指向该类型内存地址的指针。同时请注意它同时把分配的内存置为零，也就是类型的零值。

就像下面这段代码

```go
func main() {
	u := new(user)
	u.lock.Lock()
	u.name = "张三"
	u.lock.Unlock()

	fmt.Println(u)

}

type user struct {
	lock sync.Mutex
	name string
	age int
}
```

示例中的`user`类型中的`lock`字段我不用初始化，直接可以拿来用，不会有无效内存引用异常，因为它已经被零值了。

这就是`new`，它返回的永远是类型的指针，指向分配类型的内存地址。

### make

`make`也是用于内存分配的，但是和`new`不同，它只用于`chan`、`map`以及切片的内存创建，而且它返回的类型就是这三个类型本身，而不是他们的指针类型，因为这三种类型就是引用类型，所以就没有必要返回他们的指针了。

注意，因为这三种类型是引用类型，所以必须得初始化，但是不是置为零值，这个和`new`是不一样的。

```go
func make(t Type, size ...IntegerType) Type
```

从函数声明中可以看到，返回的还是该类型。

### 两者异同

二者都是内存的分配（堆上），但是`make`只用于slice、map以及channel的初始化（非零值）；而`new`用于类型的内存分配，并且内存置为零。所以在我们编写程序的时候，就可以根据自己的需要很好的选择了。

`make`返回的还是这三个引用类型本身；而`new`返回的是指向类型的指针。

但是如果我们使用new来为map这样的引用类型的变量申请空间的话，他只会申请一块只想map类型的指针空间，而不会分配键值对的内存。

```go
func main() {
	p := new(map[string]int)
	m := *p
	m["a"] = 1
	fmt.Println(m)
}
```

```go
panic: assignment to entry in nil map
```

## 自定义类型

在go语言中我们可以使用type关键字来定义我们自己的类型。

```go
type MyInt int
```

