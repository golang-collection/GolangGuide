# Mac OS

本节主要为大家讲解如何在Mac OS上安装Go语言开发包，大家可以在Go语言官网下载对应版本的的安装包（[https://golang.google.cn/dl/](https://golang.google.cn/dl/)），如下图所示。

![](../.gitbook/assets/image%20%283%29.png)

##  安装Go语言开发包

 Mac OS 的Go语言开发包是 .pkg 格式的，双击我们下载的安装包即可开始安装。

![](../.gitbook/assets/image%20%284%29.png)

 Mac OS 下是傻瓜式安装，一路点击“继续”即可，不再赘述。

![](../.gitbook/assets/image%20%285%29.png)

 安装包会默认安装在 /usr/local 目录下，如下所示。

![](../.gitbook/assets/image%20%286%29.png)

 安装完成之后，在终端运行 `go version`，如果显示类似下面的信息，表明安装成功。

```bash
go version
go 1.13.4 darwin/amd64
```



