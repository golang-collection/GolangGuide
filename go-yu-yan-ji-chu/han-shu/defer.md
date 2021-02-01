# defer延迟调用

首发于：[https://blog.csdn.net/s\_842499467/article/details/104283109](https://blog.csdn.net/s_842499467/article/details/104283109)

## 1. 为什么需要defer机制

在项目中，我们常用的操作就是释放资源等操作，为了即时的释放资源，go设计了`defer`机制。defer语句常用语资源释放，解除锁定以及错误处理等。

## 2. defer机制的简单实用

我们先通过一个案例来看一下`defer`机制的作用

```go
import "fmt"

func defer_sum(num1, num2 int) int {
	defer fmt.Println("defer1 : ", num1)
	defer fmt.Println("defer2 : ", num2)

	res := num1 + num2
	fmt.Println("sum res : ", res)
	return res
}

func main() {
	res := defer_sum(10 ,20)
	fmt.Println("main : ", res)
}
```

程序的输出结果如下

```go
sum res :  30
defer2 :  20
defer1 :  10
main :  30
```

我们可以看到通过`defer`声明的操作不会立即执行，而是当所处的代码块执行完成之后才继续执行。

## 3. defer机制的细节

那么defer机制为什么会产生这样的结果呢，是因为通过defer声明的语句不会立刻执行，**而会压入一个`defer`栈中**，因为栈遵循的是先入后出原则，所以得到如上所示的结果。

**在 `defer`将语句放入到栈时，也会将相关的值拷贝同时入栈。** 我们用一段代码来解释这个特点。这段代码与刚才的那段代码没什么不同，只是在中间加入了一段`num1`和`num2`自增和自减的操作。

```go
import "fmt"

func defer_sum(num1, num2 int) int {
	defer fmt.Println("defer1 : ", num1)
	defer fmt.Println("defer2 : ", num2)
	
	num1++
	num2--

	res := num1 + num2
	fmt.Println("sum res : ", res)
	return res
}

func main() {
	res := defer_sum(10 ,20)
	fmt.Println("main : ", res)
}
```

程序的输出结果如下

```text
sum res :  30
defer2 :  20
defer1 :  10
main :  30
```

可以看到输出结果与上一个案例的输出结果相同，这也就解释了在**`defer`语句入栈的同时也会将值拷贝一份压入栈中。**

## 4. defer与闭包结合

对闭包不熟悉的可以参看[**闭包**](bi-bao.md)。

```go
func main() {
	for i := 0;i<3 ;i++  {
		defer func() {
			fmt.Println(i)
		}()
	}
}
```

程序的输出结果为

```go
3
3
3
```

我在[Go-闭包](https://blog.csdn.net/s_842499467/article/details/104281602)这篇博客中讲过，我们可以**将闭包中的变量看作是类中的静态变量**，再结合`defer`机制的性质，得出这样的结果也就不足为奇了。

## 5. 什么时候使用defer

就像开头提到的，defer机制就是为了更好的关闭资源的，所以我们使用`defer`也是在创建资源后使用，如下例所示。

```go
func main(){
    connect := connectDB()
    defer connect.close()
    
}
```

## 6. 注意事项

需要注意一点的是如果我们在main函数中申请资源时使用了defer，要注意这个资源是main函数执行完才会被释放，如果申请的资源很大那无疑是一种错误的处理方式，更为优雅的方式是将其与匿名函数封装，这样匿名函数调用结束则释放资源。

而且defer语句也要花费更大的代价，所以在高性能的算法设计中要谨慎使用。

## 推荐阅读

{% embed url="https://www.kancloud.cn/aceld/golang/1958310" %}



