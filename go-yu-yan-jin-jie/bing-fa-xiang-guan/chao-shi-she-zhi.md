# select与超时机制

转载自：[http://c.biancheng.net/view/4361.html](http://c.biancheng.net/view/4361.html)

Go语言没有提供直接的超时处理机制，所谓超时可以理解为当我们上网浏览一些网站时，如果一段时间之后不作操作，就需要重新登录。

那么我们应该如何实现这一功能呢，这时就可以使用 select 来设置超时。

虽然 select 机制不是专门为超时而设计的，却能很方便的解决超时问题，因为 select 的特点是只要其中有一个 case 已经完成，程序就会继续往下执行，而不会考虑其他 case 的情况。

超时机制本身虽然也会带来一些问题，比如在运行比较快的机器或者高速的网络上运行正常的程序，到了慢速的机器或者网络上运行就会出问题，从而出现结果不一致的现象，但从根本上来说，解决死锁问题的价值要远大于所带来的问题。

select 的用法与 switch 语言非常类似，由 select 开始一个新的选择块，每个选择条件由 case 语句来描述。

与 switch 语句相比，select 有比较多的限制，其中最大的一条限制就是每个 case 语句里必须是一个 IO 操作，大致的结构如下：

```go
select {
     case <-chan1:
     // 如果chan1成功读到数据，则进行该case处理语句
     case chan2 <- 1:
     // 如果成功向chan2写入数据，则进行该case处理语句
     default:
     // 如果上面都没有成功，则进入default处理流程
 }
```

 在一个 select 语句中，Go语言会按顺序从头至尾评估每一个发送和接收的语句。

如果其中的任意一语句可以继续执行（即没有被阻塞），那么就从那些可以执行的语句中任意选择一条来使用。

 如果没有任意一条语句可以执行（即所有的通道都被阻塞），那么有如下两种可能的情况：

*  如果给出了 default 语句，那么就会执行 default 语句，同时程序的执行会从 select 语句后的语句中恢复；
*  如果没有 default 语句，那么 select 语句将被阻塞，直到至少有一个通信可以进行下去。

 示例代码如下所示：

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan int)
    quit := make(chan bool)

    //新开一个协程
    go func() {
        for {
            select {
            case num := <-ch:
                fmt.Println("num = ", num)
            //超过三秒即为超时    
            case <-time.After(3 * time.Second):
                fmt.Println("超时")
                quit <- true
            }
        }

    }() //别忘了()

    for i := 0; i < 5; i++ {
        ch <- i
        time.Sleep(time.Second)
    }

    <-quit
    fmt.Println("程序结束")
}
```

 运行结果如下：

 num =  0  
 num =  1  
 num =  2  
 num =  3  
 num =  4  
 超时  
 程序结束

