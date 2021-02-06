# GC

## 思维导图

![](../../.gitbook/assets/image%20%2844%29.png)

## 什么是GC

垃圾回收\(Garbage Collection，简称GC\)是编程语言中提供的自动的内存管理机制，自动释放不需要的对象，让出存储器资源，无需程序员手动执行。

​ Golang中的垃圾回收主要应用三色标记法，GC过程和其他用户goroutine可并发运行，但需要一定时间的**STW\(stop the world\)**，STW的过程中，CPU不执行用户代码，全部用于垃圾回收，这个过程的影响很大，Golang进行了多次的迭代优化来解决这个问题。

## 推荐资源

{% embed url="https://mp.weixin.qq.com/s/CsHcVpZ\_9rhO3aJy1-gBaA" %}

{% embed url="https://www.kancloud.cn/aceld/golang/1958308\#\_\_397" %}

{% embed url="https://www.bilibili.com/video/BV1hv411x7we?p=19" %}



