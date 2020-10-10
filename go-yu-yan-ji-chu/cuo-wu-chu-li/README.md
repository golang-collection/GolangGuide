# 错误与异常处理

转载自：[https://github.com/unknwon/the-way-to-go\_ZH\_CN/blob/master/eBook/13.0.md](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/13.0.md) by unknwon

## 异常与错误

Go 没有像 Java 和 .NET 那样的 `try/catch` 异常机制：不能执行抛异常操作。但是有一套 `defer-panic-and-recover` 机制。

Go 的设计者觉得 `try/catch` 机制的使用太泛滥了，而且从底层向更高的层级抛异常太耗费资源。他们给 Go 设计的机制也可以 “捕捉” 异常，但是更轻量，并且只应该作为（处理错误的）最后的手段。

Go 是怎么处理普通错误的呢？通过在函数和方法中返回错误对象作为它们的唯一或最后一个返回值——如果返回 nil，则没有错误发生——并且主调（calling）函数总是应该检查收到的错误。

**永远不要忽略错误，否则可能会导致程序崩溃！！**

处理错误并且在函数发生错误的地方给用户返回错误信息：照这样处理就算真的出了问题，你的程序也能继续运行并且通知给用户。`panic and recover` 是用来处理真正的异常（无法预测的错误）而不是普通的错误。

Go 检查和报告错误条件的惯有方式：

* 产生错误的函数会返回两个变量，一个值和一个错误码；如果后者是 nil 就是成功，非 nil 就是发生了错误。
* 为了防止发生错误时正在执行的函数（如果有必要的话甚至会是整个程序）被中止，在调用函数后必须检查错误。

下面这段来自 pack1 包的代码 Func1 测试了它的返回值：

```go
if value, err := pack1.Func1(param1); err != nil {
	fmt.Printf("Error %s in pack1.Func1 with parameter %v", err.Error(), param1)
	return    // or: return err
} else {
	// Process(value)
}
```

_为了更清晰的代码，应该总是使用包含错误值变量的 if 复合语句_

上例除了 `fmt.Printf` 还可以使用日志输出，如果程序中止也没关系的话甚至可以使用 `panic`。

## 推荐阅读

{% embed url="https://mp.weixin.qq.com/s/KPrzPP797efFUKOTTfY1Ow" %}

{% embed url="https://blog.golang.org/defer-panic-and-recover" %}



