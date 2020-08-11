# 惰性求值

转载自：[https://github.com/unknwon/the-way-to-go\_ZH\_CN/blob/master/eBook/06.6.md](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/06.6.md)

生成器是指当被调用时返回一个序列中下一个值的函数，例如：

```text
    generateInteger() => 0
    generateInteger() => 1
    generateInteger() => 2
    ....
```

生成器每次返回的是序列中下一个值而非整个序列；这种特性也称之为惰性求值：只在你需要时进行求值，同时保留相关变量资源（内存和cpu）：这是一项在需要时对表达式进行求值的技术。例如，生成一个无限数量的偶数序列：要产生这样一个序列并且在一个一个的使用可能会很困难，而且内存会溢出！但是一个含有通道和go协程的函数能轻易实现这个需求。

在14.12的例子中，我们实现了一个使用 int 型通道来实现的生成器。通道被命名为`yield`和`resume`，这些词经常在协程代码中使用。

示例 14.12 [lazy\_evaluation.go](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/examples/chapter_14/lazy_evaluation.go)：

```text
package main

import (
    "fmt"
)

var resume chan int

func integers() chan int {
    yield := make(chan int)
    count := 0
    go func() {
        for {
            yield <- count
            count++
        }
    }()
    return yield
}

func generateInteger() int {
    return <-resume
}

func main() {
    resume = integers()
    fmt.Println(generateInteger())  //=> 0
    fmt.Println(generateInteger())  //=> 1
    fmt.Println(generateInteger())  //=> 2    
}
```

有一个细微的区别是从通道读取的值可能会是稍早前产生的，并不是在程序被调用时生成的。如果确实需要这样的行为，就得实现一个请求响应机制。当生成器生成数据的过程是计算密集型且各个结果的顺序并不重要时，那么就可以将生成器放入到go协程实现并行化。但是得小心，使用大量的go协程的开销可能会超过带来的性能增益。

这些原则可以概括为：通过巧妙地使用空接口、闭包和高阶函数，我们能实现一个通用的惰性生产器的工厂函数`BuildLazyEvaluator`（这个应该放在一个工具包中实现）。工厂函数需要一个函数和一个初始状态作为输入参数，返回一个无参、返回值是生成序列的函数。传入的函数需要计算出下一个返回值以及下一个状态参数。在工厂函数中，创建一个通道和无限循环的go协程。返回值被放到了该通道中，返回函数稍后被调用时从该通道中取得该返回值。每当取得一个值时，下一个值即被计算。在下面的例子中，定义了一个`evenFunc`函数，其是一个惰性生成函数：在main函数中，我们创建了前10个偶数，每个都是通过调用`even()`函数取得下一个值的。为此，我们需要在`BuildLazyIntEvaluator`函数中具体化我们的生成函数，然后我们能够基于此做出定义。

示例 14.13 [general\_lazy\_evalution1.go](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/examples/chapter_14/general_lazy_evalution1.go)：

```text
package main

import (
    "fmt"
)

type Any interface{}
type EvalFunc func(Any) (Any, Any)

func main() {
    evenFunc := func(state Any) (Any, Any) {
        os := state.(int)
        ns := os + 2
        return os, ns
    }
    
    even := BuildLazyIntEvaluator(evenFunc, 0)
    
    for i := 0; i < 10; i++ {
        fmt.Printf("%vth even: %v\n", i, even())
    }
}

func BuildLazyEvaluator(evalFunc EvalFunc, initState Any) func() Any {
    retValChan := make(chan Any)
    loopFunc := func() {
        var actState Any = initState
        var retVal Any
        for {
            retVal, actState = evalFunc(actState)
            retValChan <- retVal
        }
    }
    retFunc := func() Any {
        return <- retValChan
    }
    go loopFunc()
    return retFunc
}

func BuildLazyIntEvaluator(evalFunc EvalFunc, initState Any) func() int {
    ef := BuildLazyEvaluator(evalFunc, initState)
    return func() int {
        return ef().(int)
    }
}
```

输出：

```text
0th even: 0
1th even: 2
2th even: 4
3th even: 6
4th even: 8
5th even: 10
6th even: 12
7th even: 14
8th even: 16
9th even: 18
```

练习14.12：[general\_lazy\_evaluation2.go](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/exercises/chapter_14/general_lazy_evalution2.go) 通过使用14.12中工厂函数生成前10个斐波那契数

提示：因为斐波那契数增长很迅速，使用`uint64`类型。 注：这种计算通常被定义为递归函数，但是在没有尾递归的语言中，例如go语言，这可能会导致栈溢出，但随着go语言中堆栈可扩展的优化，这个问题就不那么严重。这里的诀窍是使用了惰性求值。gccgo编译器在某些情况下会实现尾递归。

## 链接

* [目录](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/directory.md)
* 上一节：[新旧模型对比：任务和worker](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/14.7.md)
* 下一节：[实现 Futures 模式](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/14.9.md)

