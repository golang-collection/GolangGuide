# 常量与枚举

## 常量

常量表示运行时不可改变的一些值，我们使用常量的标识符来代替一些魔法数字，使得在调整程序中的常量值时不需要修改所有引用的代码。

**转载自：**[**https://github.com/unknwon/the-way-to-go\_ZH\_CN/blob/master/eBook/04.3.md**](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/04.3.md)\*\*\*\*

常量使用关键字 `const` 定义，用于存储不会改变的数据。存储在常量中的数据类型只可以是布尔型、数字型（整数型、浮点型和复数）和字符串型。

常量的定义格式：`const identifier [type] = value`，例如：

```go
const Pi = 3.14159
```

在 Go 语言中，你可以省略类型说明符 `[type]`，因为编译器可以根据变量的值来推断其类型。

* 显式类型定义： `const b string = "abc"`
* 隐式类型定义： `const b = "abc"`

一个没有指定类型的常量被使用时，会根据其使用环境而推断出它所需要具备的类型。换句话说，未定义类型的常量会在必要时刻根据上下文来获得相关类型。

```go
var n int
f(n + 5) // 无类型的数字型常量 “5” 它的类型在这里变成了 int
```

常量的值必须是能够在编译时就能够确定的；你可以在其赋值表达式中涉及计算过程，但是所有用于计算的值必须在编译期间就能获得。

* 正确的做法：`const c1 = 2/3`
* 错误的做法：`const c2 = getNumber()` // 引发构建错误: `getNumber() used as value`

**因为在编译期间自定义函数均属于未知，因此无法用于常量的赋值，但内置函数可以使用，如：len\(\)。**

数字型的常量是没有大小和符号的，并且可以使用任何精度而不会导致溢出：

```go
const Ln2 = 0.693147180559945309417232121458\
			176568075500134360255254120680009
const Log2E = 1/Ln2 // this is a precise reciprocal
const Billion = 1e9 // float constant
const hardEight = (1 << 100) >> 97
```

根据上面的例子我们可以看到，反斜杠 `\` 可以在常量表达式中作为多行的连接符使用。

与各种类型的数字型变量相比，你无需担心常量之间的类型转换问题，因为它们都是非常理想的数字。

不过需要注意的是，当常量赋值给一个精度过小的数字型变量时，可能会因为无法正确表达常量所代表的数值而导致溢出，这会在编译期间就引发错误。另外，常量也允许使用并行赋值的形式：

```go
const beef, two, c = "eat", 2, "veg"
const Monday, Tuesday, Wednesday, Thursday, Friday, Saturday = 1, 2, 3, 4, 5, 6
const (
	Monday, Tuesday, Wednesday = 1, 2, 3
	Thursday, Friday, Saturday = 4, 5, 6
)
```

## 枚举

转载自：[https://github.com/unknwon/the-way-to-go\_ZH\_CN/blob/master/eBook/04.3.md](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/04.3.md)

常量还可以用作枚举：

```go
const (
	Unknown = 0
	Female = 1
	Male = 2
)
```

现在，数字 0、1 和 2 分别代表未知性别、女性和男性。这些枚举值可以用于测试某个变量或常量的实际值，比如使用 switch/case 结构.

在这个例子中，`iota` 可以被用作枚举值：

```go
const (
	a = iota
	b = iota
	c = iota
)
```

第一个 `iota` 等于 0，每当 `iota` 在新的一行被使用时，它的值都会自动加 1；所以 `a=0, b=1, c=2` 可以简写为如下形式：

```go
const (
	a = iota
	b
	c
)
```

`iota` 也可以用在表达式中，如：`iota + 50`。在每遇到一个新的常量块或单个常量声明时， `iota` 都会重置为 0（ **简单地讲，每遇到一次 const 关键字，iota 就重置为 0** ）。

当然，常量之所以为常量就是恒定不变的量，因此我们无法在程序运行过程中修改它的值；如果你在代码中试图修改常量的值则会引发编译错误。

引用 `time` 包中的一段代码作为示例：一周中每天的名称。

```go
const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)
```

你也可以使用某个类型作为枚举常量的类型：

```go
type Color int

const (
	RED Color = iota // 0
	ORANGE // 1
	YELLOW // 2
	GREEN // ..
	BLUE
	INDIGO
	VIOLET // 6
)
```

下面这个例子也比较特殊

```go
const (
    Apple, Banana = iota + 1, iota + 2 //0+1，0+2
    Cherimoya, Durian //1+1， 1+2
    Elderberry, Fig //2+1， 2+2
)
```

在看下面这个例子

```go
const (
	enum_a = 'A'
	enum_b
	enum_c = iota
	enum_d
)

const enum_e = iota

func main() {
	fmt.Println(enum_a) //65
	fmt.Println(enum_b) //65
	fmt.Println(enum_c) //2
	fmt.Println(enum_d) //3
	fmt.Println(enum_e) //0
}
```

