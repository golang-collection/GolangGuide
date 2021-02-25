# map底层原理

map底层使用哈希表来实现，通过与运算的方式来进行选桶，map变量本质上是一个\*hmap指针，指向下面这样一个结构体

```go
type hmap struct {
    // Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
    // Make sure this stays in sync with the compiler's definition.
    count     int // # live cells == size of map.  Must be first (used by len() builtin)
    flags     uint8
    B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
    noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
    hash0     uint32 // hash seed

    buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
    oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
    nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

    extra *mapextra // optional fields
}

type mapextra struct {
    // If both key and elem do not contain pointers and are inline, then we mark bucket
    // type as containing no pointers. This avoids scanning such maps.
    // However, bmap.overflow is a pointer. In order to keep overflow buckets
    // alive, we store pointers to all overflow buckets in hmap.extra.overflow and hmap.extra.oldoverflow.
    // overflow and oldoverflow are only used if key and elem do not contain pointers.
    // overflow contains overflow buckets for hmap.buckets.
    // oldoverflow contains overflow buckets for hmap.oldbuckets.
    // The indirection allows to store a pointer to the slice in hiter.
    overflow    *[]*bmap
    oldoverflow *[]*bmap

    // nextOverflow holds a pointer to a free overflow bucket.
    nextOverflow *bmap
}

// A bucket for a Go map.
type bmap struct {
    // tophash generally contains the top byte of the hash value
    // for each key in this bucket. If tophash[0] < minTopHash,
    // tophash[0] is a bucket evacuation state instead.
    tophash [bucketCnt]uint8
}
// 但这只是表面(src/runtime/hashmap.go)的结构，
// 编译期间会给它加料，动态地创建一个新的结构：
type bmap struct {
    topbits  [8]uint8
    keys     [8]keytype
    values   [8]valuetype
    pad      uintptr
    overflow uintptr
}
```

其中hmap结构体变量解释如下：

* count：代表键值对数目
* B：记录hash桶数目是2的多少倍，其中桶的数量是2的倍数是为了保证防止有的桶不被选中
* noverflow：使用的溢出桶的数量
* buckets：记录桶的地址
* oldbuckets：记录旧桶的地址（扩容期间使用）
* nevacuate：渐进式扩容时即将迁移的桶的编号
* extra：记录溢出桶的信息

{% hint style="info" %}
渐进式扩容：扩容时开辟新的空间，但不是一次性将旧桶内容全部复制到新桶，而是逐渐的将旧的内容复制过去，防止系统出现性能抖动。
{% endhint %}

其中mapextra结构体变量解释如下：

* overflow：记录已经被使用的溢出桶
* olderoverflow：记录旧的已经被使用的溢出桶
* nextOverflow：指向下一个空的溢出桶

对于bmap来说，它的图示如下，内部会讲key放在一起，value放在一起，通过这样的方式使内存排列更加紧凑，第一行存储的是每个key对应的hash值的高八位。

![&#x56FE;&#x7247;&#x6765;&#x6E90;&#xFF1A;https://www.cnblogs.com/qcrao-2018/p/10903807.html](../../.gitbook/assets/image%20%2845%29.png)

为了减少扩容的次数，当桶存满时，只要还有可用的溢出桶就会在对应的bmap后面链接对应的溢出桶。

如果hmap中的B&gt;4，就认为溢出的可能性比较高，就会预先分配2^\(B-4\)个溢出桶用于备用，在内存中其与常规桶是连续的，前2^B个用于常规桶，会面的用于做溢出桶。

## 哈希函数

map 的一个关键点在于，哈希函数的选择。在程序启动时，会检测 cpu 是否支持 aes，如果支持，则使用 aes hash，否则使用 memhash。这是在函数 `alginit()` 中完成，位于路径：`src/runtime/alg.go` 下。

{% hint style="info" %}
hash 函数，有加密型和非加密型。 加密型的一般用于加密数据、数字摘要等，典型代表就是 md5、sha1、sha256、aes256 这种； 非加密型的一般就是查找。在 map 的应用场景中，用的是查找。 选择 hash 函数主要考察的是两点：性能、碰撞概率。
{% endhint %}

表示类型的结构体：

```go
type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldalign uint8
	kind       uint8
	alg        *typeAlg
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}
```

其中 `alg` 字段就和哈希相关，它是指向如下结构体的指针：

```go
// src/runtime/alg.go
type typeAlg struct {
	// (ptr to object, seed) -> hash
	hash func(unsafe.Pointer, uintptr) uintptr
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
}
```

typeAlg 包含两个函数，hash 函数计算类型的哈希值，而 equal 函数则计算两个类型是否“哈希相等”。

对于 string 类型，它的 hash、equal 函数如下：

```go
func strhash(a unsafe.Pointer, h uintptr) uintptr {
	x := (*stringStruct)(a)
	return memhash(x.str, h, uintptr(x.len))
}

func strequal(p, q unsafe.Pointer) bool {
	return *(*string)(p) == *(*string)(q)
}
```

根据 key 的类型，\_type 结构体的 alg 字段会被设置对应类型的 hash 和 equal 函数。

## map扩容规则

在go语言中，默认的负载因子是6.5，当count/\(2^B\)&gt;6.5时就会触发翻倍扩容。

但是当溢出桶使用较多时，也会触发等量扩容。此中情况一般是因为在溢出桶中删除元素较多，每个溢出桶已经没有多少元素，但是仍让占用空间，通过等量扩容的方式使的内存占用更少。

溢出桶是否使用较多判断依据

* B&lt;=15且noverflow&gt;=2^B
* B&gt;15且noverflow&gt;=2^15

## 推荐阅读

{% embed url="https://www.cnblogs.com/qcrao-2018/p/10903807.html" %}

