# 并发通信

转载自：[http://c.biancheng.net/view/4357.html](http://c.biancheng.net/view/4357.html)

通过上一节的学习，关键字 go 的引入使得在Go语言中并发编程变得简单而优雅，但我们同时也应该意识到并发编程的原生复杂性，并时刻对并发中容易出现的问题保持警惕。

 事实上，不管是什么平台，什么编程语言，不管在哪，并发都是一个大话题。并发编程的难度在于协调，而协调就要通过交流，从这个角度看来，并发单元间的通信是最大的问题。

 在工程上，有两种最常见的并发通信模型：共享数据和消息。

 共享数据是指多个并发单元分别保持对同一个数据的引用，实现对该数据的共享。被共享的数据可能有多种形式，比如内存数据块、磁盘文件、网络数据等。在实际工程应用中最常见的无疑是内存了，也就是常说的共享内存。

 先看看我们在C语言中通常是怎么处理线程间数据共享的，代码如下所示。

```c
#include 
#include 
#include 
void *count();
pthread_mutex_t mutex1 = PTHREAD_MUTEX_INITIALIZER;
int counter = 0;
int main()
{
    int rc1, rc2;
    pthread_t thread1, thread2;
    /* 创建线程，每个线程独立执行函数functionC */
    if((rc1 = pthread_create(&thread1, NULL, &count, NULL)))
    {
        printf("Thread creation failed: %d\n", rc1);
    }
    if((rc2 = pthread_create(&thread2, NULL, &count, NULL)))
    {
        printf("Thread creation failed: %d\n", rc2);
    }
    /* 等待所有线程执行完毕 */
    pthread_join( thread1, NULL);
    pthread_join( thread2, NULL);
    exit(0);
}
void *count()
{
    pthread_mutex_lock( &mutex1 );
    counter++;
    printf("Counter value: %d\n",counter);
    pthread_mutex_unlock( &mutex1 );
}
```

 现在我们尝试将这段C语言代码直接翻译为Go语言代码，代码如下所示。

```go
package main
import (
    "fmt"
    "runtime"
    "sync"
)
var counter int = 0
func Count(lock *sync.Mutex) {
    lock.Lock()
    counter++
    fmt.Println(counter)
    lock.Unlock()
}
func main() {
    lock := &sync.Mutex{}
    for i := 0; i < 10; i++ {
        go Count(lock)
    }
    for {
        lock.Lock()
        c := counter
        lock.Unlock()
        runtime.Gosched()
        if c >= 10 {
            break
        }
    }
}
```

 在上面的例子中，我们在 10 个 goroutine 中共享了变量 counter。每个 goroutine 执行完成后，会将 counter 的值加 1。因为 10 个 goroutine 是并发执行的，所以我们还引入了锁，也就是代码中的 lock 变量。每次对 n 的操作，都要先将锁锁住，操作完成后，再将锁打开。

 在 main 函数中，使用 for 循环来不断检查 counter 的值（同样需要加锁）。当其值达到 10 时，说明所有 goroutine 都执行完毕了，这时主函数返回，程序退出。

 事情好像开始变得糟糕了。实现一个如此简单的功能，却写出如此臃肿而且难以理解的代码。想象一下，在一个大的系统中具有无数的锁、无数的共享变量、无数的业务逻辑与错误处理分支，那将是一场噩梦。这噩梦就是众多 C/C++ 开发者正在经历的，其实 Java 和 C\# 开发者也好不到哪里去。

 Go语言既然以并发编程作为语言的最核心优势，当然不至于将这样的问题用这么无奈的方式来解决。Go语言提供的是另一种通信模型，即以消息机制而非共享内存作为通信方式。

 消息机制认为每个并发单元是自包含的、独立的个体，并且都有自己的变量，但在不同并发单元间这些变量不共享。每个并发单元的输入和输出只有一种，那就是消息。这有点类似于进程的概念，每个进程不会被其他进程打扰，它只做好自己的工作就可以了。不同进程间靠消息来通信，它们不会共享内存。

 Go语言提供的消息通信机制被称为 channel，关于 channel 的介绍将在后续的学习中为大家讲解。

