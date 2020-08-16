# sync.Map

转载自：[http://c.biancheng.net/view/34.html](http://c.biancheng.net/view/34.html)

Go语言中的 map 在并发情况下，只读是线程安全的，同时读写是线程不安全的。

 下面来看下并发情况下读写 map 时会出现的问题，代码如下：

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

 fatal error: concurrent map read and map write

错误信息显示，并发的 map 读和 map 写，也就是说使用了两个并发函数不断地对 map 进行读和写而发生了竞态问题，map 内部会对这种并发操作进行检查并提前发现。

需要并发读写时，一般的做法是加锁，但这样性能并不高，Go语言在 1.9 版本中提供了一种效率较高的并发安全的 sync.Map，sync.Map 和 map 不同，不是以语言原生形态提供，而是在 sync 包下的特殊结构。

 sync.Map 有以下特性：

*  无须初始化，直接声明即可。
*  sync.Map 不能使用 map 的方式进行取值和设置等操作，而是使用 sync.Map 的方法进行调用，Store 表示存储，Load 表示获取，Delete 表示删除。
*  使用 Range 配合一个回调函数进行遍历操作，通过回调函数返回内部遍历出来的值，Range 参数中回调函数的返回值在需要继续迭代遍历时，返回 true，终止迭代遍历时，返回 false。

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

