# expvar

> 转载自：[https://github.com/polaris1119/The-Golang-Standard-Library-by-Example/blob/master/chapter13/13.3.md](https://github.com/polaris1119/The-Golang-Standard-Library-by-Example/blob/master/chapter13/13.3.md)

expvar 挺简单的，然而，它也是很有用的。但不幸的是，貌似了解它的人不多。来自 godoc.org 的数据表明，没有多少人知道这个包。截止目前（2017-6-18），该包被公开的项目 import 2207 次，相比较而言，连 image 包都被 import 3491 次之多。

如果你看到了这里，希望以后你的项目中能使用上 expvar 这个包。

### 包简介

包 expvar 为公共变量提供了一个标准化的接口，如服务器中的操作计数器。它以 JSON 格式通过 `/debug/vars` 接口以 HTTP 的方式公开这些公共变量。

设置或修改这些公共变量的操作是原子的。

除了为程序增加 HTTP handler，此包还注册以下变量：

```text
cmdline   os.Args
memstats  runtime.Memstats
```

导入该包有时只是为注册其 HTTP handler 和上述变量。 要以这种方式使用，请将此包通过如下形式引入到程序中：

```text
import _ "expvar"
```

### 例子

在继续介绍此包的详细信息之前，我们演示使用 expvar 包可以做什么。以下代码创建一个在监听 8080 端口的 HTTP 服务器。每个请求到达 hander\(\) 后，在向访问者发送响应消息之前增加计数器。

```text
    package main
    
    import (
        "expvar"
        "fmt"
        "net/http"
    )
    
    var visits = expvar.NewInt("visits")
    
    func handler(w http.ResponseWriter, r *http.Request) {
        visits.Add(1)
        fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
    }
    
    func main() {
        http.HandleFunc("/", handler)
        http.ListenAndServe(":8080", nil)
    }
```

导入 expvar 包后，它将为 `http.DefaultServeMux` 上的 PATH `/debug/vars` 注册一个处理函数。此处理程序返回已在 expvar 包中注册的所有公共变量。运行代码并访问 `http://localhost:8080/debug/vars`，您将看到如下所示的内容（输出被截断以增加可读性）：

```text
{
  "cmdline": [
    "/var/folders/qv/2jztyc09357ddtxn_bvgh8j00000gn/T/go-build146580631/command-line-arguments/_obj/exe/test"
  ],
  "memstats": {
    "Alloc": 414432,
    "TotalAlloc": 414432,
    "Sys": 3084288,
    "Lookups": 13,
    "Mallocs": 5111,
    "Frees": 147,
    "HeapAlloc": 414432,
    "HeapSys": 1703936,
    "HeapIdle": 835584,
    "HeapInuse": 868352,
    "HeapReleased": 0,
    "HeapObjects": 4964,
    "StackInuse": 393216,
    "StackSys": 393216,
    "MSpanInuse": 15504,
    "MSpanSys": 16384,
    "MCacheInuse": 4800,
    "MCacheSys": 16384,
    "BuckHashSys": 2426,
    "GCSys": 137216,
    "OtherSys": 814726,
    "NextGC": 4473924,
    "LastGC": 0,
    "PauseTotalNs": 0,
    "PauseNs": [
      0,
      0,
    ],
    "PauseEnd": [
      0,
      0
    ],
    "GCCPUFraction": 0,
    "EnableGC": true,
    "DebugGC": false,
    "BySize": [
      {
        "Size": 16640,
        "Mallocs": 0,
        "Frees": 0
      },
      {
        "Size": 17664,
        "Mallocs": 0,
        "Frees": 0
      }
    ]
  },
  "visits": 0
}
```

信息真不少。这是因为默认情况下该包注册了 `os.Args` 和 `runtime.Memstats` 两个指标。因为我们还没有访问到增加 visits 的路径，所以它的值仍然为 0。现在通过访问 `http:// localhost:8080/golang` 来增加计数器，然后返回。计数器不再为 0。

### expvar.Publish 函数

expvar 包相当小且容易理解。它主要由两个部分组成。第一个是函数 `expvar.Publish(name string，v expvar.Var)`。该函数可用于在未导出的全局注册表中注册具有特定名称（name）的 v。以下代码段显示了具体实现。接下来的 3 个代码段是从 expvar 包的源代码中截取的。

先看下全局注册表：

```text
	var (
		mutex   sync.RWMutex
		vars    = make(map[string]Var)
		varKeys []string // sorted
	)
```

全局注册表实际就是一个 map：vars。

Publish 函数的实现：

```text
    // Publish declares a named exported variable. This should be called from a
    // package's init function when it creates its Vars. If the name is already
    // registered then this will log.Panic.
    func Publish(name string, v Var) {
        mutex.Lock()
        defer mutex.Unlock()
    
        // Check if name has been taken already. If so, panic.
        if _, existing := vars[name]; existing {
            log.Panicln("Reuse of exported var name:", name)
        }
    
         // vars is the global registry. It is defined somewhere else in the
         // expvar package like this:
         //
         //  vars = make(map[string]Var)
        vars[name] = v
        // 一方面，该包中所有公共变量，放在 vars 中，同时，通过 varKeys 保存了所有变量名，并且按字母序排序，即实现了一个有序的、线程安全的 map
        varKeys = append(varKeys, name)
        sort.Strings(varKeys)
    }
```

expvar 包内置的两个公共变量就是通过 Publish 注册的：

```text
Publish("cmdline", Func(cmdline))
Publish("memstats", Func(memstats))
```

### expvar.Var 接口

另一个重要的组成部分是 `expvar.Var` 接口。 这个接口只有一个方法：

```text
    // Var is an abstract type for all exported variables.
    type Var interface {
            // String returns a valid JSON value for the variable.
            // Types with String methods that do not return valid JSON
            // (such as time.Time) must not be used as a Var.
            String() string
    }
```

所以你可以在有 String\(\) string 方法的所有类型上调用 Publish\(\) 函数，但需要注意的是，这里的 String\(\) 要求返回的是一个有效的 JSON 字符串。

### expvar.Int 类型

expvar 包提供了其他几个类型，它们实现了 expvar.Var 接口。其中一个是 expvar.Int，我们已经在演示代码中通过 expvar.NewInt\("visits"\) 方式使用它了，它会创建一个新的 expvar.Int，并使用 expvar.Publish 注册它，然后返回一个指向新创建的 expvar.Int 的指针。

```text
    func NewInt(name string) *Int {
        v := new(Int)
        Publish(name, v)
        return v
    }
```

expvar.Int 包装一个 int64，并有两个函数 Add\(delta int64\) 和 Set\(value int64\)，它们以线程安全的方式（通过 `atomic` 包实现）修改包装的 int64。另外通过 `Value() int64` 函数获得包装的 int64。

```text
	type Int struct {
		i int64
	}
```

### 其他类型

除了 expvar.Int，该包还提供了一些实现 expvar.Var 接口的其他类型：

* [expvar.Float](http://docs.studygolang.com/pkg/expvar/#Float)
* [expvar.String](http://docs.studygolang.com/pkg/expvar/#String)
* [expvar.Map](http://docs.studygolang.com/pkg/expvar/#Map)
* [expvar.Func](http://docs.studygolang.com/pkg/expvar/#Func)

前两个类型包装了 float64 和 string。后两种类型需要稍微解释下。

`expvar.Map` 类型可用于使公共变量出现在某个名称空间下。可以这样用：

```text
    var stats = expvar.NewMap("http")
    var requests, requestsFailed expvar.Int
    
    func init() {
        stats.Set("req_succ", &requests)
        stats.Set("req_failed", &requestsFailed)
    }
```

这段代码使用名称空间 http 注册了两个指标 req\_succ 和 req\_failed。它将显示在 JSON 响应中，如下所示：

```text
{
    "http": {
        "req_succ": 18,
        "req_failed": 21
    }
}
```

当要注册某个函数的执行结果到某个公共变量时，您可以使用 `expvar.Func`。假设您希望计算应用程序的正常运行时间，每次有人访问 `http://localhost:8080/debug/vars` 时，都必须重新计算此值。

```text
    var start = time.Now()
    
    func calculateUptime() interface {
        return time.Since(start).String()
    }
    
    expvar.Publish("uptime", expvar.Func(calculateUptime))
```

实际上，内置的两个指标 `cmdline` 和 `memstats` 就是通过这种方式注册的。注意，函数签名有如下要求：没有参数，返回 interface{}

```text
type Func func() interface{}
```

### 关于 Handler 函数

本文开始时提到，可以简单的导入 expvar 包，然后使用其副作用，导出路径 `/debug/vars`。然而，如果我们使用了一些框架，并非使用 `http.DefaultServeMux`，而是框架自己定义的 `Mux`，这时直接导入使用副作用可能不会生效。我们可以按照使用的框架，定义自己的路径，比如也叫 `/debug/vars`，然后，这对应的处理程序中，按如下的两种方式处理：

1）将处理直接交给 `expvar.Handler`，比如：

```text
handler := expvar.Handler()
handler.ServeHTTP(w, req)
```

2）自己遍历 expvar 中的公共变量，构造输出，甚至可以过滤 expvar 默认提供的 cmdline 和 memstats，我们看下 expvarHandler 的源码就明白了：（通过 expvar.Do 函数来遍历的）

```text
    func expvarHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        fmt.Fprintf(w, "{\n")
        first := true
        Do(func(kv KeyValue) {
            if !first {
                fmt.Fprintf(w, ",\n")
            }
            first = false
            fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
        })
        fmt.Fprintf(w, "\n}\n")
    }
```

[Go 语言中文网](https://github.com/studygolang/studygolang/blob/V3.0/src/http/controller/admin/metrics.go#L42) 因为使用了 Echo 框架，使用第 1 种方式来处理的。

### 定义自己的 expvar.Var 类型

expvar 包提供了 int、float 和 string 这三种基本数据类型的 expvar.Var 实现，以及 Func 和 Map。有时，我们自己有一个复杂的类型，想要实现 expvar.Var 接口，怎么做呢？

从上面的介绍，应该很容易实现吧，如果您遇到了具体的需求，可以试试。

### 总结

综上所述，通过 expvar 包，可以非常容易的展示应用程序指标。建议您在您的每个应用程序中使用它来展示一些指示应用程序运行状况的指标，通过它和其他的一些工具来监控应用程序。

