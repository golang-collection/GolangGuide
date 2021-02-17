# defer延迟调用

首发于：[https://blog.csdn.net/s\_842499467/article/details/104283109](https://blog.csdn.net/s_842499467/article/details/104283109)

## 1. 为什么需要defer机制

在项目中，我们常用的操作就是释放资源等操作，为了即时的释放资源，go设计了`defer`机制。defer语句常用语资源释放，解除锁定以及错误处理等。

## 2. defer机制的简单实用

我们先通过一个案例来看一下`defer`机制的作用

```go
import "fmt"

func defer_sum(num1, num2 int) int {
	defer fmt.Println("defer1 : ", num1)
	defer fmt.Println("defer2 : ", num2)

	res := num1 + num2
	fmt.Println("sum res : ", res)
	return res
}

func main() {
	res := defer_sum(10 ,20)
	fmt.Println("main : ", res)
}
```

程序的输出结果如下

```go
sum res :  30
defer2 :  20
defer1 :  10
main :  30
```

我们可以看到通过`defer`声明的操作不会立即执行，而是当所处的代码块执行完成之后才继续执行。

## 3. defer机制的细节

那么defer机制为什么会产生这样的结果呢，是因为通过defer声明的语句不会立刻执行，**而会压入一个`defer`栈中**，因为栈遵循的是先入后出原则，所以得到如上所示的结果。

**在 `defer`将语句放入到栈时，也会将相关的值拷贝同时入栈。** 我们用一段代码来解释这个特点。这段代码与刚才的那段代码没什么不同，只是在中间加入了一段`num1`和`num2`自增和自减的操作。

```go
import "fmt"

func defer_sum(num1, num2 int) int {
	defer fmt.Println("defer1 : ", num1)
	defer fmt.Println("defer2 : ", num2)
	
	num1++
	num2--

	res := num1 + num2
	fmt.Println("sum res : ", res)
	return res
}

func main() {
	res := defer_sum(10 ,20)
	fmt.Println("main : ", res)
}
```

程序的输出结果如下

```text
sum res :  30
defer2 :  20
defer1 :  10
main :  30
```

可以看到输出结果与上一个案例的输出结果相同，这也就解释了在**`defer`语句入栈的同时也会将值拷贝一份压入栈中。**

## 4. 什么时候使用defer

就像开头提到的，defer机制就是为了更好的关闭资源的，所以我们使用`defer`也是在创建资源后使用，如下例所示。

```go
func main(){
    connect := connectDB()
    defer connect.close()
    
}
```

## 5. 注意事项

需要注意一点的是如果我们在main函数中申请资源时使用了defer，要注意这个资源是main函数执行完才会被释放，如果申请的资源很大那无疑是一种错误的处理方式，更为优雅的方式是将其与匿名函数封装，这样匿名函数调用结束则释放资源。

而且defer语句也要花费更大的代价，所以在高性能的算法设计中要谨慎使用。

## 6. 底层实现

defer函数对应的结构体如下，defer函数会注册到一个链表中，每个goroutine都会持有该链表的头指针，新注册的defer函数会添加到链表的头部，所以defer函数执行起来是倒序执行的。

```go
// A _defer holds an entry on the list of deferred calls.
// If you add a field here, add code to clear it in freedefer and deferProcStack
// This struct must match the code in cmd/compile/internal/gc/reflect.go:deferstruct
// and cmd/compile/internal/gc/ssa.go:(*state).call.
// Some defers will be allocated on the stack and some on the heap.
// All defers are logically part of the stack, so write barriers to
// initialize them are not required. All defers must be manually scanned,
// and for heap defers, marked.
type _defer struct {
    siz     int32 // includes both arguments and results
    started bool //defer是否已经执行
    heap    bool //标识当前_defer是否是堆分配
    // openDefer indicates that this _defer is for a frame with open-coded
    // defers. We have only one defer record for the entire frame (which may
    // currently have 0, 1, or more defers active).
    openDefer bool //1.14版本新增加内容，用于判断当前_defer是否需要栈扫描才能执行
    sp        uintptr  // sp at time of defer 调用者栈指针
    pc        uintptr  // pc at time of defer 返回地址
    fn        *funcval // can be nil for open-coded defers 要注册的funcval
    _panic    *_panic  // panic that is running defer
    link      *_defer //对应的defer链表

    // If openDefer is true, the fields below record values about the stack
    // frame and associated function that has the open-coded defer(s). sp
    // above will be the sp for the frame, and pc will be address of the
    // deferreturn call in the function.
    fd   unsafe.Pointer // funcdata for the function associated with the frame
    varp uintptr        // value of varp for the stack frame
    // framepc is the current pc associated with the stack frame. Together,
    // with sp above (which is the sp associated with the stack frame),
    // framepc/sp can be used as pc/sp pair to continue a stack trace via
    // gentraceback().
    framepc uintptr
}
```

在生命一个defer函数时，会调用下面这样一个函数

```go
func deferproc(siz int32, fn *funcval)
```

注册完成后go会在对中为defer函数开辟对应的\_derfer结构体空间，用来存储其对应的值，但是在go中会预先分配一个defer池，在注册defer函数时会优先从defer池中取出defer对象，没有空闲的或者大小合适的才会进行堆分配。

上面介绍的是在1.12版本中defer的执行逻辑，该方式有两个问题

* \_defer结构体是在堆上进行分配，而且还存在变量在栈与堆上的变换，比较复杂
* \_defer是使用链表进行保存，操作起来不方便

在1.13版本中将defer函数的内容保存到当前函数栈帧的局部栈空间中，然后再通过下面这个函数将其注册到defer链表中

```go
func deferprocStack(d *_defer)
```

但是像一些循环中的defer仍然需要1.12版本中的处理方式，所以在\_defer结构体中增加了heap字段，用于标记其是否在堆上分配。

在1.14版本中直接在编译过程中将defer函数修改为直接调用的方式，如下图所示

![](../../.gitbook/assets/image%20%2859%29.png)

通过df变量来标志某个defer函数是否需要被执行，在1.14中就是通过在编译阶段将defer展开到所属函数的函数体内，从而提高执行效率，称这种方式为open coded defer。但是像循环中的defer仍然需要通过1.12的方式来执行。

但是如果程序执行过程中出现panic等问题，那么一些defer函数就没有办法执行了，所以需要通过栈扫描的方式来使defer函数正常执行。

## 推荐阅读

{% embed url="https://www.kancloud.cn/aceld/golang/1958310" %}

{% embed url="https://www.bilibili.com/video/BV1E5411x7NC" %}

{% embed url="https://www.bilibili.com/video/BV1b5411W7ih" %}



