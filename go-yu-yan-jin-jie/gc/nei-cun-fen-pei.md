# 内存分配

> 转载自：[https://www.cnblogs.com/qcrao-2018/p/10520785.html](https://www.cnblogs.com/qcrao-2018/p/10520785.html)

Go语言的runtime抛弃了传统的内存分配方式，改为自主管理，。这样可以自主地实现更好的内存使用模式，比如内存池、预分配等等。这样，不会每次内存分配都需要进行系统调用。

Golang运行时的内存分配算法主要源自Google为 C 语言开发的`TCMalloc算法`，全称`Thread-Caching Malloc`。核心思想就是把内存分为多级管理，从而降低锁的粒度。它将可用的堆内存采用二级分配的方式进行管理：每个线程都会自行维护一个独立的内存池，进行内存分配时优先从该内存池中分配，当内存池不足时才会向全局内存池申请，以避免不同线程对全局内存池的频繁竞争。

## 基础概念 <a id="&#x57FA;&#x7840;&#x6982;&#x5FF5;"></a>

Go在程序启动的时候，会先向操作系统申请一块内存（注意这时还只是一段虚拟的地址空间，并不会真正地分配内存），切成小块后自己进行管理。

申请到的内存块被分配了三个区域，在X64上分别是512MB，16GB，512GB大小。

`arena区域`就是我们所谓的堆区，Go动态分配的内存都是在这个区域，它把内存分割成`8KB`大小的页，一些页组合起来称为`mspan`。

`bitmap区域`标识`arena`区域哪些地址保存了对象，并且用`4bit`标志位表示对象是否包含指针、`GC`标记信息。`bitmap`中一个`byte`大小的内存对应`arena`区域中4个指针大小（指针大小为 8B ）的内存，所以`bitmap`区域的大小是`512GB/(4*8B)=16GB`。  


## 推荐阅读

{% embed url="https://www.cnblogs.com/qcrao-2018/p/10520785.html" %}

{% embed url="https://gfw.go101.org/article/memory-block.html" %}

{% embed url="https://gfw.go101.org/article/memory-layout.html" %}

{% embed url="https://draveness.me/golang/docs/part3-runtime/ch07-memory/golang-memory-allocator/" %}

{% embed url="https://mp.weixin.qq.com/s/3gGbJaeuvx4klqcv34hmmw" %}



