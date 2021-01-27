# 空结构体 struct{}

> 来源：[https://geektutu.com/post/hpg-empty-struct.html](https://geektutu.com/post/hpg-empty-struct.html)

## 空结构体占用空间么 <a id="1-&#x7A7A;&#x7ED3;&#x6784;&#x4F53;&#x5360;&#x7528;&#x7A7A;&#x95F4;&#x4E48;"></a>

在 Go 语言中，我们可以使用 `unsafe.Sizeof` 计算出一个数据类型实例需要占用的字节数。

```go
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println(unsafe.Sizeof(struct{}{}))
}
```

运行上面的例子将会输出：

```go
$ go run main.go
0
```

也就是说，空结构体 struct{} 实例不占据任何的内存空间。

## 空结构体的作用 <a id="2-&#x7A7A;&#x7ED3;&#x6784;&#x4F53;&#x7684;&#x4F5C;&#x7528;"></a>

因为空结构体不占据内存空间，因此被广泛作为各种场景下的占位符使用。一是节省资源，二是空结构体本身就具备很强的语义，即这里不需要任何值，仅作为占位符。

### 实现集合\(Set\) <a id="2-1-&#x5B9E;&#x73B0;&#x96C6;&#x5408;-Set"></a>

Go 语言标准库没有提供 Set 的实现，通常使用 map 来代替。事实上，对于集合来说，只需要 map 的键，而不需要值。即使是将值设置为 bool 类型，也会多占据 1 个字节，那假设 map 中有一百万条数据，就会浪费 1MB 的空间。

因此呢，将 map 作为集合\(Set\)使用时，可以将值类型定义为空结构体，仅作为占位符使用即可。

```go
type Set map[string]struct{}

func (s Set) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) Delete(key string) {
	delete(s, key)
}

func main() {
	s := make(Set)
	s.Add("Tom")
	s.Add("Sam")
	fmt.Println(s.Has("Tom"))
	fmt.Println(s.Has("Jack"))
}
```

### 不发送数据的信道\(channel\) <a id="2-2-&#x4E0D;&#x53D1;&#x9001;&#x6570;&#x636E;&#x7684;&#x4FE1;&#x9053;-channel"></a>

```go
func worker(ch chan struct{}) {
	<-ch
	fmt.Println("do something")
	close(ch)
}

func main() {
	ch := make(chan struct{})
	go worker(ch)
	ch <- struct{}{}
}
```

有时候使用 channel 不需要发送任何的数据，只用来通知子协程\(goroutine\)执行任务，或只用来控制协程并发度。这种情况下，使用空结构体作为占位符就非常合适了。

### 仅包含方法的结构体 <a id="2-3-&#x4EC5;&#x5305;&#x542B;&#x65B9;&#x6CD5;&#x7684;&#x7ED3;&#x6784;&#x4F53;"></a>

```go
type Door struct{}

func (d Door) Open() {
	fmt.Println("Open the door")
}

func (d Door) Close() {
	fmt.Println("Close the door")
}
```

在部分场景下，结构体只包含方法，不包含任何的字段。例如上面例子中的 `Door`，在这种情况下，`Door` 事实上可以用任何的数据结构替代。例如：

```go
type Door int
type Door bool
```

无论是 `int` 还是 `bool` 都会浪费额外的内存，因此呢，这种情况下，声明为空结构体是最合适的。

