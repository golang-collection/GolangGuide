# flag包

转载自：[http://c.biancheng.net/view/5573.html](http://c.biancheng.net/view/5573.html)

在编写命令行程序（工具、server）时，我们有时需要对命令参数进行解析，各种编程语言一般都会提供解析命令行参数的方法或库，以方便程序员使用。在Go语言中的 flag 包中，提供了命令行参数解析的功能。

下面我们就来看一下 flag 包可以做什么，它具有什么样的能力。

这里介绍几个概念：

* 命令行参数（或参数）：是指运行程序时提供的参数；
* 已定义命令行参数：是指程序中通过 flag.Type 这种形式定义了的参数；
* 非 flag（non-flag）命令行参数（或保留的命令行参数）：可以简单理解为 flag 包不能解析的参数。

## flag 包概述

Go语言内置的 flag 包实现了命令行参数的解析，flag 包使得开发命令行工具更为简单。若要使用 flag 包，首先需要使用 import 关键字导入 flag 包，如下所示：

```go
import "flag"
```

## flag 参数类型

flag 包支持的命令行参数类型有 bool、int、int64、uint、uint64、float、float64、string、duration，如下表所示：

| flag 参数 | 有效值 |
| :--- | :--- |
| 字符串 flag | 合法字符串 |
| 整数 flag | 1234、0664、0x1234 等类型，也可以是负数 |
| 浮点数 flag | 合法浮点数 |
| bool 类型 flag | 1、0、t、f、T、F、true、false、TRUE、FALSE、True、False |
| 时间段 flag | 任何合法的时间段字符串，如“300ms”、“-1.5h”、“2h45m”， 合法的单位有“ns”、“us”、“µs”、“ms”、“s”、“m”、“h” |

## flag 包基本使用

有以下两种常用的定义命令行 flag 参数的方法：

### **1\) flag.Type\(\)**

基本格式如下：

flag.Type\(flag 名, 默认值, 帮助信息\) \*TypeType 可以是 Int、String、Bool 等，返回值为一个相应类型的指针，例如我们要定义姓名、年龄、婚否三个命令行参数，我们可以按如下方式定义：

```go
name := flag.String("name", "张三", "姓名")
age := flag.Int("age", 18, "年龄")
married := flag.Bool("married", false, "婚否")
delay := flag.Duration("d", 0, "时间间隔")
```

需要注意的是，此时 name、age、married、delay 均为对应类型的指针。

### **2\) flag.TypeVar\(\)**

基本格式如下：

flag.TypeVar\(Type 指针, flag 名, 默认值, 帮助信息\)TypeVar 可以是 IntVar、StringVar、BoolVar 等，其功能为将 flag 绑定到一个变量上，例如我们要定义姓名、年龄、婚否三个命令行参数，我们可以按如下方式定义：  


```go
var name string
var age int
var married bool
var delay time.Duration
flag.StringVar(&name, "name", "张三", "姓名")
flag.IntVar(&age, "age", 18, "年龄")
flag.BoolVar(&married, "married", false, "婚否")
flag.DurationVar(&delay, "d", 0, "时间间隔")
```

### **flag.Parse\(\)**

通过以上两种方法定义好命令行 flag 参数后，需要通过调用 flag.Parse\(\) 来对命令行参数进行解析。

支持的命令行参数格式有以下几种：

* -flag：只支持 bool 类型；
* -flag=x；
* -flag x：只支持非 bool 类型。

其中，布尔类型的参数必须使用等号的方式指定。

flag 包的其他函数：

flag.Args\(\)  //返回命令行参数后的其他参数，以 \[\]string 类型  
flag.NArg\(\)  //返回命令行参数后的其他参数个数  
flag.NFlag\(\) //返回使用的命令行参 数个数结合上面的介绍知识，我们来看一个实例，代码如下：

```go
package main

import (    
  "flag"    
  "fmt"
)

var Input_pstrName = flag.String("name", "gerry", "input ur name")
var Input_piAge = flag.Int("age", 20, "input ur age")
var Input_flagvar int

func Init() {    
  flag.IntVar(&Input_flagvar, "flagname", 1234, "help message for flagname")
}

func main() {    
  Init()    
  flag.Parse()    
  // After parsing, the arguments after the flag are available as the slice flag.Args() or individually as flag.Arg(i). The arguments are indexed from 0 through flag.NArg()-1    
  // Args returns the non-flag command-line arguments    
  // NArg is the number of arguments remaining after flags have been processed    
  fmt.Printf("args=%s, num=%d\n", flag.Args(), flag.NArg())    
  for i := 0; i != flag.NArg(); i++ {        
    fmt.Printf("arg[%d]=%s\n", i, flag.Arg(i))    
  }    
  fmt.Println("name=", *Input_pstrName)    
  fmt.Println("age=", *Input_piAge)    
  fmt.Println("flagname=", Input_flagvar)
}
```

运行结果如下：

```go
go run main.go -name "aaa" -age=123 -flagname=999
args=[], num=0
name= aaa
age= 123
flagname= 999
```

### **自定义 Value**

另外，我们还可以创建自定义 flag，只要实现 flag.Value 接口即可（要求 receiver 是指针类型），这时候可以通过如下方式定义该 flag：

flag.Var\(&flagVal, "name", "help message for flagname"\)

【示例】解析喜欢的编程语言，并直接解析到 slice 中，我们可以定义如下 sliceValue 类型，然后实现 Value 接口：

```go
package main
import (
    "flag"
    "fmt"
    "strings"
)
//定义一个类型，用于增加该类型方法
type sliceValue []string
//new一个存放命令行参数值的slice
func newSliceValue(vals []string, p *[]string) *sliceValue {
    *p = vals
    return (*sliceValue)(p)
}
/*
Value接口：
type Value interface {
    String() string
    Set(string) error
}
实现flag包中的Value接口，将命令行接收到的值用,分隔存到slice里
*/
func (s *sliceValue) Set(val string) error {
    *s = sliceValue(strings.Split(val, ","))
    return nil
}
//flag为slice的默认值default is me,和return返回值没有关系
func (s *sliceValue) String() string {
    *s = sliceValue(strings.Split("default is me", ","))
    return "It's none of my business"
}
/*
可执行文件名 -slice="java,go"  最后将输出[java,go]
可执行文件名 最后将输出[default is me]
*/
func main(){
    var languages []string
    flag.Var(newSliceValue([]string{}, &languages), "slice", "I like programming `languages`")
    flag.Parse()
    //打印结果slice接收到的值
    fmt.Println(languages)
}
```

通过`-slice go,php` 这样的形式传递参数，languages 得到的就是 \[go, php\]，如果不加`-slice` 参数则打印默认值`[default is me]`，如下所示：

```go
go run main.go -slice go,php,java
[go php java]
```

flag 中对 Duration 这种非基本类型的支持，使用的就是类似这样的方式，即同样实现了 Value 接口。

