# 单向通道

转载自：[http://c.biancheng.net/view/99.html](http://c.biancheng.net/view/99.html)

Go语言的类型系统提供了单方向的 channel 类型，顾名思义，单向 channel 就是只能用于写入或者只能用于读取数据。当然 channel 本身必然是同时支持读写的，否则根本没法用。

 假如一个 channel 真的只能读取数据，那么它肯定只会是空的，因为你没机会往里面写数据。同理，如果一个 channel 只允许写入数据，即使写进去了，也没有丝毫意义，因为没有办法读取到里面的数据。所谓的单向 channel 概念，其实只是对 channel 的一种使用限制。

##  单向通道的声明格式

 我们在将一个 channel 变量传递到一个函数时，可以通过将其指定为单向 channel 变量，从而限制该函数中可以对此 channel 的操作，比如只能往这个 channel 中写入数据，或者只能从这个 channel 读取数据。

 单向 channel 变量的声明非常简单，只能写入数据的通道类型为`chan<-`，只能读取数据的通道类型为`<-chan`，格式如下：

 var 通道实例 chan&lt;- 元素类型    // 只能写入数据的通道  
 var 通道实例 &lt;-chan 元素类型    // 只能读取数据的通道

*  元素类型：通道包含的元素类型。
*  通道实例：声明的通道变量。

##  单向通道的使用例子

 示例代码如下：

```go
ch := make(chan int)
// 声明一个只能写入数据的通道类型, 并赋值为ch
var chSendOnly chan<- int = ch
//声明一个只能读取数据的通道类型, 并赋值为ch
var chRecvOnly <-chan int = ch
```

 上面的例子中，chSendOnly 只能写入数据，如果尝试读取数据，将会出现如下报错：

 invalid operation: &lt;-chSendOnly \(receive from send-only type chan&lt;- int\)

 同理，chRecvOnly 也是不能写入数据的。

 当然，使用 make 创建通道时，也可以创建一个只写入或只读取的通道：

```go
ch := make(<-chan int)

var chReadOnly <-chan int = ch
<-chReadOnly
```

 上面代码编译正常，运行也是正确的。但是，一个不能写入数据只能读取的通道是毫无意义的。

##  time包中的单向通道

 time 包中的计时器会返回一个 timer 实例，代码如下：

```go
timer := time.NewTimer(time.Second)
```

 timer的Timer类型定义如下：

```go
type Timer struct {
    C <-chan Time
    r runtimeTimer
}
```

 第 2 行中 C 通道的类型就是一种只能读取的单向通道。如果此处不进行通道方向约束，一旦外部向通道写入数据，将会造成其他使用到计时器的地方逻辑产生混乱。

 因此，单向通道有利于代码接口的严谨性。

##  关闭 channel

 关闭 channel 非常简单，直接使用Go语言内置的 close\(\) 函数即可：

 close\(ch\)

 在介绍了如何关闭 channel 之后，我们就多了一个问题：如何判断一个 channel 是否已经被关闭？我们可以在读取的时候使用多重返回值的方式：

```go
x, ok := <-ch
```

 这个用法与 map 中的按键获取 value 的过程比较类似，只需要看第二个 bool 返回值即可，如果返回值是 false 则表示 ch 已经被关闭。

