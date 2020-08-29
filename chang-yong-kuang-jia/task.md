# Task

## 官网

{% embed url="https://taskfile.dev/\#/" %}

[Task](https://github.com/go-task/task)是一个任务运行器/构建工具，它的目标是比GNU Make等工具更简单更容易使用。

## 安装

* Mac：

```bash
brew install go-task/tap/go-task
```

* Linux

```bash
sudo snap install task --classic
```

* Windows

```bash
scoop bucket add extras
scoop install task
```

## 使用

安装成功后，在项目的根目录创建`Taskfile.yml`文件。

```yaml
version: '3'

tasks:
  build:
    cmds:
      - go build -v -i main.go
```

之后在终端输入以下命令即可成功构建项目

```bash
task build
```

注意`build`为yaml文件中tasks对应的名称，如果yaml文件中有多个task则使用以下命令即可

```bash
task task1 task2
```

