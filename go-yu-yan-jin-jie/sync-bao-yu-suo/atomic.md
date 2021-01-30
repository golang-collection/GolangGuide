# atomic

> 转载自：[https://github.com/polaris1119/The-Golang-Standard-Library-by-Example/blob/master/chapter16/16.02.md](https://github.com/polaris1119/The-Golang-Standard-Library-by-Example/blob/master/chapter16/16.02.md)

对于并发操作而言，原子操作是个非常现实的问题。典型的就是 i++ 的问题。 当两个 CPU 同时对内存中的 i 进行读取，然后把加一之后的值放入内存中，可能两次 i++ 的结果，这个 i 只增加了一次。 如何保证多 CPU 对同一块内存的操作是原子的。 golang 中 sync/atomic 就是做这个使用的。

具体的原子操作在不同的操作系统中实现是不同的。比如在 Intel 的 CPU 架构机器上，主要是使用总线锁的方式实现的。 大致的意思就是当一个 CPU 需要操作一个内存块的时候，向总线发送一个 LOCK 信号，所有 CPU 收到这个信号后就不对这个内存块进行操作了。 等待操作的 CPU 执行完操作后，发送 UNLOCK 信号，才结束。 在 AMD 的 CPU 架构机器上就是使用 MESI 一致性协议的方式来保证原子操作。 所以我们在看 atomic 源码的时候，我们看到它针对不同的操作系统有不同汇编语言文件。

如果我们善用原子操作，它会比锁更为高效。

## CAS

原子操作中最经典的 CAS\(compare-and-swap\) 在 atomic 包中是 Compare 开头的函数。

* func CompareAndSwapInt32\(addr \*int32, old, new int32\) \(swapped bool\)
* func CompareAndSwapInt64\(addr \*int64, old, new int64\) \(swapped bool\)
* func CompareAndSwapPointer\(addr \*unsafe.Pointer, old, new unsafe.Pointer\) \(swapped bool\)
* func CompareAndSwapUint32\(addr \*uint32, old, new uint32\) \(swapped bool\)
* func CompareAndSwapUint64\(addr \*uint64, old, new uint64\) \(swapped bool\)
* func CompareAndSwapUintptr\(addr \*uintptr, old, new uintptr\) \(swapped bool\)

CAS 的意思是判断内存中的某个值是否等于 old 值，如果是的话，则赋 new 值给这块内存。CAS 是一个方法，并不局限在 CPU 原子操作中。 CAS 比互斥锁乐观，但是也就代表 CAS 是有赋值不成功的时候，调用 CAS 的那一方就需要处理赋值不成功的后续行为了。

这一系列的函数需要比较后再进行交换，也有不需要进行比较就进行交换的原子操作。

* func SwapInt32\(addr \*int32, new int32\) \(old int32\)
* func SwapInt64\(addr \*int64, new int64\) \(old int64\)
* func SwapPointer\(addr \*unsafe.Pointer, new unsafe.Pointer\) \(old unsafe.Pointer\)
* func SwapUint32\(addr \*uint32, new uint32\) \(old uint32\)
* func SwapUint64\(addr \*uint64, new uint64\) \(old uint64\)
* func SwapUintptr\(addr \*uintptr, new uintptr\) \(old uintptr\)

## 增加或减少

对一个数值进行增加或者减少的行为也需要保证是原子的，它对应于 atomic 包的函数就是

* func AddInt32\(addr \*int32, delta int32\) \(new int32\)
* func AddInt64\(addr \*int64, delta int64\) \(new int64\)
* func AddUint32\(addr \*uint32, delta uint32\) \(new uint32\)
* func AddUint64\(addr \*uint64, delta uint64\) \(new uint64\)
* func AddUintptr\(addr \*uintptr, delta uintptr\) \(new uintptr\)

## 读取或写入

当我们要读取一个变量的时候，很有可能这个变量正在被写入，这个时候，我们就很有可能读取到写到一半的数据。 所以读取操作是需要一个原子行为的。在 atomic 包中就是 Load 开头的函数群。

* func LoadInt32\(addr \*int32\) \(val int32\)
* func LoadInt64\(addr \*int64\) \(val int64\)
* func LoadPointer\(addr \*unsafe.Pointer\) \(val unsafe.Pointer\)
* func LoadUint32\(addr \*uint32\) \(val uint32\)
* func LoadUint64\(addr \*uint64\) \(val uint64\)
* func LoadUintptr\(addr \*uintptr\) \(val uintptr\)

好了，读取我们是完成了原子性，那写入呢？也是同样的，如果有多个 CPU 往内存中一个数据块写入数据的时候，可能导致这个写入的数据不完整。 在 atomic 包对应的是 Store 开头的函数群。

* func StoreInt32\(addr \*int32, val int32\)
* func StoreInt64\(addr \*int64, val int64\)
* func StorePointer\(addr \*unsafe.Pointer, val unsafe.Pointer\)
* func StoreUint32\(addr \*uint32, val uint32\)
* func StoreUint64\(addr \*uint64, val uint64\)
* func StoreUintptr\(addr \*uintptr, val uintptr\)

