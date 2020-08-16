# 空接口

转载自：[http://c.biancheng.net/view/84.html](http://c.biancheng.net/view/84.html)

空接口是接口类型的特殊形式，空接口没有任何方法，因此任何类型都无须实现空接口。从实现的角度看，任何值都满足这个接口的需求。因此空接口类型可以保存任何值，也可以从空接口中取出原值。

{% hint style="info" %}
空接口类型类似于 [C\#](http://c.biancheng.net/csharp/) 或 [Java](http://c.biancheng.net/java/) 语言中的 Object、C语言中的 void\*、[C++](http://c.biancheng.net/cplus/) 中的 std::any。在泛型和模板出现前，空接口是一种非常灵活的数据抽象保存和使用的方法。
{% endhint %}

空接口的内部实现保存了对象的类型和指针。使用空接口保存一个数据的过程会比直接用数据对应类型的变量保存稍慢。因此在开发中，应在需要的地方使用空接口，而不是在所有地方使用空接口。

##  将值保存到空接口

 空接口的赋值如下：

```text
var any interface{}

any = 1
fmt.Println(any)

any = "hello"
fmt.Println(any)

any = false
fmt.Println(any)
```

 代码输出如下：

```text
1
hello
false
```

 对代码的说明：

*  第 1 行，声明 any 为 interface{} 类型的变量。
*  第 3 行，为 any 赋值一个整型 1。
*  第 4 行，打印 any 的值，提供给 fmt.Println 的类型依然是 interface{}。
*  第 6 行，为 any 赋值一个字符串 hello。此时 any 内部保存了一个字符串。但类型依然是 interface{}。
*  第 9 行，赋值布尔值。

##  从空接口获取值

 保存到空接口的值，如果直接取出指定类型的值时，会发生编译错误，代码如下：

```text
// 声明a变量, 类型int, 初始值为1
var a int = 1

// 声明i变量, 类型为interface{}, 初始值为a, 此时i的值变为1
var i interface{} = a

// 声明b变量, 尝试赋值i
var b int = i
```

 第8行代码编译报错：

 cannot use i \(type interface {}\) as type int in assignment: need type assertion

 编译器告诉我们，不能将i变量视为int类型赋值给b。

 在代码第 15 行中，将 a 的值赋值给 i 时，虽然 i 在赋值完成后的内部值为 int，但 i 还是一个 interface{} 类型的变量。类似于无论集装箱装的是茶叶还是烟草，集装箱依然是金属做的，不会因为所装物的类型改变而改变。

 为了让第 8 行的操作能够完成，编译器提示我们得使用 type assertion，意思就是类型断言。

 使用类型断言修改第 8 行代码如下：

```text
var b int = i.(int)
```

 修改后，代码可以编译通过，并且 b 可以获得 i 变量保存的 a 变量的值：1。

##  空接口的值比较

 空接口在保存不同的值后，可以和其他变量值一样使用`==`进行比较操作。空接口的比较有以下几种特性。

###  1\) 类型不同的空接口间的比较结果不相同

 保存有类型不同的值的空接口进行比较时，Go语言会优先比较值的类型。因此类型不同，比较结果也是不相同的，代码如下：

```text
// a保存整型
var a interface{} = 100

// b保存字符串
var b interface{} = "hi"

// 两个空接口不相等
fmt.Println(a == b)
```

 代码输出如下：

 false

###  2\) 不能比较空接口中的动态值

 当接口中保存有动态类型的值时，运行时将触发错误，代码如下：

```text
// c保存包含10的整型切片
var c interface{} = []int{10}

// d保存包含20的整型切片
var d interface{} = []int{20}

// 这里会发生崩溃
fmt.Println(c == d)
```

 代码运行到第8行时发生崩溃：

 panic: runtime error: comparing uncomparable type \[\]int

 这是一个运行时错误，提示 \[\]int 是不可比较的类型。下表中列举出了类型及比较的几种情况。

|   |  |
| :--- | :--- |
|  类  型 |  说  明 |
|  map |  宕机错误，不可比较 |
|  切片（\[\]T） |  宕机错误，不可比较 |
|  通道（channel） |  可比较，必须由同一个 make 生成，也就是同一个通道才会是 true，否则为 false |
|  数组（\[容量\]T） |  可比较，编译期知道两个数组是否一致 |
|  结构体 |  可比较，可以逐个比较结构体的值 |
|  函数 |  可比较 |

