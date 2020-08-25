# os包

转载自：[http://c.biancheng.net/view/5572.html](http://c.biancheng.net/view/5572.html)

Go语言的 os 包中提供了操作系统函数的接口，是一个比较重要的包。顾名思义，os 包的作用主要是在服务器上进行系统的基本操作，如文件操作、目录操作、执行命令、信号与中断、进程、系统状态等等。

##  os 包中的常用函数

###  1\) Hostname

 函数定义:

```go
func Hostname() (name string, err error)
```

 Hostname 函数会返回内核提供的主机名。

###  2\) Environ

 函数定义:

```go
func Environ() []string
```

 Environ 函数会返回所有的环境变量，返回值格式为“key=value”的字符串的切片拷贝。

###  3\) Getenv

 函数定义:

```go
func Getenv(key string) string
```

 Getenv 函数会检索并返回名为 key 的环境变量的值。如果不存在该环境变量则会返回空字符串。

###  4\) Setenv

 函数定义:

```go
func Setenv(key, value string) error
```

 Setenv 函数可以设置名为 key 的环境变量，如果出错会返回该错误。

###  5\) Exit

 函数定义:

```go
func Exit(code int)
```

 Exit 函数可以让当前程序以给出的状态码 code 退出。

一般来说，状态码 0 表示成功，非 0 表示出错。程序会立刻终止，并且 defer 的函数不会被执行。

###  6\) Getuid

 函数定义:

```go
func Getuid() int
```

 Getuid 函数可以返回调用者的用户 ID。

###  7\) Getgid

 函数定义:

```go
func Getgid() int
```

 Getgid 函数可以返回调用者的组 ID。

###  8\) Getpid

 函数定义:

```go
func Getpid() int
```

 Getpid 函数可以返回调用者所在进程的进程 ID。

###  9\) Getwd

 函数定义:

```go
func Getwd() (dir string, err error)
```

 Getwd 函数可以返回一个对应当前工作目录的根路径。如果当前目录可以经过多条路径抵达（因为硬链接），Getwd 会返回其中一个。

###  10\) Mkdir

 函数定义:

```go
func Mkdir(name string, perm FileMode) error
```

 Mkdir 函数可以使用指定的权限和名称创建一个目录。如果出错，会返回 \*PathError 底层类型的错误。

###  11\) MkdirAll

 函数定义:

```go
func MkdirAll(path string, perm FileMode) error
```

 MkdirAll 函数可以使用指定的权限和名称创建一个目录，包括任何必要的上级目录，并返回 nil，否则返回错误。权限位 perm 会应用在每一个被该函数创建的目录上。如果 path 指定了一个已经存在的目录，MkdirAll 不做任何操作并返回 nil。

###  12\) Remove

 函数定义:

```go
func Remove(name string) error
```

 Remove 函数会删除 name 指定的文件或目录。如果出错，会返回 \*PathError 底层类型的错误。

 RemoveAll 函数跟 Remove 用法一样，区别是会递归的删除所有子目录和文件。

 在 os 包下，有 exec，signal，user 三个子包，下面来分别介绍一下。

##  os/exec 执行外部命令

exec 包可以执行外部命令，它包装了 os.StartProcess 函数以便更容易的修正输入和输出，使用管道连接 I/O，以及作其它的一些调整。

```go
func LookPath(file string) (string, error)
```

在环境变量 PATH 指定的目录中搜索可执行文件，如果 file 中有斜杠，则只在当前目录搜索。返回完整路径或者相对于当前目录的一个相对路径。

 示例代码如下：

```go
package main

import (
    "fmt"
    "os/exec"
)

func main() {
    f, err := exec.LookPath("main")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(f)
}
```

 运行结果如下：

 main.exe

## os/user 获取当前用户信息

可以通过 os/user 包中的 Current\(\) 函数来获取当前用户信息，该函数会返回一个 User 结构体，结构体中的 Username、Uid、HomeDir、Gid 分别表示当前用户的名称、用户 id、用户主目录和用户所属组 id，函数原型如下：

```go
func Current() (*User, error)
```

示例代码如下：

```go
package main

import (
    "log"
    "os/user"
)

func main() {
    u, _ := user.Current()
    log.Println("用户名：", u.Username)
    log.Println("用户id", u.Uid)
    log.Println("用户主目录：", u.HomeDir)
    log.Println("主组id：", u.Gid)

    // 用户所在的所有的组的id
    s, _ := u.GroupIds()
    log.Println("用户所在的所有组：", s)
}
```

 运行结果如下：

 2019/12/13 15:12:14 用户名： LENOVO-PC\Administrator  
 2019/12/13 15:12:14 用户id S-1-5-21-711400000-2334436127-1750000211-000  
 2019/12/13 15:12:14 用户主目录： C:\Users\Administrator  
 2019/12/13 15:12:14 主组id： S-1-5-22-766000000-2300000100-1050000262-000  
 2019/12/13 15:12:14 用户所在的所有组： \[S-1-5-32-544 S-1-5-22-000 S-1-5-21-777400999-2344436111-1750000262-003\]

## os/signal 信号处理

 一个运行良好的程序在退出（正常退出或者强制退出，如 Ctrl+C，kill 等）时是可以执行一段清理代码的，将收尾工作做完后再真正退出。一般采用系统 Signal 来通知系统退出，如 kill pid，在程序中针对一些系统信号设置了处理函数，当收到信号后，会执行相关清理程序或通知各个子进程做自清理。

 Go语言中对信号的处理主要使用 os/signal 包中的两个方法，一个是 Notify 方法用来监听收到的信号，一个是 stop 方法用来取消监听。

```go
func Notify(c chan<- os.Signal, sig ...os.Signal)
```

 其中，第一个参数表示接收信号的 channel，第二个及后面的参数表示设置要监听的信号，如果不设置表示监听所有的信号。

 【示例 1】使用 Notify 方法来监听收到的信号：

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
)

func main() {
    c := make(chan os.Signal, 0)
    signal.Notify(c)

    // Block until a signal is received.
    s := <-c
    fmt.Println("Got signal:", s)
}
```

 运行该程序，然后在 CMD 窗口中通过 Ctrl+C 来结束该程序，便会得到输出结果：

```go
Got signal: interrupt
```

 【示例 2】使用 stop 方法来取消监听：

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
)

func main() {
    c := make(chan os.Signal, 0)
    signal.Notify(c)

    signal.Stop(c) //不允许继续往c中存入内容
    s := <-c       //c无内容，此处阻塞，所以不会执行下面的语句，也就没有输出
    fmt.Println("Got signal:", s)
}
```

 因为使用 Stop 方法取消了 Notify 方法的监听，所以运行程序没有输出结果。

