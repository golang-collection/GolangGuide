# 并发调度模型

Go 语言的并发同步模型来自一个叫作通信顺序进程（Communicating Sequential Processes，CSP） 的范型（paradigm）。

CSP 是一种消息传递模型，通过在 goroutine 之间传递数据来传递消息，而不是对数据进行加锁来实现同步访问。用于在 goroutine 之间同步和传递数据的关键数据类型叫作通道 （channel）。

Go 语言的运行时会在逻辑处理器上调度 goroutine来运行。每个逻辑处理器都分别绑定到单个操作系统线程。

在 1.5 版本 上，Go语言的运行时默认会为每个可用的物理处理器分配一个逻辑处理器。

在 1.5 版本之前的版本中，默认给 整个应用程序只分配一个逻辑处理器。这些逻辑处理器会用于执行所有被创建的goroutine。

即便只有一个逻辑处理器，Go也可以以神奇的效率和性能，并发调度无数个goroutine。

## 推荐阅读

{% embed url="https://draveness.me/system-design-scheduler/" %}



