# sync.Map

## 什么是sync.map

go原生的map不是线程安全的，如果对他进行并发操作需要加锁，而sync.map则是并发安全的map。

## 怎么用

Go语言原生的 map 在并发情况下，只读是线程安全的，同时读写是线程不安全的。下面来看下并发情况下读写 map 时会出现的问题，代码如下：

```go
// 创建一个int到int的映射
m := make(map[int]int)

// 开启一段并发代码
go func() {

    // 不停地对map进行写入
    for {
        m[1] = 1
    }

}()

// 开启一段并发代码
go func() {

    // 不停地对map进行读取
    for {
        _ = m[1]
    }

}()

// 无限循环, 让并发程序在后台执行
for {

}
```

运行代码会报错，输出如下：

```go
fatal error: concurrent map read and map write
```

错误信息显示，并发的 map 读和 map 写，也就是说使用了两个并发函数不断地对 map 进行读和写而发生了竞态问题，map 内部会对这种并发操作进行检查并提前发现。

 sync.Map 有以下特性：

* 无须初始化，直接声明即可。
* sync.Map 不能使用 map 的方式进行取值和设置等操作，而是使用 sync.Map 的方法进行调用，Store 表示存储，Load 表示获取，Delete 表示删除。
* 使用 Range 配合一个回调函数进行遍历操作，通过回调函数返回内部遍历出来的值，Range 参数中回调函数的返回值在需要继续迭代遍历时，返回 true，终止迭代遍历时，返回 false。

 并发安全的 sync.Map 演示代码如下：

```go
package main

import (
      "fmt"
      "sync"
)

func main() {

    var scene sync.Map

    // 将键值对保存到sync.Map
    scene.Store("greece", 97)
    scene.Store("london", 100)
    scene.Store("egypt", 200)

    // 从sync.Map中根据键取值
    fmt.Println(scene.Load("london"))

    // 根据键删除对应的键值对
    scene.Delete("london")

    // 遍历所有sync.Map中的键值对
    scene.Range(func(k, v interface{}) bool {

        fmt.Println("iterate:", k, v)
        return true
    })

}
```

 代码输出如下：

```go
100 true
iterate: egypt 200
iterate: greece 97
```

 代码说明如下：

*  第 10 行，声明 scene，类型为 sync.Map，注意，sync.Map 不能使用 make 创建。
*  第 13～15 行，将一系列键值对保存到 sync.Map 中，sync.Map 将键和值以 interface{} 类型进行保存。
*  第 18 行，提供一个 sync.Map 的键给 scene.Load\(\) 方法后将查询到键对应的值返回。
*  第 21 行，sync.Map 的 Delete 可以使用指定的键将对应的键值对删除。
*  第 24 行，Range\(\) 方法可以遍历 sync.Map，遍历需要提供一个匿名函数，参数为 k、v，类型为 interface{}，每次 Range\(\) 在遍历一个元素时，都会调用这个匿名函数把结果返回。

 sync.Map 没有提供获取 map 数量的方法，替代方法是在获取 sync.Map 时遍历自行计算数量，sync.Map 为了保证并发安全有一些性能损失，因此在非并发情况下，使用 map 相比使用 sync.Map 会有更好的性能。

## 底层实现

数据结构

```go
// Map is like a Go map[interface{}]interface{} but is safe for concurrent use
// by multiple goroutines without additional locking or coordination.
// Loads, stores, and deletes run in amortized constant time.
//
// The Map type is specialized. Most code should use a plain Go map instead,
// with separate locking or coordination, for better type safety and to make it
// easier to maintain other invariants along with the map content.
//
// The Map type is optimized for two common use cases: (1) when the entry for a given
// key is only ever written once but read many times, as in caches that only grow,
// or (2) when multiple goroutines read, write, and overwrite entries for disjoint
// sets of keys. In these two cases, use of a Map may significantly reduce lock
// contention compared to a Go map paired with a separate Mutex or RWMutex.
//
// The zero Map is empty and ready for use. A Map must not be copied after first use.
type Map struct {
	mu Mutex

	// read contains the portion of the map's contents that are safe for
	// concurrent access (with or without mu held).
	//
	// The read field itself is always safe to load, but must only be stored with
	// mu held.
	//
	// Entries stored in read may be updated concurrently without mu, but updating
	// a previously-expunged entry requires that the entry be copied to the dirty
	// map and unexpunged with mu held.
	read atomic.Value // readOnly

	// dirty contains the portion of the map's contents that require mu to be
	// held. To ensure that the dirty map can be promoted to the read map quickly,
	// it also includes all of the non-expunged entries in the read map.
	//
	// Expunged entries are not stored in the dirty map. An expunged entry in the
	// clean map must be unexpunged and added to the dirty map before a new value
	// can be stored to it.
	//
	// If the dirty map is nil, the next write to the map will initialize it by
	// making a shallow copy of the clean map, omitting stale entries.
	dirty map[interface{}]*entry

	// misses counts the number of loads since the read map was last updated that
	// needed to lock mu to determine whether the key was present.
	//
	// Once enough misses have occurred to cover the cost of copying the dirty
	// map, the dirty map will be promoted to the read map (in the unamended
	// state) and the next store to the map will make a new dirty copy.
	misses int
}

// readOnly is an immutable struct stored atomically in the Map.read field.
type readOnly struct {
	m       map[interface{}]*entry
	amended bool // true if the dirty map contains some key not in m.
}
```

对于上面的数据结构中字段解释如下

* mu是互斥量也就是锁，用于保护dirty
* read是atomic.Value类型，可以并发的读取，在使用中必须要保证load和store的原子性，其实际存储的是上面代码段中的readOnly结构体。atomic.Value对应的定义如下

![](../../.gitbook/assets/image%20%2865%29.png)

* dirty是一个非线程安全的原始map。包含新写入的key，并且包含read中的所有未被删除的key。这样可以快速的将dirty提升为read对外提供服务。
* misses：计数作用。每次从read中读失败，则计数+1。

在read和dirty中都包含entry，他是这样的一个结构体

```go
// An entry is a slot in the map corresponding to a particular key.
type entry struct {
	// p points to the interface{} value stored for the entry.
	//
	// If p == nil, the entry has been deleted and m.dirty == nil.
	//
	// If p == expunged, the entry has been deleted, m.dirty != nil, and the entry
	// is missing from m.dirty.
	//
	// Otherwise, the entry is valid and recorded in m.read.m[key] and, if m.dirty
	// != nil, in m.dirty[key].
	//
	// An entry can be deleted by atomic replacement with nil: when m.dirty is
	// next created, it will atomically replace nil with expunged and leave
	// m.dirty[key] unset.
	//
	// An entry's associated value can be updated by atomic replacement, provided
	// p != expunged. If p == expunged, an entry's associated value can be updated
	// only after first setting m.dirty[key] = e so that lookups using the dirty
	// map find the entry.
	p unsafe.Pointer // *interface{}
}
```

其内部仅保留了一个指针指向的就是具体的value值。这个指针的状态有三种

* p==nil，表示该键值对已被删除且m.dirty==nil。
* p==expunged，说明该键值对已经被删除，并且m.dirty!=nil，并且m.dirty中不含这个entity。
* 其他情况，p指向一个正常的值，并且被记录在m.read.m\[key\]中，如果m.dirty也不为空，那么该entity也会被存储在m.dirty\[key\]中。实际上就是read和dirty会指向同一个值。

对于sync的map操作是这样的，写会直接将内容写入dirty中，读取时会先从read中读，没有得到的话再从dirty中读取。

![](../../.gitbook/assets/image%20%2866%29.png)



