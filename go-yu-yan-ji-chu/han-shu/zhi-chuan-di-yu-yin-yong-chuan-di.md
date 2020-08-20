# 值传递与引用传递

## 1. 值传递

**基本数据类型和数组作为参数会进行值传递**

接下来看一个最简单的例子

```go
import "fmt"

func test(num int) {
    num = num +1
    fmt.Println("test : ", num)
}

func main() {
    num := 10
    test(num)
    fmt.Println("main : ", num)
}
```

如果大家有其他语言的开发经验，应该很容易看出来程序的输出结果是

```go
test : 11
main : 10
```

我们再来看看引用传递是什么结果

## 2. 引用传递

还是与上面的例子相同，只不过函数参数改为指针类型

```go
import "fmt"

func test(ptr *int) {
    *ptr = *ptr +1
    fmt.Println("test : ", *ptr)
}

func main() {
    num := 10
    ptr := &num
    test(ptr)
    fmt.Println("main : ", num)
}
```

我们再来看一下输出结果

```go
test : 11
main : 11
```

这两种不同的输出结果，底层到底是如何实现的呢，

## 3. 程序运行时的内存分析

对于程序而言，在运行是操作系统会为其分配一块内存，以满足程序的运行需要，程序的进程会将这块内存分为三个部分，分别是：1. 栈区 2. 堆区 3. 代码区。**这是人为的逻辑上的划分。**

* 栈区：一般来说存储基本数据类型
* 堆区：一般来说存储引用数据类型
* 代码区： 存储代码本身

### 3.1 值传递内存分析

下面我们就来对照一下代码与逻辑上的分区来看一下到底是什么输出结果  

![main &#x51FD;&#x6570;](https://img-blog.csdnimg.cn/20200212122607648.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3NfODQyNDk5NDY3,size_16,color_FFFFFF,t_70)

首先，代码从main函数开始执行，会在栈中开辟一块区域用来存储main函数的相关变量，注意这里只是人为的逻辑分区。此时就在main函数的栈区中开辟出一块空间用来存储变量`num`。 

![test&#x51FD;&#x6570;](https://img-blog.csdnimg.cn/20200212123342278.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3NfODQyNDk5NDY3,size_16,color_FFFFFF,t_70)

此时程序运行到`test`函数处，此时在栈中开辟一块空间作为test函数的栈区，其中也会开辟一块空间存储一个`num`变量，只不过此时的`test`函数栈区中的`num`与`main`函数栈区中的`num`是两个变量，其中`test`函数栈去中的`num`的值为`main`函数栈区中`num`值的拷贝。  

![test&#x51FD;&#x6570;2](https://img-blog.csdnimg.cn/20200212123956293.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3NfODQyNDk5NDY3,size_16,color_FFFFFF,t_70)

此时程序运行到`test`函数中，将`test`函数栈区中`num`做`+1`操作，操作结束之后执行输出语句，输出的是`test`函数栈区中的`num`，故输出`test : 11`，之后程序返回到`main`函数中继续执行。  

![main&#x51FD;&#x6570;2](https://img-blog.csdnimg.cn/20200212124328449.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3NfODQyNDk5NDY3,size_16,color_FFFFFF,t_70)

注意此时`test`函数已经执行完毕，所以将`test`函数的栈去已经自动从栈区中删除了，此时再执行输出语句，输出的是`main`函数栈区中保存的`num`，故输出`num : 10`。

### 3.2 引用传递内存分析

 

![main&#x51FD;&#x6570;3](https://img-blog.csdnimg.cn/20200212131955928.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3NfODQyNDk5NDY3,size_16,color_FFFFFF,t_70)

与前面类似，也首先在栈去中开辟一块空间存储`main`函数的相关变量。  

![main&#x5806;&#x533A;](https://img-blog.csdnimg.cn/20200212133104159.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3NfODQyNDk5NDY3,size_16,color_FFFFFF,t_70)

程序向下运行，此时会在堆区中分配一块空间，用来存放`main`函数的一些引用变量，其中的`ptr`为`int`类型的指针，其值为`main`函数中`num`的地址，也就是`ptr`为指向`num`的指针。  

![test&#x51FD;&#x6570;&#x5806;](https://img-blog.csdnimg.cn/20200212133440211.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3NfODQyNDk5NDY3,size_16,color_FFFFFF,t_70)

函数执行到这里的时候，会在栈区和堆区也分别创建`test`函数的空间，其中在`test`函数堆中存储了一个`ptr`指针，其值为`main`函数堆中`ptr`值的拷贝，所以这两个`ptr`指针保存的都是`main`函数中`num`的地址，也就是说`test`函数堆区中的`ptr`指针也指向`main`函数栈区的`num`变量。  

![test&#x51FD;&#x6570;&#x5806;](https://img-blog.csdnimg.cn/20200212133836946.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3NfODQyNDk5NDY3,size_16,color_FFFFFF,t_70)

此时操作`test`堆区中的`ptr`指针，也就像相当于操作`main`函数中的`num`变量，将其执行`+1`操作，程序继续执行故输出为`test : 11`。  

![main&#x5806;](https://img-blog.csdnimg.cn/2020021213410358.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3NfODQyNDk5NDY3,size_16,color_FFFFFF,t_70)

程序继续执行到`main`函数中的输出语句，此时已将`test`栈区所占空间自动释放，`test`堆区由`GC`机制决定何时释放。`main`函数输出的值就是`main : 11`。

## 4. 总结

以上就是值传递与引用传递的分析。注意其中的栈，堆等均为人为的逻辑分区，每个程序在运行过程中也未必会严格按照此进行内存的分配，这里是为了解释方便，后期在做`GC`等分析的时候也会有更详细的说明。

