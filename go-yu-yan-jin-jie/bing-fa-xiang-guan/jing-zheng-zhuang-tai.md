# 竞争状态

转载自：[http://c.biancheng.net/view/4358.html](http://c.biancheng.net/view/4358.html)

有并发，就有资源竞争，如果两个或者多个 goroutine 在没有相互同步的情况下，访问某个共享的资源，比如同时对该资源进行读写时，就会处于相互竞争的状态，这就是并发中的资源竞争。

 并发本身并不复杂，但是因为有了资源竞争的问题，就使得我们开发出好的并发程序变得复杂起来，因为会引起很多莫名其妙的问题。

 下面的代码中就会出现竞争状态：

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
)

var (
    count int32
    wg    sync.WaitGroup
)

func main() {
    wg.Add(2)
    go incCount()
    go incCount()
    wg.Wait()
    fmt.Println(count)
}

func incCount() {
    defer wg.Done()
    for i := 0; i < 2; i++ {
        value := count
        runtime.Gosched()
        value++
        count = value
    }
}
```

 这是一个资源竞争的例子，大家可以将程序多运行几次，会发现结果可能是 2，也可以是 3，还可能是 4。这是因为 count 变量没有任何同步保护，所以两个 goroutine 都会对其进行读写，会导致对已经计算好的结果被覆盖，以至于产生错误结果。

 代码中的 runtime.Gosched\(\) 是让当前 goroutine 暂停的意思，退回执行队列，让其他等待的 goroutine 运行，目的是为了使资源竞争的结果更明显。

 下面我们来分析一下程序的运行过程，将两个 goroutine 分别假设为 g1 和 g2：

*  g1 读取到 count 的值为 0；
*  然后 g1 暂停了，切换到 g2 运行，g2 读取到 count 的值也为 0；
*  g2 暂停，切换到 g1，g1 对 count+1，count 的值变为 1；
*  g1 暂停，切换到 g2，g2 刚刚已经获取到值 0，对其 +1，最后赋值给 count，其结果还是 1；
*  可以看出 g1 对 count+1 的结果被 g2 给覆盖了，两个 goroutine 都 +1 而结果还是 1。

 通过上面的分析可以看出，之所以出现上面的问题，是因为两个 goroutine 相互覆盖结果。

 所以我们对于同一个资源的读写必须是原子化的，也就是说，同一时间只能允许有一个 goroutine 对共享资源进行读写操作。

 共享资源竞争的问题，非常复杂，并且难以察觉，好在 Go 为我们提供了一个工具帮助我们检查，这个就是`go build -race` 命令。在项目目录下执行这个命令，生成一个可以执行文件，然后再运行这个可执行文件，就可以看到打印出的检测信息。

 在`go build`命令中多加了一个`-race` 标志，这样生成的可执行程序就自带了检测资源竞争的功能，运行生成的可执行文件，效果如下所示：

```text
==================
 WARNING: DATA RACE
 Read at 0x000000619cbc by goroutine 8:
   main.incCount()
       D:/code/src/main.go:25 +0x80
 Previous write at 0x000000619cbc by goroutine 7:
   main.incCount()
       D:/code/src/main.go:28 +0x9f
 Goroutine 8 (running) created at:
   main.main()
       D:/code/src/main.go:17 +0x7e
 Goroutine 7 (finished) created at:
   main.main()
       D:/code/src/main.go:16 +0x66
 ==================
 4
 Found 1 data race(s)
```

 通过运行结果可以看出 goroutine 8 在代码 25 行读取共享资源 `value := count`，而这时 goroutine 7 在代码 28 行修改共享资源 `count = value`，而这两个 goroutine 都是从 main 函数的 16、17 行通过 go 关键字启动的。

##  锁住共享资源

 Go语言提供了传统的同步 goroutine 的机制，就是对共享资源加锁。atomic 和 sync 包里的一些函数就可以对共享的资源进行加锁操作。

###  原子函数

 原子函数能够以很底层的加锁机制来同步访问整型变量和指针，示例代码如下所示：

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "sync/atomic"
)

var (
    counter int64
    wg      sync.WaitGroup
)

func main() {
    wg.Add(2)
    go incCounter(1)
    go incCounter(2)

    wg.Wait() //等待goroutine结束
    fmt.Println(counter)
}

func incCounter(id int) {
    defer wg.Done()
    for count := 0; count < 2; count++ {
        atomic.AddInt64(&counter, 1) //安全的对counter加1

        runtime.Gosched()
    }
}
```

 上述代码中使用了 atmoic 包的 AddInt64 函数，这个函数会同步整型值的加法，方法是强制同一时刻只能有一个 gorountie 运行并完成这个加法操作。当 goroutine 试图去调用任何原子函数时，这些 goroutine 都会自动根据所引用的变量做同步处理。

 另外两个有用的原子函数是 LoadInt64 和 StoreInt64。这两个函数提供了一种安全地读和写一个整型值的方式。下面是代码就使用了 LoadInt64 和 StoreInt64 函数来创建一个同步标志，这个标志可以向程序里多个 goroutine 通知某个特殊状态。

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

var (
    shutdown int64
    wg       sync.WaitGroup
)

func main() {
    wg.Add(2)

    go doWork("A")
    go doWork("B")

    time.Sleep(1 * time.Second)
    fmt.Println("Shutdown Now")
    atomic.StoreInt64(&shutdown, 1)
    wg.Wait()
}

func doWork(name string) {
    defer wg.Done()

    for {
        fmt.Printf("Doing %s Work\n", name)
        time.Sleep(250 * time.Millisecond)

        if atomic.LoadInt64(&shutdown) == 1 {
            fmt.Printf("Shutting %s Down\n", name)
            break
        }
    }
}
```

 上面代码中 main 函数使用 StoreInt64 函数来安全地修改 shutdown 变量的值。如果哪个 doWork goroutine 试图在 main 函数调用 StoreInt64 的同时调用 LoadInt64 函数，那么原子函数会将这些调用互相同步，保证这些操作都是安全的，不会进入竞争状态。

###  互斥锁

 另一种同步访问共享资源的方式是使用互斥锁，互斥锁这个名字来自互斥的概念。互斥锁用于在代码上创建一个临界区，保证同一时间只有一个 goroutine 可以执行这个临界代码。

 示例代码如下所示：

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
)

var (
    counter int64
    wg      sync.WaitGroup
    mutex   sync.Mutex
)

func main() {
    wg.Add(2)

    go incCounter(1)
    go incCounter(2)

    wg.Wait()
    fmt.Println(counter)
}

func incCounter(id int) {
    defer wg.Done()

    for count := 0; count < 2; count++ {
        //同一时刻只允许一个goroutine进入这个临界区
        mutex.Lock()
        {
            value := counter
            runtime.Gosched()
            value++
            counter = value
        }
        mutex.Unlock() //释放锁，允许其他正在等待的goroutine进入临界区
    }
}
```

 同一时刻只有一个 goroutine 可以进入临界区。之后直到调用 Unlock 函数之后，其他 goroutine 才能进去临界区。当调用 runtime.Gosched 函数强制将当前 goroutine 退出当前线程后，调度器会再次分配这个 goroutine 继续运行。

