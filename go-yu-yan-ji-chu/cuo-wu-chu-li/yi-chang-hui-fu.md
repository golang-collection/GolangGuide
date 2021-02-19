# 异常恢复

转载自：[https://github.com/unknwon/the-way-to-go\_ZH\_CN/blob/master/eBook/13.3.md](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/13.3.md) by unknwon

## 简介

正如名字一样，这个（recover）内建函数被用于从 panic 或 错误场景中恢复：让程序可以从 panicking 重新获得控制权，停止终止过程进而恢复正常执行。

`recover` 只能在 defer 修饰的函数中使用：用于取得 panic 调用中传递过来的错误值，如果是正常执行，调用 `recover` 会返回 nil，且没有其它效果。

下面例子中的 protect 函数调用函数参数 g 来保护调用者防止从 g 中抛出的运行时 panic，并展示 panic 中的信息：

```go
func protect(g func()) {
	defer func() {
		log.Println("done")
		// Println executes normally even if there is a panic
		if err := recover(); err != nil {
			log.Printf("run time panic: %v", err)
		}
	}()
	log.Println("start")
	g() //   possible runtime-error
}
```

这跟 Java 和 .NET 这样的语言中的 catch 块类似。

log 包实现了简单的日志功能：默认的 log 对象向标准错误输出中写入并打印每条日志信息的日期和时间。除了 `Println` 和 `Printf` 函数，其它的致命性函数都会在写完日志信息后调用 os.Exit\(1\)，那些退出函数也是如此。而 Panic 效果的函数会在写完日志信息后调用 panic；可以在程序必须中止或发生了临界错误时使用它们，就像当 web 服务器不能启动时那样。

{% hint style="info" %}
但是需要注意的是，不能所有错误都一味的recover，很有可能会导致程序假死，已经挂掉了但是心跳还在的情况，倒不如直接通过守护线程直接重启的方式更有效。
{% endhint %}

## 案例

这是一个展示 panic，defer 和 recover 怎么结合使用的完整例子：

```go
// panic_recover.go
package main

import (
	"fmt"
)

func badCall() {
	panic("bad end")
}

func test() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("Panicing %s\r\n", e)
		}
	}()
	badCall()
	fmt.Printf("After bad call\r\n") // <-- wordt niet bereikt
}

func main() {
	fmt.Printf("Calling test\r\n")
	test()
	fmt.Printf("Test completed\r\n")
}
```

输出：

```go
Calling test
Panicing bad end
Test completed
```

`defer-panic-recover` 在某种意义上也是一种像 `if`，`for` 这样的控制流机制。

Go 标准库中许多地方都用了这个机制，例如，json 包中的解码和 regexp 包中的 Complie 函数。Go 库的原则是即使在包的内部使用了 panic，在它的对外接口（API）中也必须用 recover 处理成返回显式的错误。

## 底层实现

在每个goroutine结构体中，有一个\_panic字段用于保存panic的链表，与\_defer相同也是采用采用头插的方式保存panic信息，panic对应的结构体底层如下

```go
type _panic struct {
    argp      unsafe.Pointer // pointer to arguments of deferred call run during panic; cannot move - known to liblink
    arg       interface{}    // argument to panic
    link      *_panic        // link to earlier panic
    pc        uintptr        // where to return to in runtime if this panic is bypassed
    sp        unsafe.Pointer // where to return to in runtime if this panic is bypassed
    recovered bool           // whether this panic is over
    aborted   bool           // the panic was aborted
    goexit    bool
}
```

下图是一个简单的例子

![](../../.gitbook/assets/image%20%2861%29.png)

当发生panic时，会触发对应的defer函数，但与正常的defer执行不同，会将defer结构体内的started标记为true，代表当前defer函数开始执行，并将其\_panic指针指向panic结构体。如果A2正常结束，那么就会将A2从defer链表移除，继续执行下一个defer函数，这样主要是为了防止defer函数没有正常退出的情况。

![](../../.gitbook/assets/image%20%2863%29.png)

在上图中A1函数继续发生了panic，那么会讲A1添加到panic链表中，然后执行defer链表，但是发现当前A1指向的panic为A，那么就会将panicA结构体中的aborted设置为true。然后将A1从defer链表中移除，现在defer链表为空，但是panic链表不为空，所以执行打印panic信息，会从panic链表尾部开始逐个向前输出。

下面来看一下recover的处理逻辑

![](../../.gitbook/assets/image%20%2862%29.png)

当执行函数A2时，内部调用了recover函数，所以会将panicA中的recovered字段设置为true，然后就将panicA从链表中移除，也将deferA2函数从链表中移除，不过再移除之前会保存sp和pc的寄存器值，来跳出panicA的流程。

![](../../.gitbook/assets/image%20%2860%29.png)

但是在上图中如果函数A2在执行时又发生了panic，就会将panicA2添加到panic链表中，然后去执行defer函数，但是发现A2已经由之前的panicA执行了，所以就会将panicA的aborted字段设置为true，然后就将A2函数移除，继续执行A1，A1结束后就输出panic信息，但是会在输出panicA时标记为recovered。

