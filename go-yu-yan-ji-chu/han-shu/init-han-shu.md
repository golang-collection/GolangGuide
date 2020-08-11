# init函数

## init函数简介

在go语言中，除了main函数外，还有一个特殊的函数叫做`init`函数。init函数在每个package是可选的，可有可无，甚至可以有多个\(但是强烈建议一个package中一个init函数\)，init函数在你导入该package时程序会自动调用init函数，所以init函数不用我们手动调用，另外它只会被调用一次，因为当一个package被多次引用时，它只会被导入一次。

## init函数特点

以下内容摘录自：[https://zhuanlan.zhihu.com/p/34211611](https://zhuanlan.zhihu.com/p/34211611)

* init函数先于main函数自动执行，不能被其他函数调用；
* init函数没有输入参数、返回值；
* 每个包可以有多个init函数；
* **包的每个源文件也可以有多个init函数**，这点比较特殊；
* 同一个包的init执行顺序，golang没有明确定义，编程时要注意程序不要依赖这个执行顺序。
* 不同包的init函数按照包导入的依赖关系决定执行顺序。

## 示例1：

{% code title="main.go" %}
```go
package main                                                                                                                     

import (
   "fmt"              
)

var T int64 = a()

func init() {
   fmt.Println("init in main.go ")
}

func a() int64 {
   fmt.Println("calling a()")
   return 2
}
func main() {                  
   fmt.Println("calling main")     
}
```
{% endcode %}

输出：

```go
calling a()
init in main.go
calling main
```

初始化顺序：**变量初始化-&gt;init\(\)-&gt;main\(\)**

## 示例2：

{% code title="pack.go" %}
```go
package pack                                                                                                                     

import (
   "fmt"
   "test_util"
)

var Pack int = 6               

func init() {
   a := test_util.Util        
   fmt.Println("init pack ", a)    
} 
```
{% endcode %}

{% code title="test\_util.go" %}
```go
package test_util                                                                                                                

import "fmt"

var Util int = 5

func init() {
   fmt.Println("init test_util")
}  
```
{% endcode %}

{% code title="main.go" %}
```go
package main                                                                                                                     

import (
   "fmt"
   "pack"
   "test_util"                
)

func main() {                  
   fmt.Println(pack.Pack)     
   fmt.Println(test_util.Util)
}
```
{% endcode %}

输出：

```go
init test_util
init pack  5
6
5
```

**由于pack包的初始化依赖test\_util，因此运行时先初始化test\_util再初始化pack包；**

## 示例3：

{% code title="sandbox.go" %}
```go
package main
import "fmt"
var _ int64 = s()
func init() {
   fmt.Println("init in sandbox.go")
}
func s() int64 {
   fmt.Println("calling s() in sandbox.go")
   return 1
}
func main() {
   fmt.Println("main")
}
```
{% endcode %}

{% code title="a.go" %}
```go
package main
import "fmt"
var _ int64 = a()
func init() {
   fmt.Println("init in a.go")
}
func a() int64 {
   fmt.Println("calling a() in a.go")
   return 2
}
```
{% endcode %}

{% code title="z.go" %}
```go
package main
import "fmt"
var _ int64 = z()
func init() {
   fmt.Println("init in z.go")
}
func z() int64 {
   fmt.Println("calling z() in z.go")
   return 3
}
```
{% endcode %}

输出：

```go
calling a() in a.go
calling s() in sandbox.go
calling z() in z.go
init in a.go
init in sandbox.go
init in z.go
main
```

**同一个包不同源文件的init函数执行顺序，golang spec没做说明，以上述程序输出来看，执行顺序是源文件名称的字典序。**

## 示例4：

```go
package main
import "fmt"
func init() {
   fmt.Println("init")
}
func main() {
   init()
}
```

**init函数不可以被调用，上面代码会提示：undefined: init**

## 示例5：

```go
package main
import "fmt"
func init() {
   fmt.Println("init 1")
}
func init() {
   fmt.Println("init 2")
}
func main() {
   fmt.Println("main")
}
```

输出：

```go
init 1
init 2
main
```

**init函数比较特殊，可以在包里被多次定义。**

## 示例6：

```go
var initArg [20]int

func init() {
   initArg[0] = 10
   for i := 1; i < len(initArg); i++ {
       initArg[i] = initArg[i-1] * 2
   }
}
```

**init函数的主要用途：初始化不能使用初始化表达式初始化的变量**

## 示例7：

```go
import _ "net/http/pprof"
```

**golang对没有使用的导入包会编译报错，但是有时我们只想调用该包的init函数，不使用包导出的变量或者方法，这时就采用上面的导入方案。**



