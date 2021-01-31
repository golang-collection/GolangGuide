# 变量的逃逸分析

简单来说，只要编译器发现该变量具备逃逸性质就会将其分配到堆上。

## 如何做逃逸分析

```go
package main

import "fmt"

func foo() *int {
	t := 3
	return &t;
}

func main() {
	x := foo()
	fmt.Println(*x)
}
```

foo函数返回一个局部变量的指针，main函数里变量x接收它。执行如下命令：

```bash
go build -gcflags '-m -l' main.go
```

加`-l`是为了不让foo函数被内联。得到如下输出：

{% hint style="info" %}
内联函数：内联函数是c语言的一个函数，其目的是防止多次调用函数导致函数进栈消耗资源，最终使程序崩溃，foo函数虽然被定义为单独的函数，但是在执行时可以将x看作是直接执行了t := 3的操作。
{% endhint %}

```bash
# command-line-arguments
src/main.go:7:9: &t escapes to heap
src/main.go:6:7: moved to heap: t
src/main.go:12:14: *x escapes to heap
src/main.go:12:13: main ... argument does not escape
```

foo函数里的变量`t`逃逸了，和我们预想的一致。让我们不解的是为什么main函数里的`x`也逃逸了？这是因为有些函数参数为interface类型，比如fmt.Println\(a ...interface{}\)，编译期间很难确定其参数的具体类型，也会发生逃逸。

## 推荐阅读

{% embed url="https://mp.weixin.qq.com/s/Qvm3ARplkZQSNERtm3TXHA" %}

{% embed url="https://www.kancloud.cn/aceld/golang/1958306" %}

{% embed url="https://www.cnblogs.com/qcrao-2018/p/10453260.html" %}



