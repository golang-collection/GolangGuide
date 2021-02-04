# 接口底层实现

## 类型元数据

在go语言中存在内置类型\(例如int，string等\)与自定义类型，在自定义类型上我们可以定义方法。无论是自定义类型还是内置类型，在go语言内部都存在一个唯一的类型描述信息，他就是runtime.\_type结构体也称为类型元数据，对于每一种类型的描述是全局唯一的，该类型的定义如下：

```go
type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32 //类型的hash值
	tflag      tflag
	align      uint8  //对其边界
	fieldAlign uint8
	kind       uint8  //是否是自定义类型
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}
```

每种类型除了包含\_type类型元数据外，还存在一些其他描述信息，例如用于描述slice类型的slicetype

```go
type slicetype struct {
    typ  _type
    elem *_type
}
```

该elem会指向slice存储的数据类型的类型元数据，如果当前slice存储的是string那么elem就会指向stringtype。

对于自定义类型，除了上面的信息外还会包括一个uncommon结构体，用于描述想包路径，方法数量等，其定义如下

```go
type uncommontype struct {
    pkgpath nameOff //包路径
    mcount  uint16 // number of methods
    xcount  uint16 // number of exported methods
    moff    uint32 // offset from this uncommontype to [mcount]method 方法元数据数组的偏移值
    _       uint32 // unused
}
```

对于方法的结构体其定义如下

```go
type method struct {
    name nameOff // name of method
    mtyp typeOff // method type (without receiver)
    ifn  textOff // fn used in interface call (one-word receiver)
    tfn  textOff // fn used for normal method call
}
```

## 接口

接口存在两种，一种是空接口，另一种是内部含有方法列表的接口也就成为非空接口。

## 空接口

对于空接口类型，go语言中用eface结构体来描述

```go
type eface struct {
    _type *_type //接口对应的类型元数据
    data  unsafe.Pointer //接口指向的具体的动态值
}
```

有下面这样一段代码

```go
var e interface{}
f, _ := os.Open("test.txt")
e = f
```

那么e中的\_type会指向\*os.File的类型元数据，在uncommontype中可以找到对应的方法元数据信息，而e中的data就会指向f。

## 非空接口

对于非空接口类型，go语言中用iface结构体来描述

```go
type iface struct {
    tab  *itab //
    data unsafe.Pointer
}

type itab struct {
    inter *interfacetype
    _type *_type //动态类型元数据
    hash  uint32 // copy of _type.hash. Used for type switches. 类型的hash值
    _     [4]byte
    fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter. 动态类型实现接口中的方法的地址数组
}

type interfacetype struct {
    typ     _type //类型
    pkgpath name //包名
    mhdr    []imethod //接口要求的方法列表
}
```

下面通过这样一段代码来进行分析

```go
var rw io.ReadWriter
f, _ := os.Open("test.txt")
rw = f
```

那么rw中的data就会指向f，tab中的inter就会指向io.ReadWriter类型，tab中的\_type就会指向os.File，同时会将_\__type中uncommontype实现的方法地址拷贝到fun数组中，使得通过rw变量可以快速访问到方法。

{% hint style="info" %}
这里需要注意的是，如果fun\[0\]=0那么就意味着当前动态类型\_type并没有实现该接口类型。
{% endhint %}

一旦接口类型确定了，其对应的实现类型也确定了，那么itab内容就不会改变也就是可以复用的 ，所以go语言会将其保存在一个hash表中，其中将接口类型与动态类型组合为key\(接口类型的hash值异或动态类型的hash值\)，将这些信息缓存起来，提高效率。在查询时如果能在itab缓存中找到则返回对应的itab，否则创建新的itab结构体并将其添加到itab中。

## 推荐阅读

{% embed url="https://www.cnblogs.com/qcrao-2018/p/10766091.html" %}

{% embed url="https://gfw.go101.org/article/interface.html" %}



