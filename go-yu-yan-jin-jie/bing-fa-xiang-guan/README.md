# 并发相关

转载自：[http://c.biancheng.net/view/4356.html](http://c.biancheng.net/view/4356.html)

## Goroutine 介绍

goroutine 是一种非常轻量级的实现，可在单个进程里执行成千上万的并发任务，它是Go语言并发设计的核心。

说到底 goroutine 其实就是线程，但是它比线程更小，十几个 goroutine 可能体现在底层就是五六个线程，而且Go语言内部也实现了 goroutine 之间的内存共享。

使用 go 关键字就可以创建 goroutine，将 go 声明放到一个需调用的函数之前，在相同地址空间调用运行这个函数，这样该函数执行时便会作为一个独立的并发线程，这种线程在Go语言中则被称为 goroutine。

goroutine 的用法如下：

```go
//go 关键字放在方法调用前新建一个 goroutine 并执行方法体
go GetThingDone(param1, param2);
//新建一个匿名方法并执行
go func(param1, param2) {
}(val1, val2)
//直接新建一个 goroutine 并在 goroutine 中执行代码块
go {
    //do someting...
}
```

因为 goroutine 在多核 cpu 环境下是并行的，如果代码块在多个 goroutine 中执行，那么我们就实现了代码的并行。

如果需要了解程序的执行情况，怎么拿到并行的结果呢？需要配合使用channel进行。

## channel

channel 是Go语言在语言级别提供的 goroutine 间的通信方式。我们可以使用 channel 在两个或多个 goroutine 之间传递消息。

channel 是进程内的通信方式，因此通过 channel 传递对象的过程和调用函数时的参数传递行为比较一致，比如也可以传递指针等。如果需要跨进程通信，我们建议用分布式系统的方法来解决，比如使用 Socket 或者 HTTP 等通信协议。Go语言对于网络方面也有非常完善的支持。

channel 是类型相关的，也就是说，一个 channel 只能传递一种类型的值，这个类型需要在声明 channel 时指定。如果对 Unix 管道有所了解的话，就不难理解 channel，可以将其认为是一种类型安全的管道。

定义一个 channel 时，也需要定义发送到 channel 的值的类型，注意，必须使用 make 创建 channel，代码如下所示：

```go
ci := make(chan int)
cs := make(chan string)
cf := make(chan interface{})
```

回到在 Windows 和 Linux 出现之前的古老年代，在开发程序时并没有并发的概念，因为命令式程序设计语言是以串行为基础的，程序会顺序执行每一条指令，整个程序只有一个执行上下文，即一个调用栈，一个堆。

并发则意味着程序在运行时有多个执行上下文，对应着多个调用栈。我们知道每一个进程在运行时，都有自己的调用栈和堆，有一个完整的上下文，而操作系统在调度进程的时候，会保存被调度进程的上下文环境，等该进程获得时间片后，再恢复该进程的上下文到系统中。

从整个操作系统层面来说，多个进程是可以并发的，那么并发的价值何在？下面我们先看以下几种场景。

1\) 一方面我们需要灵敏响应的图形用户界面，一方面程序还需要执行大量的运算或者 IO 密集操作，而我们需要让界面响应与运算同时执行。

2\) 当我们的 Web 服务器面对大量用户请求时，需要有更多的“Web 服务器工作单元”来分别响应用户。

3\) 我们的事务处于分布式环境上，相同的工作单元在不同的计算机上处理着被分片的数据，计算机的 CPU 从单内核（core）向多内核发展，而我们的程序都是串行的，计算机硬件的能力没有得到发挥。

4\) 我们的程序因为 IO 操作被阻塞，整个程序处于停滞状态，其他 IO 无关的任务无法执行。

从以上几个例子可以看到，串行程序在很多场景下无法满足我们的要求。下面我们归纳了并发程序的几条优点，让大家认识到并发势在必行：

* 并发能更客观地表现问题模型；
* 并发可以充分利用 CPU 核心的优势，提高程序的执行效率；
* 并发能充分利用 CPU 与其他硬件设备固有的异步性。

## 参考资料

\[1\] [https://mp.weixin.qq.com/s/nTSpQkE6As5YrfO2NevKcw](https://mp.weixin.qq.com/s/nTSpQkE6As5YrfO2NevKcw)

\[2\] [https://mp.weixin.qq.com/s/Tq1Pi9NeOhKvr62S-X2Jdw](https://mp.weixin.qq.com/s/Tq1Pi9NeOhKvr62S-X2Jdw)



