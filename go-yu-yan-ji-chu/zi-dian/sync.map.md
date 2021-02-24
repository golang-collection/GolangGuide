# sync.Map

## 什么是sync.map

go原生的map不是线程安全的，如果对他进行并发操作需要加锁，而sync.map则是并发安全的map。对于sync.map适合应用在大量读，少量写的场景下。

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

### Load读取

```go
func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
	// 如果没在 read 中找到，并且 amended 为 true，即 dirty 中存在 read 中没有的 key
	if !ok && read.amended {
		m.mu.Lock() // dirty map 不是线程安全的，所以需要加上互斥锁
		// double check。避免在上锁的过程中 dirty map 提升为 read map。
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
		// 仍然没有在 read 中找到这个 key，并且 amended 为 true
		if !ok && read.amended {
			e, ok = m.dirty[key] // 从 dirty 中找
			// 不管 dirty 中有没有找到，都要"记一笔"，因为在 dirty 提升为 read 之前，都会进入这条路径
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok { // 如果没找到，返回空，false
		return nil, false
	}
	return e.load()
}

func (m *Map) missLocked() {
    m.misses++
    if m.misses < len(m.dirty) {
        return
    }
    
    // 将dirty置给read，因为穿透概率太大了(原子操作，耗时很小)
    m.read.Store(readOnly{m: m.dirty})
    m.dirty = nil
    m.misses = 0
}
```

算法流程如下：

![](../../.gitbook/assets/image%20%2869%29.png)

* 如何设置阀值？这里采用**miss计数和dirty长度**的比较，来进行阀值的设定。
* 为什么dirty可以直接换到read？因为写操作只会操作dirty，所以保证了dirty是最新的，并且数据集是肯定包含read的。 （可能有同学疑问，dirty不是下一步就置为nil了，为何还包含？后文会有解释。）
* 为什么dirty置为nil？我不确定这个原因。猜测：一方面是当read完全等于dirty的时候，读的话read没有就是没有了，即使穿透也是一样的结果，所以存的没啥用。另一方是当存的时候，如果元素比较多，影响插入效率。

### 存储Store

```go
// 它是一个指向任意类型的指针，用来标记从 dirty map 中删除的 entry。
var expunged = unsafe.Pointer(new(interface{}))

// Store sets the value for a key.
func (m *Map) Store(key, value interface{}) {
	// 如果 read map 中存在该 key  则尝试直接更改(由于修改的是 entry 内部的 pointer，因此 dirty map 也可见)
	read, _ := m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			// 如果 read map 中存在该 key，但 p == expunged，则说明 m.dirty != nil 并且 m.dirty 中不存在该 key 值 此时:
			//    a. 将 p 的状态由 expunged  更改为 nil
			//    b. dirty map 插入 key
			m.dirty[key] = e
		}
		// 更新 entry.p = value (read map 和 dirty map 指向同一个 entry)
		e.storeLocked(&value)
	} else if e, ok := m.dirty[key]; ok {
		// 如果 read map 中不存在该 key，但 dirty map 中存在该 key，直接写入更新 entry(read map 中仍然没有这个 key)
		e.storeLocked(&value)
	} else {
		// 如果 read map 和 dirty map 中都不存在该 key，则：
		//	  a. 如果 dirty map 为空，则需要创建 dirty map，并从 read map 中拷贝未删除的元素到新创建的 dirty map
		//    b. 更新 amended 字段，标识 dirty map 中存在 read map 中没有的 key
		//    c. 将 kv 写入 dirty map 中，read 不变
		if !read.amended {
		    // 到这里就意味着，当前的 key 是第一次被加到 dirty map 中。
			// store 之前先判断一下 dirty map 是否为空，如果为空，就把 read map 浅拷贝一次。
			m.dirtyLocked()
			m.read.Store(readOnly{m: read.m, amended: true})
		}
		// 写入新 key，在 dirty 中存储 value
		m.dirty[key] = newEntry(value)
	}
	m.mu.Unlock()
}

// 如果 entry 没被删，tryStore 存储值到 entry 中。如果 p == expunged，即 entry 被删，那么返回 false。
func (e *entry) tryStore(i *interface{}) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expunged {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

// unexpungeLocked 函数确保了 entry 没有被标记成已被清除。
// 如果 entry 先前被清除过了，那么在 mutex 解锁之前，它一定要被加入到 dirty map 中
func (e *entry) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expunged, nil)
}
```

