# Context

> 转载自：[https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-context/](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-context/)

上下文 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 是用来设置截止日期、同步信号，传递请求相关值的结构体。上下文与 Goroutine 有比较密切的关系。[`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 是 Go 语言中独特的设计，在其他编程语言中我们很少见到类似的概念。

[`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 是 Go 语言在 1.7 版本中引入标准库的接口[1]()，该接口定义了四个需要实现的方法，其中包括：

1. `Deadline` — 返回 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 被取消的时间，也就是完成工作的截止日期；
2. `Done` — 返回一个 Channel，这个 Channel 会在当前工作完成或者上下文被取消之后关闭，多次调用 `Done` 方法会返回同一个 Channel；
3. `Err` — 返回 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 结束的原因，它只会在 `Done` 返回的 Channel 被关闭时才会返回非空的值；
   1. 如果 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 被取消，会返回 `Canceled` 错误；
   2. 如果 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 超时，会返回 `DeadlineExceeded` 错误；
4. `Value` — 从 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 中获取键对应的值，对于同一个上下文来说，多次调用 `Value` 并传入相同的 `Key` 会返回相同的结果，该方法可以用来传递请求特定的数据；

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}
```

[`context`](https://github.com/golang/go/tree/master/src/context) 包中提供的 [`context.Background`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L208-L210)、[`context.TODO`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L216-L218)、[`context.WithDeadline`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L232-L236) 和 [`context.WithValue`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L513-L521) 函数会返回实现该接口的私有结构体，我们会在后面详细介绍它们的工作原理。

## 设计原理 <a id="611-&#x8BBE;&#x8BA1;&#x539F;&#x7406;"></a>

在 Goroutine 构成的树形结构中对信号进行同步以减少计算资源的浪费是 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 的最大作用。Go 服务的每一个请求的都是通过单独的 Goroutine 处理的[2]()，HTTP/RPC 请求的处理器会启动新的 Goroutine 访问数据库和其他服务。

如下图所示，我们可能会创建多个 Goroutine 来处理一次请求，而 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 的作用就是在不同 Goroutine 之间同步请求特定数据、取消信号以及处理请求的截止日期。

![golang-context-usage](https://img.draveness.me/golang-context-usage.png)

**图 6-1 Context 与 Goroutine 树**

每一个 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 都会从最顶层的 Goroutine 一层一层传递到最下层。[`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 可以在上层 Goroutine 执行出现错误时，将信号及时同步给下层。

