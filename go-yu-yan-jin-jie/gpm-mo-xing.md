# GPM模型

## 什么是GPM

* G代表一个goroutine，即我们要执行的具体任务
* P代表CPU的一个核，一个CPU的最大核数，即_GOMAXPROCS_限制了P的个数，也就设定最大并发数。
* M代表操作系统的线程，对 M 来说，P 提供了相关的执行环境\(Context\)，如内存分配状态\(mcache\)，任务队列\(G\)等。

而且他们三个的关系如下

* _G_需要绑定在_M_上才能运行；
* _M_需要绑定_P_才能运行；
* 程序中的多个_M_并不会同时都处于执行状态，最多只有_GOMAXPROCS_个_M_在执行。

通过引入_P_，实现了一种叫做_work-stealing_的调度算法：

* 每个_P_维护一个_G_队列；
* 当一个_G_被创建出来，或者变为可执行状态时，就把他放到_P_的可执行队列中；
* 当一个_G_执行结束时，_P_会从队列中把该_G_取出；如果此时_P_的队列为空，即没有其他_G_可以执行， 就随机选择另外一个_P_，从其可执行的_G_队列中偷取一半。

该算法避免了在Goroutine调度时使用全局锁。模型结构如下：

![&#x6765;&#x6E90;&#xFF1A;https://wudaijun.com/2018/01/go-scheduler/](../.gitbook/assets/image%20%2824%29.png)

## 调度流程

在M与P绑定后，M会不断从P的Local队列\(runq\)中取出G\(无锁操作\)，切换到G的堆栈并执行，当P的Local队列中没有G时，再从Global队列中返回一个G\(有锁操作，因此实际还会从Global队列批量转移一批G到P Local队列\)，当Global队列中也没有待运行的G时，则尝试从其它的P窃取\(steal\)部分G来执行。

当某个M陷入系统调用时，P则”抛妻弃子”，与M解绑，让阻塞的M和G等待被OS唤醒，而P则带着local queue剩下的G去找一个\(或新建一个\)idle的M，当阻塞的M被唤醒时，它会尝试给G找一个新的归宿\(idle的P，或扔到global queue，等待被领养\)。

## 参考文献

{% embed url="https://www.zhihu.com/question/20862617/answer/36191625" %}

{% embed url="https://morsmachine.dk/go-scheduler" %}

{% embed url="https://juejin.im/entry/6844903609813958669" %}