算法流程

![](../../.gitbook/assets/image%20%2868%29.png)

1. 如果在 read 里能够找到待存储的 key，并且对应的 entry 的 p 值不为 expunged，也就是没被删除时，直接更新对应的 entry 即可。
2. 第一步没有成功：要么 read 中没有这个 key，要么 key 被标记为删除。则先加锁，再进行后续的操作。
3. 再次在 read 中查找是否存在这个 key，也就是 double check 一下，这也是 lock-free 编程里的常见套路。如果 read 中存在该 key，但 `p == expunged`，说明 m.dirty != nil 并且 m.dirty 中不存在该 key 值 此时: a. 将 p 的状态由 expunged 更改为 nil；b. dirty map 插入 key。然后，直接更新对应的 value。
4. 如果 read 中没有此 key，那就查看 dirty 中是否有此 key，如果有，则直接更新对应的 value，这时 read 中还是没有此 key。
5. 最后一步，如果 read 和 dirty 中都不存在该 key，则：a. 如果 dirty 为空，则需要创建 dirty，并从 read 中拷贝未被删除的元素；b. 更新 amended 字段，标识 dirty map 中存在 read map 中没有的 key；c. 将 k-v 写入 dirty map 中，read.m 不变。最后，更新此 key 对应的 value。

### 删除Delete

```go
// Delete deletes the value for a key.
func (m *Map) Delete(key interface{}) {
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
	// 如果 read 中没有这个 key，且 dirty map 不为空
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
		if !ok && read.amended {
			delete(m.dirty, key) // 直接从 dirty 中删除这个 key
		}
		m.mu.Unlock()
	}
	if ok {
		e.delete() // 如果在 read 中找到了这个 key，将 p 置为 nil
	}
}

func (e *entry) delete() (hadValue bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expunged {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return true
		}
	}
}
```

算法流程

![](../../.gitbook/assets/image%20%2867%29.png)

可以看到，基本套路还是和 Load，Store 类似，都是先从 read 里查是否有这个 key，如果有则执行 `entry.delete` 方法，将 p 置为 nil，这样 read 和 dirty 都能看到这个变化。

如果没在 read 中找到这个 key，并且 dirty 不为空，那么就要操作 dirty 了，操作之前，还是要先上锁。然后进行 double check，如果仍然没有在 read 里找到此 key，则从 dirty 中删掉这个 key。但不是真正地从 dirty 中删除，而是更新 entry 的状态。

delete函数真正做的事情是将正常状态（指向一个 interface{}）的 p 设置成 nil。没有设置成 expunged 的原因是，当 p 为 expunged 时，表示它已经不在 dirty 中了。这是 p 的状态机决定的，在 `tryExpungeLocked` 函数中，会将 nil 原子地设置成 expunged。

`tryExpungeLocked` 是在新创建 dirty 时调用的，会将已被删除的 entry.p 从 nil 改成 expunged，这个 entry 就不会写入 dirty 了。

```go
func (e *entry) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		// 如果原来是 nil，说明原 key 已被删除，则将其转为 expunged。
		if atomic.CompareAndSwapPointer(&e.p, nil, expunged) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expunged
}
```

注意到如果 key 同时存在于 read 和 dirty 中时，删除只是做了一个标记，将 p 置为 nil；而如果仅在 dirty 中含有这个 key 时，会直接删除这个 key。原因在于，若两者都存在这个 key，仅做标记删除，可以在下次查找这个 key 时，命中 read，提升效率。若只有在 dirty 中存在时，read 起不到“缓存”的作用，直接删除。



## 推荐阅读

{% embed url="https://www.cnblogs.com/qcrao-2018/p/12833787.html" %}

{% embed url="https://juejin.cn/post/6844903895227957262" %}



