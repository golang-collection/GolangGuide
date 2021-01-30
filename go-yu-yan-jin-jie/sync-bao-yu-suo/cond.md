# Cond

sync.Cond 是用来控制某个条件下，goroutine 进入等待时期，等待信号到来，然后重新启动。比如：

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    locker := new(sync.Mutex)
    cond := sync.NewCond(locker)
    done := false

    cond.L.Lock()

    go func() {
        time.Sleep(2e9)
        done = true
        cond.Signal()
    }()

    if (!done) {
        cond.Wait()
    }

    fmt.Println("now done is ", done);
}
```

这里当主 goroutine 进入 cond.Wait 的时候，就会进入等待，当从 goroutine 发出信号之后，主 goroutine 才会继续往下面走。

sync.Cond 还有一个 BroadCast 方法，用来通知唤醒所有等待的 gouroutine。

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

var locker = new(sync.Mutex)
var cond = sync.NewCond(locker)

func test(x int) {

    cond.L.Lock() // 获取锁
    cond.Wait()   // 等待通知  暂时阻塞
    fmt.Println(x)
    time.Sleep(time.Second * 1)
    cond.L.Unlock() // 释放锁，不释放的话将只会有一次输出
}
func main() {
    for i := 0; i < 40; i++ {
        go test(i)
    }
    fmt.Println("start all")
    cond.Broadcast() //  下发广播给所有等待的 goroutine
    time.Sleep(time.Second * 60)
}
```

主 gouroutine 开启后，可以创建多个从 gouroutine，从 gouroutine 获取锁后，进入 cond.Wait 状态，当主 gouroutine 执行完任务后，通过 BroadCast 广播信号。 处于 cond.Wait 状态的所有 gouroutine 收到信号后将全部被唤醒并往下执行。需要注意的是，从 gouroutine 执行完任务后，需要通过 cond.L.Unlock 释放锁， 否则其它被唤醒的 gouroutine 将没法继续执行。 通过查看 cond.Wait 的源码就明白为什么需要需要释放锁了

```go
func (c *Cond) Wait() {
    c.checker.check()
    if raceenabled {
        raceDisable()
    }
    atomic.AddUint32(&c.waiters, 1)
    if raceenabled {
        raceEnable()
    }
    c.L.Unlock()
    runtime_Syncsemacquire(&c.sema)
    c.L.Lock()
}
```

Cond.Wait 会自动释放锁等待信号的到来，当信号到来后，第一个获取到信号的 Wait 将继续往下执行并从新上锁，如果不释放锁， 其它收到信号的 gouroutine 将阻塞无法继续执行。 由于各个 Wait 收到信号的时间是不确定的，因此每次的输出顺序也都是随机的。

