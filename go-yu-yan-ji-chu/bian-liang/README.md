# 变量

## 变量

go语言是一种静态类型的语言，也就意味着如果一个变量被声明了，那么我们只能修改它的值，而不能修改它的类型。不像python语言那样，变量的值和类型是在运行时确定的。

## 变量定义

go语言提供了多种变量的定义方式

```go
var a int
var b = false
var x, y = 100, "test"
// 等同于
var(
    x = 100
    y = "test"
)

//简短方式
a := 1
```

以上是几种go语言中变量的定义方式。在go语言中变量会自动赋予初始值，而且变量声明之后不使用会报错。

在变量的简短定义方式中，变量的类型go语言会自动推导，但是像例子中的a只能被推导为int类型，要是需要将其转换为float类型需要强制转换。

```go
b := float64(a)
```

在go语言中还有一种特性比较有趣

```go
func main() {
	x := 100
	fmt.Println(&x, x)

	x, y := 200, "abc"
	fmt.Println(&x, x)
	fmt.Println(&y, y)
}
```

在上面这段代码中我们首先定义了一个变量x，在之后我们好像又定义了一个x同时定义了一个y，但是从程序的输出我们可以看到，在下面的那个变量定义中，x只是被赋予了新的值，而不是重新定义，也叫做退化赋值。

```go
0xc00001e0b8 100
0xc00001e0b8 200
0xc000010200 abc
```

而且需要注意的是，上面这中退化赋值操作只能在同一作用域下起作用，像下面这样就不可以了。

```go
func main() {
	x := 100
	fmt.Println(&x, x)
	{
		x, y := 200, "abc"
		fmt.Println(&x, x)
		fmt.Println(&y, y)
	}
}
```

程序的输出为：

```go
0xc00001e0b8 100
0xc00001e0c8 200
0xc000010200 abc
```

上面的退化操作在go语言中允许我们重复使用err来代表错误，以免生成非常多的变量。

## 空标识符

在go语言中有一个特殊的成员`_`，它可以表示空占位符，比如有的函数有两个返回值，但我只想使用其中的一个，我们就可以用下面这种方式使用。

```go
x, _ = strconv.Atoi("12")
```

此时Atoi函数的error返回值就被我们忽略了。

## 变量交换

在go语言中交换两个变量非常的方便，可以使用如下的方法。

```go
var a int = 100
var b int = 200
b, a = a, b
fmt.Println(a, b)

//输出
200 100
```

在上面的多重赋值时，变量的左值和右值按从左到右的顺序赋值。

## 变量作用域

### 局部变量

在函数体内声明的变量称之为局部变量，它们的作用域只在函数体内，函数的参数和返回值变量都属于局部变量。局部变量不是一直存在的，它只在定义它的函数被调用后存在，函数调用结束后这个局部变量就会被销毁。

### 全局变量

全局变量只需要在一个文件中定义，就可以在所有文件中使用，当然，不包含这个全局变量的文件需要使用“import”关键字引入全局变量所在的源文件之后才能使用这个全局变量。全局变量声明必须以 var 关键字开头，如果想要在外部包中使用全局变量的首字母必须大写。

```go
var g = 2.3

func main() {
	fmt.Println(g)
	i := 5
	local(i)
}

func local(i int){
	i = 10
	fmt.Println(i)
	fmt.Println(g)
}
```

下面的例子中也说明了变量的作用域问题

```go
func main() {
	b := 2
	fmt.Println(&b)
	{
		a, b := 1, 3
		fmt.Println(&a, &b)
	}

	if a, b := 5, 7; a < 6 || b < 6 {
		fmt.Println(&a, &b)
	}
}
```

程序输出

```go
0xc00009a008
0xc00009a020 0xc00009a028
0xc00009a030 0xc00009a038
```

所以函数中常用的变量推荐使用这种方式统一进行声明，不仅能减少变量的创建还能防止出现一些可能[被忽略的问题](../../go-yu-yan-chang-jian-keng/if-fu-zhi-yu-ju.md)。

```go
func main() {
	var (
		a int
		b int
	)
	fmt.Println(&a, &b)
	{
		a, b = 1, 3
		fmt.Println(&a, &b)
	}

	if a, b = 5, 7; a < 6 || b < 6 {
		fmt.Println(&a, &b)
	}
}
```
