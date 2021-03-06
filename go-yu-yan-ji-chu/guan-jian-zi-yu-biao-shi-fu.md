# 关键字与标识符

## 关键字

golang中有25个关键字或保留字

![](../.gitbook/assets/image%20%281%29.png)

## 标识符

有效的标识符必须以字符（可以使用任何 UTF-8 编码的字符或 `_`）开头，然后紧跟着 0 个或多个字符或 Unicode 数字，如：X56、group1、\_x23、i、өԑ12。

以下是无效的标识符：

```text
1ab（以数字开头） 
case（Go 语言的关键字） 
a+b（运算符是不允许的）
```

 `_` 本身就是一个特殊的标识符，被称为空白标识符。它可以像其他标识符那样用于变量的声明或赋值（任何类型都可以赋值给它），但任何赋给这个标识符的值都将被抛弃，因此这些值不能在后续的代码中使用，也不可以使用这个标识符作为变量对其它变量进行赋值或运算。

## 命名建议

* 区分大小写
* 使用驼峰命名法
* 在能表达意思的情况下越短越好
* 专有名词建议全部大写

