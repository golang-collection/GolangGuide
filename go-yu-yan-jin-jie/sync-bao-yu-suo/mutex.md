# Mutex

转载自：[http://c.biancheng.net/view/107.html](http://c.biancheng.net/view/107.html)

Mutex 是最简单的一种锁类型，同时也比较暴力，当一个 goroutine 获得了 Mutex 后，其他 goroutine 就只能乖乖等到这个 goroutine 释放该 Mutex。  
  
RWMutex 相对友好些，是经典的单写多读模型。在读锁占用的情况下，会阻止写，但不阻止读，也就是多个 goroutine 可同时获取读锁（调用 RLock\(\) 方法；而写锁（调用 Lock\(\) 方法）会阻止任何其他 goroutine（无论读和写）进来，整个锁相当于由该 goroutine 独占。从 RWMutex 的实现看，RWMutex 类型其实组合了 Mutex：

```go
type RWMutex struct {
    w Mutex
    writerSem uint32
    readerSem uint32
    readerCount int32
    readerWait int32
}
```

对于这两种锁类型，任何一个 Lock\(\) 或 RLock\(\) 均需要保证对应有 Unlock\(\) 或 RUnlock\(\) 调用与之对应，否则可能导致等待该锁的所有 goroutine 处于饥饿状态，甚至可能导致死锁。锁的典型使用模式如下：

```go
package main
import (
    "fmt"
    "sync"
)
var (
    // 逻辑中使用的某个变量
    count int
    // 与变量对应的使用互斥锁
    countGuard sync.Mutex
)
func GetCount() int {
    // 锁定
    countGuard.Lock()
    // 在函数退出时解除锁定
    defer countGuard.Unlock()
    return count
}
func SetCount(c int) {
    countGuard.Lock()
    count = c
    countGuard.Unlock()
}
func main() {
    // 可以进行并发安全的设置
    SetCount(1)
    // 可以进行并发安全的获取
    fmt.Println(GetCount())
}
```

代码说明如下：

* 第 10 行是某个逻辑步骤中使用到的变量，无论是包级的变量还是结构体成员字段，都可以。
* 第 13 行，一般情况下，建议将互斥锁的粒度设置得越小越好，降低因为共享访问时等待的时间。这里笔者习惯性地将互斥锁的变量命名为以下格式：

  变量名+Guard以表示这个互斥锁用于保护这个变量。

* 第 16 行是一个获取 count 值的函数封装，通过这个函数可以并发安全的访问变量 count。
* 第 19 行，尝试对 countGuard 互斥量进行加锁。一旦 countGuard 发生加锁，如果另外一个 goroutine 尝试继续加锁时将会发生阻塞，直到这个 countGuard 被解锁。
* 第 22 行使用 defer 将 countGuard 的解锁进行延迟调用，解锁操作将会发生在 GetCount\(\) 函数返回时。
* 第 27 行在设置 count 值时，同样使用 countGuard 进行加锁、解锁操作，保证修改 count 值的过程是一个原子过程，不会发生并发访问冲突。

## Mutex的模式



## Mutex可以用作自旋锁吗？



