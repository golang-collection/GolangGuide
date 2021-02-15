# Go项目基本结构

## 项目结构

* Go程序是通过**package**来组织的
* 只有package **名称为main的包**可以包含main函数
* 一个可执行程序**有且仅有**一个main包，_初步练习时可以使用这种结构_

![&#x7EC3;&#x4E60;](https://imgconvert.csdnimg.cn/aHR0cDovL2ltZy5ibG9nLmNzZG4ubmV0LzIwMTcxMTI5MTE1NTAwODQ2?x-oss-process=image/format,png)

## 程序结构

1. package 包名
2. import 要导入的包
3. 多个包可以用import\(包名1 包名2 …\)
4. 如果导入包但是并没有对包进行相关操作就会编译异常
5. import another\_name “包名” 为当前导入的包起别名
6. const 用来定义常量
7. var 用来定义全局变量
8. type 定义普通的变量 语法：type variable\_name int
9. type 定义结构体类型 type variable\_name struct{}
10. type 定义接口类型 type variable\_name interface{}
11. main函数 func main\(\){}

```go
package main

import "fmt"

func main(){
    fmt.Println("hello go")
}
```

## 包的导入

在go语言中可以使用这样的结构进行导包

![](https://img-blog.csdnimg.cn/20190916211320842.png)

但是如果你嫌这样写多个import比较麻烦可以这样写  


![](https://img-blog.csdnimg.cn/20190916211531580.png)

## package的别名

在go语言中我们导包的过程中为了 防止第三包可能会引起冲突，或者为了增强代码的可读性，我们可以使用别名来进行区分。

例如：  


![](https://img-blog.csdnimg.cn/20190916212112612.png)

