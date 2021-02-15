# Go 环境变量

在刚开始接触go语言时，一定会对像GOROOT， GOPATH这些内容感到非常疑惑，下面我们就来一起梳理一下。

## GOROOT

GOROOT代表的就是go的安装目录。

## GOPATH

GOPATH可以理解为项目的工作目录，在这个目录下一般会有三个子文件夹。分别

* bin
* pkg
* src

其中src用于存放项目的源码，bin用于存放项目的二进制可执行文件。pkg是golang编译包时，生成的.a文件存放路径。

