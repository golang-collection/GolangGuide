# 闭包

首发在：[https://blog.csdn.net/s\_842499467/article/details/104281602](https://blog.csdn.net/s_842499467/article/details/104281602)

关于什么是闭包，我相信大家初次听到这个概念的时候一定是非常迷茫，我也是查找了很多资料也还是没有那么透彻的理解，但是在之后的实践中通过逐渐的使用好像更了解了闭包是什么，所以希望通过此篇文章来和大家分享一下我对闭包的理解。

## 什么是闭包

我认为的闭包就是**一个函数与这个函数外部变量的一个封装**。粗略的可以理解为一个类，类里面有变量和方法，**其中闭包所包含的外部变量对应着类中的静态变量。** 为什么这么理解，首先让我们来看一个例子。

```go
func add() func(int) int {
	n := 10
	str := "string"
	return func(x int) int {
		n = n + x
		str += strconv.Itoa(x)
		fmt.Println(str)
		return n
	}
}

func main() {
	f := add()
	fmt.Println(f(1))
	fmt.Println(f(2))
	fmt.Println(f(3))
}
```

这个程序的输出结果是

```go
string1
11
string12
13
string123
16
```

如果不了解的闭包的朋友肯定会觉得很奇怪，为什么会输出这样的结果。这就要用到我最开始的解释。**闭包就是一个函数和一个函数外的变量的封装，而且这个变量就对应着类中的静态变量。** 这样就可以将这个程序的输出结果解释的通了。

* 最开始我们先声明一个函数`add`，在函数体内返回一个匿名函数  ![&#x95ED;&#x5305;](https://img-blog.csdnimg.cn/20200212175537911.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3NfODQyNDk5NDY3,size_16,color_FFFFFF,t_70)
* 其中的`n`,`str`与下面的匿名函数构成了整个的闭包，`n`与`str`就像类中的静态变量只会初始化一次，所以说尽管后面多次调用这个整体函数，里面都不会再重新初始化了
* 而且对于外部变量的操作是累加的，这与类中的静态变量也是一致的

在go语言学习笔记中，雨痕提到在汇编代码中，闭包返回的不仅仅是匿名函数，还包括所引用的环境变量指针，这与我们之前的解释也是类似的，闭包通过操作指针来调用对应的变量。

通过以上的解释我们再来理解一下开头的那个程序，就能解释的通了。

## 延迟求值

```go
func main() {
	for _, f := range testClouser(){
		f()
	}
}

func testClouser() []func(){
	var s []func()

	for i := 0; i<2; i++{
		s = append(s, func() {
			fmt.Println(&i, i)
		})
	}

	return s
}
```

根据上面的解释，大家可以猜一猜结果是什么？

```go
0xc00009c008 2
0xc00009c008 2
```

## 闭包的好处

对于一个函数来说，如果不调用闭包也能实现绝大部分功能，但是为什么要使用闭包呢，我认为有以下原因

* 通过封装的方式使代码更简洁，减少全局变量的使用
* 通过闭包我们可以减少参数的传递

我们可以通过这个例子来对比一下，这个函数是判断一个字符串的后缀是否为`.jpg`，如果不是的话，则在字符串的末尾加上`.jpg`，否则返回原值，我们分别使用闭包和普通函数实现以下，我们来看一下区别。

```go

func makeSuffix(suffix string) func (string) string{
	return func (name string) string{
		if !strings.HasSuffix(name, suffix){
			return name + suffix
		}
		return name
	}
}


func makeSuffix2(name , suffix string) string {
	if !strings.HasSuffix(name, suffix){
		return name + suffix
	}
	return name
}

func main() {
	f := makeSuffix(".jpg")
	fmt.Println(f("winter"))
	fmt.Println(f("hello.jpg"))
	
    fmt.Println(makeSuffix2("winter", ".jpg"))
	fmt.Println(makeSuffix2("hello.jpg", ".jpg"))
}
```

通过这个例子我想大家能够更好的理解闭包的好处了。

## 底层原理

再介绍闭包底层原理之前再强调一下，在go语言中函数是一等公民，所以你可能会看见下面这样的代码

```go
func B(f func()) {}
func c() func() {return func(){}}
var f func() = c()
```

但在实现层面每一个函数在执行时会通过runtime.FuncVal结构体来实保存对应的地址

```go
type funcval struct {
    fn uintptr
}
```

其中fn会保存函数指令的入口地址，下面这个例子中演示了正常情况下函数的执行过程

![](../../.gitbook/assets/image%20%2847%29.png)

在这种情况下，编译器会作出优化，让f1和f2共用一个funcval结构体，然后在编译阶段，会在只读数据段中分配一个funcval结构体对应fn来记录func A\(\)的入口地址，而其本身的地址addr2会在执行阶段赋值给f1和f2。

下面来看一个闭包的案例

![](../../.gitbook/assets/image%20%2852%29.png)

在函数执行过程中，会在内存中按照如下方式分配内存空间

![](../../.gitbook/assets/image%20%2850%29.png)

继续向下执行，就会将addr2的地址返回给main栈的返回值处，然后赋值给f1。

![](../../.gitbook/assets/image%20%2854%29.png)

函数继续向下执行就会将addr3的地址赋值给f2，接下来通过f1和f2调用函数时会调用对应的addr1，但是f1和f2会使用各自的捕获列表。

在go语言中会将funcval结构体的地址存储到特定的寄存器中，通过funcval结构体地址+对应偏移就可以得到对应的捕获列表。

不仅如此，如果在闭包中还会修改对应的变量，那么go编译器首先会将该变量分配到堆上，然后栈与对应的funcval结构体位置存储该变量的地址。

## 总结

总的来说，在没有类这种概念的情况下，golang通过使用闭包的概念仍然可以实现一些简单的封装，使得代码更加简洁，使用起来也更加方便，但是闭包可能会存在占用内存等问题，所以在使用的过程中也要注意。

## **推荐阅读**

{% embed url="https://www.bilibili.com/video/BV1ma4y1e7R5" %}

\*\*\*\*