![golang-without-context](https://img.draveness.me/golang-without-context.png)

**图 6-2 不使用 Context 同步信号**

如上图所示，当最上层的 Goroutine 因为某些原因执行失败时，下层的 Goroutine 由于没有接收到这个信号所以会继续工作；但是当我们正确地使用 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 时，就可以在下层及时停掉无用的工作以减少额外资源的消耗：

![golang-with-context](https://img.draveness.me/golang-with-context.png)

**图 6-3 使用 Context 同步信号**

我们可以通过一个代码片段了解 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 是如何对信号进行同步的。在这段代码中，我们创建了一个过期时间为 1s 的上下文，并向上下文传入 `handle` 函数，该方法会使用 500ms 的时间处理传入的『请求』：

```text
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go handle(ctx, 500*time.Millisecond)
	select {
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
	}
}

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}
```

因为过期时间大于处理时间，所以我们有足够的时间处理该『请求』，运行上述代码会打印出如下所示的内容：

```text
$ go run context.go
process request with 500ms
main context deadline exceeded
```

`handle` 函数没有进入超时的 `select` 分支，但是 `main` 函数的 `select` 却会等待 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 的超时并打印出 `main context deadline exceeded`。

如果我们将处理『请求』时间增加至 1500ms，整个程序都会因为上下文的过期而被中止，：

```text
$ go run context.go
main context deadline exceeded
handle context deadline exceeded
```

相信这两个例子能够帮助各位读者理解 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 的使用方法和设计原理 — 多个 Goroutine 同时订阅 `ctx.Done()` 管道中的消息，一旦接收到取消信号就立刻停止当前正在执行的工作。

## 默认上下文 <a id="612-&#x9ED8;&#x8BA4;&#x4E0A;&#x4E0B;&#x6587;"></a>

[`context`](https://github.com/golang/go/tree/master/src/context) 包中最常用的方法还是 [`context.Background`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L208-L210)、[`context.TODO`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L216-L218)，这两个方法都会返回预先初始化好的私有变量 `background` 和 `todo`，它们会在同一个 Go 程序中被复用：

```text
func Background() Context {
	return background
}

func TODO() Context {
	return todo
}
```

这两个私有变量都是通过 `new(emptyCtx)` 语句初始化的，它们是指向私有结构体 [`context.emptyCtx`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L171) 的指针，这是最简单、最常用的上下文类型：

```text
type emptyCtx int

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*emptyCtx) Done() <-chan struct{} {
	return nil
}

func (*emptyCtx) Err() error {
	return nil
}

func (*emptyCtx) Value(key interface{}) interface{} {
	return nil
}
```

从上述代码，我们不难发现 [`context.emptyCtx`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L171) 通过返回 `nil` 实现了 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 接口，它没有任何特殊的功能。

![golang-context-hierarchy](https://img.draveness.me/golang-context-hierarchy.png)

**图 6-4 Context 层级关系**

从源代码来看，[`context.Background`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L208-L210) 和 [`context.TODO`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L216-L218) 函数其实也只是互为别名，没有太大的差别。它们只是在使用和语义上稍有不同：

* [`context.Background`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L208-L210) 是上下文的默认值，所有其他的上下文都应该从它衍生（Derived）出来；
* [`context.TODO`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L216-L218) 应该只在不确定应该使用哪种上下文时使用；

在多数情况下，如果当前函数没有上下文作为入参，我们都会使用 [`context.Background`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L208-L210) 作为起始的上下文向下传递。

## 取消信号 <a id="613-&#x53D6;&#x6D88;&#x4FE1;&#x53F7;"></a>

[`context.WithCancel`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L232-L236) 函数能够从 [`context.Context`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L62-L154) 中衍生出一个新的子上下文并返回用于取消该上下文的函数（CancelFunc）。一旦我们执行返回的取消函数，当前上下文以及它的子上下文都会被取消，所有的 Goroutine 都会同步收到这一取消信号。

![golang-parent-cancel-context](https://img.draveness.me/2020-01-20-15795072700927-golang-parent-cancel-context.png)

**图 6-5 Context 子树的取消**

我们直接从 [`context.WithCancel`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L232-L236) 函数的实现来看它到底做了什么：

```text
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	c := newCancelCtx(parent)
	propagateCancel(parent, &c)
	return &c, func() { c.cancel(true, Canceled) }
}
```

* [`context.newCancelCtx`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L239-L241) 将传入的上下文包装成私有结构体 [`context.cancelCtx`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L341-L348)；
* [`context.propagateCancel`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L247-L283) 会构建父子上下文之间的关联，当父上下文被取消时，子上下文也会被取消：

```text
func propagateCancel(parent Context, child canceler) {
	done := parent.Done()
	if done == nil {
		return // 父上下文不会触发取消信号
	}
	select {
	case <-done:
		child.cancel(false, parent.Err()) // 父上下文已经被取消
		return
	default:
	}

	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock()
		if p.err != nil {
			child.cancel(false, p.err)
		} else {
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else {
		go func() {
			select {
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			}
		}()
	}
}
```

上述函数总共与父上下文相关的三种不同的情况：

1. 当 `parent.Done() == nil`，也就是 `parent` 不会触发取消事件时，当前函数会直接返回；
2. 当 `child` 的继承链包含可以取消的上下文时，会判断 `parent` 是否已经触发了取消信号；
   * 如果已经被取消，`child` 会立刻被取消；
   * 如果没有被取消，`child` 会被加入 `parent` 的 `children` 列表中，等待 `parent` 释放取消信号；
3. 在默认情况下
   1. 运行一个新的 Goroutine 同时监听 `parent.Done()` 和 `child.Done()` 两个 Channel
   2. 在 `parent.Done()` 关闭时调用 `child.cancel` 取消子上下文；

[`context.propagateCancel`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L247-L283) 的作用是在 `parent` 和 `child` 之间同步取消和结束的信号，保证在 `parent` 被取消时，`child` 也会收到对应的信号，不会发生状态不一致的问题。

[`context.cancelCtx`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L341-L348) 实现的几个接口方法也没有太多值得分析的地方，该结构体最重要的方法是 `cancel`，这个方法会关闭上下文中的 Channel 并向所有的子上下文同步取消信号：

```text
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return
	}
	c.err = err
	if c.done == nil {
		c.done = closedchan
	} else {
		close(c.done)
	}
	for child := range c.children {
		child.cancel(false, err)
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.Context, c)
	}
}
```

除了 [`context.WithCancel`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L232-L236) 之外，[`context`](https://github.com/golang/go/tree/master/src/context) 包中的另外两个函数 [`context.WithDeadline`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L427-L450) 和 [`context.WithTimeout`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L496-L498) 也都能创建可以被取消的计时器上下文 [`context.timerCtx`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L455-L460)：

```text
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
		return WithCancel(parent)
	}
	c := &timerCtx{
		cancelCtx: newCancelCtx(parent),
		deadline:  d,
	}
	propagateCancel(parent, c)
	dur := time.Until(d)
	if dur <= 0 {
		c.cancel(true, DeadlineExceeded) // 已经过了截止日期
		return c, func() { c.cancel(false, Canceled) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		c.timer = time.AfterFunc(dur, func() {
			c.cancel(true, DeadlineExceeded)
		})
	}
	return c, func() { c.cancel(true, Canceled) }
}
```

[`context.WithDeadline`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L427-L450) 方法在创建 [`context.timerCtx`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L455-L460) 的过程中，判断了父上下文的截止日期与当前日期，并通过 [`time.AfterFunc`](https://github.com/golang/go/blob/001fe7f33f1d7aed9e3a047bd8e784bdc103c28c/src/time/sleep.go#L155-L165) 创建定时器，当时间超过了截止日期后会调用 [`context.timerCtx.cancel`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L472-L484) 方法同步取消信号。

[`context.timerCtx`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L455-L460) 结构体内部不仅通过嵌入了[`context.cancelCtx`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L341-L348) 结构体继承了相关的变量和方法，还通过持有的定时器 `timer` 和截止时间 `deadline` 实现了定时取消这一功能：

```text
type timerCtx struct {
	cancelCtx
	timer *time.Timer // Under cancelCtx.mu.

	deadline time.Time
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c *timerCtx) cancel(removeFromParent bool, err error) {
	c.cancelCtx.cancel(false, err)
	if removeFromParent {
		removeChild(c.cancelCtx.Context, c)
	}
	c.mu.Lock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	c.mu.Unlock()
}
```

[`context.timerCtx.cancel`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L472-L484) 方法不仅调用了 [`context.cancelCtx.cancel`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L391-L416)，还会停止持有的定时器减少不必要的资源浪费。

## 传值方法 <a id="614-&#x4F20;&#x503C;&#x65B9;&#x6CD5;"></a>

在最后我们需要了解如何使用上下文传值，[`context`](https://github.com/golang/go/tree/master/src/context) 包中的 [`context.WithValue`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L513-L521) 函数能从父上下文中创建一个子上下文，传值的子上下文使用 [`context.valueCtx`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L525-L528) 类型：

```text
func WithValue(parent Context, key, val interface{}) Context {
	if key == nil {
		panic("nil key")
	}
	if !reflectlite.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parent, key, val}
}
```

[`context.valueCtx`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L525-L528) 结构体会将除了 `Value` 之外的 `Err`、`Deadline` 等方法代理到父上下文中，它只会响应 [`context.valueCtx.Value`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L549-L554) 方法，这个方法的实现也很简单：

```text
type valueCtx struct {
	Context
	key, val interface{}
}

func (c *valueCtx) Value(key interface{}) interface{} {
	if c.key == key {
		return c.val
	}
	return c.Context.Value(key)
}
```

如果 [`context.valueCtx`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L525-L528) 中存储的键值对与 [`context.valueCtx.Value`](https://github.com/golang/go/blob/df2999ef43ea49ce1578137017949c0ee660608a/src/context/context.go#L549-L554) 方法中传入的参数不匹配，就会从父上下文中查找该键对应的值直到在某个父上下文中返回 `nil` 或者查找到对应的值。

## 小结 <a id="615-&#x5C0F;&#x7ED3;"></a>

Go 语言中的 [`context.Context`](https://github.com/golang/go/blob/71bbffbc48d03b447c73da1f54ac57350fc9b36a/src/context/context.go#L62-L154) 的主要作用还是在多个 Goroutine 组成的树中同步取消信号以减少对资源的消耗和占用，虽然它也有传值的功能，但是这个功能我们还是很少用到。

在真正使用传值的功能时我们也应该非常谨慎，使用 [`context.Context`](https://github.com/golang/go/blob/71bbffbc48d03b447c73da1f54ac57350fc9b36a/src/context/context.go#L62-L154) 进行传递参数请求的所有参数一种非常差的设计，比较常见的使用场景是传递请求对应用户的认证令牌以及用于进行分布式追踪的请求 ID。

## 延伸阅读 <a id="616-&#x5EF6;&#x4F38;&#x9605;&#x8BFB;"></a>

* [Package context · Golang](https://golang.org/pkg/context/)
* [Go Concurrency Patterns: Context](https://blog.golang.org/context)
* [Using context cancellation in Go](https://www.sohamkamani.com/blog/golang/2018-06-17-golang-using-context-cancellation/)
* [https://www.flysnow.org/2017/05/12/go-in-action-go-context.html](https://www.flysnow.org/2017/05/12/go-in-action-go-context.html)

