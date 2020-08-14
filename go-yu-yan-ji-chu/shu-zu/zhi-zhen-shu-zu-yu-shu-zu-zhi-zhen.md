# 指针数组与数组指针

对于指针数组和数组指针在c或c++中也经常被讨论，尤其对于初学者来说会分辨不清楚。其实在每个词中间添加一个“的“就很好理解了，指针的数组，数组的指针。

## 指针数组

对于指针数组来说，就是：一个数组里面装的都是指针，在go语言中数组默认是值传递的，所以如果我们在函数中修改传递过来的数组对原来的数组是没有影响的。

```go
func main() {
	var a [5]int
	fmt.Println(a)
	test(a)
	fmt.Println(a)
}

func test(a [5]int) {
	a[1] = 2
	fmt.Println(a)
}
```

输出

```go
[0 0 0 0 0]
[0 2 0 0 0]
[0 0 0 0 0]
```

但是如果我们一个函数传递的是指针数组，情况会有什么不一样呢？

```go
func main() {
	var a [5]*int
	fmt.Println(a)
	for i := 0; i < 5; i++ {
		temp := i
		a[i] = &temp
	}
	for i := 0; i < 5; i++ {
		fmt.Print(" ", *a[i])
	}
	fmt.Println()
	test1(a)
	for i := 0; i < 5; i++ {
		fmt.Print(" ", *a[i])
	}
}

func test1(a [5]*int) {
	*a[1] = 2
	for i := 0; i < 5; i++ {
		fmt.Print(" ", *a[i])
	}
	fmt.Println()
}
```

我们先来看一下程序的运行结果

```go
[<nil> <nil> <nil> <nil> <nil>]
 0 1 2 3 4
 0 2 2 3 4
 0 2 2 3 4
```

可以看到初始化值全是nil，也就验证了指针数组内部全都是一个一个指针，之后我们将其初始化，内部的每个指针指向一块内存空间。

{% hint style="info" %}
在初始化的时候如果直接另a\[i\] = &i那么指针数组内部存储的全是同一个地址，所以输出结果也一定是相同的
{% endhint %}

然后我们将这个指针数组传递给test1函数，**对于数组的参数传递仍然是复制的形式也就是值传递**，但是因为数组中每个元素是一个指针，所以test1函数复制的新数组中的值仍然是这些指针指向的具体地址值，这时改变a\[1\]这块存储空间地址指向的值，那么原实参指向的值也会变为2，具体流程如下图所示。

![](../../.gitbook/assets/image%20%287%29.png)

## 数组指针

了解了指针数组之后，再来看一下数组指针，数组指针的全称应该叫做，指向数组的指针，在go语言中我们可以如下操作。

```go
func main() {
	var a [5]int
	var aPtr *[5]int
	aPtr = &a
	//这样简短定义也可以aPtr := &a
	fmt.Println(aPtr)
	test(aPtr)
	fmt.Println(aPtr)
}

func test(aPtr *[5]int) {
	aPtr[1] = 5
	fmt.Println(aPtr)
}
```

我们先定义了一个数组a，然后定一个指向数组a的指针aPtr，然后将这个指针传入一个函数，在函数中我们改变了具体的值，程序的输出结果

```go
&[0 0 0 0 0] 
&[0 5 0 0 0] 
&[0 5 0 0 0]
```

![](../../.gitbook/assets/image%20%288%29.png)

通过上面的图我们可以看见虽然main和test函数中的aPtr是不同的指针，但是他们都指向同一块数组的内存空间，所以无论在main函数还是在test函数中对数组的操作都会直接改变数组。

