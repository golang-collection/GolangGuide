# make

## 官网

{% embed url="https://www.gnu.org/software/make/" %}

以下内容转载自：[http://c.biancheng.net/makefile/](http://c.biancheng.net/makefile/)

## 什么是Makefile

Makefile 可以简单的认为是一个工程文件的编译规则，描述了整个工程的编译和链接等规则。其中包含了那些文件需要编译，那些文件不需要编译，那些文件需要先编译，那些文件需要后编译，那些文件需要重建等等。编译整个工程需要涉及到的，在 Makefile 中都可以进行描述。换句话说，Makefile 可以使得我们的项目工程的编译变得自动化，不需要每次都手动输入一堆源文件和参数。

## Makefile编写规则

Makefile 描述的是文件编译的相关规则，它的规则主要是两个部分组成，分别是依赖的关系和执行的命令，其结构如下所示：

```text
targets : prerequisites
    command
```

或者是

```text
targets : prerequisites; command
    command
```

相关说明如下：

* targets：规则的目标，可以是 Object File（一般称它为中间文件），也可以是可执行文件，还可以是一个标签；
* prerequisites：是我们的依赖文件，要生成 targets 需要的文件或者是目标。可以是多个，也可以是没有；
* command：make 需要执行的命令（任意的 shell 命令）。可以有多条命令，每一条命令占一行。

{% hint style="info" %}
注意：我们的目标和依赖文件之间要使用冒号分隔开，命令的开始一定要使用`Tab`键。
{% endhint %}

通过下面的例子来具体使用一下 Makefile 的规则，Makefile文件中添代码如下：

```text
test:test.c
    gcc -o test test.c
```

上述代码实现的功能就是编译 test.c 文件，通过这个实例可以详细的说明 Makefile 的具体的使用。其中 test 是目标文件，也是我们的最终生成的可执行文件。依赖文件就是 test.c 源文件，重建目标文件需要执行的操作是`gcc -o test test.c`。这就是 Makefile 的基本的语法规则的使用。

{% hint style="info" %}
使用 Makefile 的方式：首先需要编写好 Makefile 文件，然后在 shell 中执行 make 命令，程序就会自动执行，得到最终的目标文件。
{% endhint %}

## Makefile中的内容

主要包含有五个部分，分别是：

### **1\) 显式规则**

显式规则说明了，如何生成一个或多的的目标文件。这是由 Makefile 的书写者明显指出，要生成的文件，文件的依赖文件，生成的命令。

### **2\) 隐晦规则**

由于我们的 make 命名有自动推导的功能，所以隐晦的规则可以让我们比较粗糙地简略地书写 Makefile，这是由 make 命令所支持的。

### **3\) 变量的定义**

在 Makefile 中我们要定义一系列的变量，变量一般都是字符串，这个有点像C语言中的宏，当 Makefile 被执行时，其中的变量都会被扩展到相应的引用位置上。

### **4\) 文件指示**

其包括了三个部分，一个是在一个 Makefile 中引用另一个 Makefile，就像C语言中的 include 一样；另一个是指根据某些情况指定 Makefile 中的有效部分，就像C语言中的预编译 \#if 一样；还有就是定义一个多行的命令。有关这一部分的内容，我会在后续的部分中讲述。

### **5\) 注释**

Makefile 中只有行注释，和 UNIX 的 Shell 脚本一样，其注释是用“\#”字符，这个就像 C/[C++](http://c.biancheng.net/cplus/) 中的“//”一样。如果你要在你的 Makefile 中使用“\#”字符，可以用反斜框进行转义，如：“\\#”。

## 参考文献

\[1\] [Makefile教程](http://c.biancheng.net/makefile/)



