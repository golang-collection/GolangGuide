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

## 总结

总的来说，在没有类这种概念的情况下，golang通过使用闭包的概念仍然可以实现一些简单的封装，使得代码更加简洁，使用起来也更加方便，但是闭包可能会存在占用内存等问题，所以在使用的过程中也要注意。

**我也找到了js中闭包的一个优秀解答**[**链接在这里**](https://stackoverflow.com/questions/111102/how-do-javascript-closures-work)

