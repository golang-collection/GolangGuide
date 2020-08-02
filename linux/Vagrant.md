# 添加box

下载box之后，通过命令

```
vagrant box add name/version path
```

添加命令之后通过命令查看是否添加成功

```
vagrant box list
```

# 加载Vagrant

## 方式一

通过vagrant直接创建虚拟机

```
vagrant init centos/7
```

就会在当前目录下创建一个centos7的文件夹，并会在其中创建一个Vagrantfile的文件，用于创建虚拟机。

## 方式二

直接创建Vagrantfile然后在目录下启动即可

# vagrant相关命令

## 启动

```
vagrant up
```

就可以启动虚拟机

## 挂起

```
vagrant halt
```



## 删除

```
vagrant destroy
```







