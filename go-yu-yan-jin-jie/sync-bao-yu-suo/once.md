# sync.Once

## 什么是sync.Once

`sync.Once` 是 Go 标准库提供的使函数只执行一次的实现。

## 使用场景

像单例模式，初始化配置，保持数据库的连接等都可以使用该函数。有的时候，我们多个 goroutine 都要过一个操作，但是这个操作我只希望被执行一次，这个时候 Once 就上场了。比如下面的例子 :

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    var once sync.Once
    onceBody := func() {
        fmt.Println("Only once")
    }
    for i := 0; i < 10; i++ {
        go func() {
            once.Do(onceBody)
        }()
    }
    time.Sleep(3e9)
}
```

该程序执行后只会打出一次 "Only once"。

## 与init函数的区别

* init 函数是当所在的 package 首次被加载时执行，若迟迟未被使用，则既浪费了内存，又延长了程序加载时间。
* sync.Once 可以在代码的任意位置初始化和调用，因此可以延迟到使用时再执行，并发场景下是线程安全的。

## sync.Once 的原理

> 转载自：[https://geektutu.com/post/hpg-sync-once.html](https://geektutu.com/post/hpg-sync-once.html)

首先：保证变量仅被初始化一次，需要有个标志来判断变量是否已初始化过，若没有则需要初始化。

第二：线程安全，支持并发，无疑需要互斥锁来实现。

### 源码实现

以下是 `sync.Once` 的源码实现，代码位于 `$(dirname $(which go))/../src/sync/once.go`：

```go
package sync

import (
    "sync/atomic"
)

type Once struct {
    done uint32
    m    Mutex
}

func (o *Once) Do(f func()) {
    if atomic.LoadUint32(&o.done) == 0 {
        o.doSlow(f)
    }
}

func (o *Once) doSlow(f func()) {
    o.m.Lock()
    defer o.m.Unlock()
    if o.done == 0 {
        defer atomic.StoreUint32(&o.done, 1)
        f()
    }
}
```

`sync.Once` 的实现与一开始的猜测是一样的，使用 `done` 标记是否已经初始化，使用锁 `m Mutex` 实现线程安全。

### done 为什么是第一个字段

字段 `done` 的注释也非常值得一看：

```go
type Once struct {
    // done indicates whether the action has been performed.
    // It is first in the struct because it is used in the hot path.
    // The hot path is inlined at every call site.
    // Placing done first allows more compact instructions on some architectures (amd64/x86),
    // and fewer instructions (to calculate offset) on other architectures.
    done uint32
    m    Mutex
}
```

其中解释了为什么将 done 置为 Once 的第一个字段：done 在热路径中，done 放在第一个字段，能够减少 CPU 指令，也就是说，这样做能够提升性能。

简单解释下这句话：

1. 热路径\(hot path\)是程序非常频繁执行的一系列指令，sync.Once 绝大部分场景都会访问 `o.done`，在热路径上是比较好理解的，如果 hot path 编译后的机器码指令更少，更直接，必然是能够提升性能的。
2. 为什么放在第一个字段就能够减少指令呢？因为结构体第一个字段的地址和结构体的指针是相同的，如果是第一个字段，直接对结构体的指针解引用即可。如果是其他的字段，除了结构体指针外，还需要计算与第一个值的偏移\(calculate offset\)。在机器码中，偏移量是随指令传递的附加值，CPU 需要做一次偏移值与指针的加法运算，才能获取要访问的值的地址。因为，访问第一个字段的机器代码更紧凑，速度更快。

