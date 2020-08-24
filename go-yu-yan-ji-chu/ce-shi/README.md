# 测试

转载自：[https://github.com/unknwon/the-way-to-go\_ZH\_CN/blob/master/eBook/13.7.md](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/13.7.md) by unknwon

## 简介

首先所有的包都应该有一定的必要文档，然后同样重要的是对包的测试。名为 testing 的包被专门用来进行自动化测试，日志和错误报告。并且还包含一些基准测试函数的功能。

单元测试的目的：

* 确保编写的函数是可以运行的，而且逻辑结果是正确的
* 确保函数性能在预期范围

对一个包做（单元）测试，需要写一些可以频繁（每次更新后）执行的小块测试单元来检查代码的正确性。于是我们必须写一些 Go 源文件来测试代码。测试程序必须属于被测试的包，并且文件名满足这种形式 `*_test.go`，所以测试代码和包中的业务代码是分开的。

`_test` 程序不会被普通的 Go 编译器编译，所以当放应用部署到生产环境时它们不会被部署；只有 go test 会编译所有的程序：普通程序和测试程序。

测试文件中必须导入 "testing" 包，并写一些名字以 `TestZzz` 打头的全局函数，这里的 `Zzz` 是被测试函数的字母描述，如 TestFmtInterface，TestPayEmployees 等。

测试函数必须有这种形式的头部：

```go
func TestAbcde(t *testing.T)
```

T 是传给测试函数的结构类型，用来管理测试状态，支持格式化测试日志，如 t.Log，t.Error，t.ErrorF 等。在函数的结尾把输出跟想要的结果对比，如果不等就打印一个错误。成功的测试则直接返回。

## 常用函数

用下面这些函数来通知测试失败：

1）`func (t *T) Fail()`

```text
	标记测试函数为失败，然后继续执行（剩下的测试）。
```

2）`func (t *T) FailNow()`

```text
	标记测试函数为失败并中止执行；文件中别的测试也被略过，继续执行下一个文件。
```

3）`func (t *T) Log(args ...interface{})`

```text
	args 被用默认的格式格式化并打印到错误日志中。
```

4）`func (t *T) Fatal(args ...interface{})`

```text
	结合 先执行 3），然后执行 2）的效果。
```

运行 go test 来编译测试程序，并执行程序中所有的 TestZZZ 函数。如果所有的测试都通过会打印出 PASS。

go test 可以接收一个或多个函数程序作为参数，并指定一些选项。

结合 --chatty 或 -v 选项，每个执行的测试函数以及测试状态会被打印。

例如：

```go
go test fmt_test.go --chatty
=== RUN fmt.TestFlagParser
--- PASS: fmt.TestFlagParser
=== RUN fmt.TestArrayPrinter
--- PASS: fmt.TestArrayPrinter
...
```

## 实例

```go
//创建一个文件名：operation_test.go 的文件
package main

import "testing"

//一个简单的相加返回结果的函数
func AddUpper(n int) int {
    res := 0
    for i := 1; i <= n; i++{
        res += i
    }
    return res
}

//对AddUpper执行测试
func TestAddUpper(t *testing.T){
    res := Test_addUpper(10)
    if res != 55{
        t.Fatalf("执行错误")
    }
    t.Log("正确")
}

程序运行结果：
=== RUN   TestAddUpper
--- PASS: TestAddUpper (0.00s)
    operation_test.go:10: 正确
PASS
```

## 细节与注意事项

1. 测试用例文件名必须以 \_test.go 结尾。
2. 测试用例函数必须以 `Test` 开头，一般来说就是 Test+被测试的函数名，比如 `TestAddUpper`。**而且Test后接的字母必须是大写。**
3. `TestAddUpper(t *tesing.T)`的形参类型必须是`*testing.T`。
4. 一个测试用例文件中，可以有多个测试用例函数。
5. 运行指令`go test [-v]`。加-v会不管正确和错误都输出日志。
6. 当出现错误时，可以使用 `t.Fatalf`来格式化输出错误信息，并退出程序
7. `t.Logf`方法可以输出相应的日志。
8. `PASS` 表示测试用例运行成功，`FAIL`表示测试用例运行失败。
9. 测试单个文件，一定要带上被测试的原文件 `go test -v add_test.go add.go`。
10. 测试文件并不需要放在`main`包下，只要保证被测试函数所在的包和单元测试在同一包下即可

## 单元测试原理

![](../../.gitbook/assets/image%20%282%29.png)

## 使用表驱动测试

编写测试代码时，一个较好的办法是把测试的输入数据和期望的结果写在一起组成一个数据表：表中的每条记录都是一个含有输入和期望值的完整测试用例，有时还可以结合像测试名字这样的额外信息来让测试输出更多的信息。

可以抽象为下面的代码段：

```go
var tests = []struct{ 	// Test table
	in  string
	out string

}{
	{"in1", "exp1"},
	{"in2", "exp2"},
	{"in3", "exp3"},
...
}

func TestFunction(t *testing.T) {
	for i, tt := range tests {
		s := FuncToBeTested(tt.in)
		if s != tt.out {
			t.Errorf("%d. %q => %q, wanted: %q", i, tt.in, s, tt.out)
		}
	}
}
```

如果大部分函数都可以写成这种形式，那么写一个帮助函数 verify 对实际测试会很有帮助：

```go
func verify(t *testing.T, testnum int, testcase, input, output, expected string) {
	if expected != output {
		t.Errorf("%d. %s with input = %s: output %s != %s", testnum, testcase, input, output, expected)
	}
}
```

TestFunction 则变为：

```go
func TestFunction(t *testing.T) {
	for i, tt := range tests {
		s := FuncToBeTested(tt.in)
		verify(t, i, “FuncToBeTested: “, tt.in, s, tt.out)
	}
}
```

## 基准测试

testing 包中有一些类型和函数可以用来做简单的基准测试；测试代码中必须包含以 `BenchmarkZzz` 打头的函数并接收一个 `*testing.B` 类型的参数，比如：

```go
func BenchmarkReverse(b *testing.B) {
	for i := 0; i<b.N; i++{
		_, _ = AesEncrypt("hello world", "123456781234567812345678")
	}
}
```

命令 `go test –test.bench=.*` 会运行所有的基准测试函数；代码中的函数会被调用 N 次（N是非常大的数，如 N = 1000000），并展示 N 的值和函数执行的平均时间，单位为 ns（纳秒，ns/op）。如果是用 testing.Benchmark 调用这些函数，直接运行程序即可。

```go
BenchmarkAesEncrypt-8   	 1642801	       720 ns/op
PASS
```

## 代码覆盖率

单元测试和性能测试可以看作是测试代码的质量，那么代码覆盖率就是测试代码的编写质量，检查测试代码是否覆盖了足够多的分支语句等。

但是测试还是要以质量为主，而不能过分的追求代码的覆盖率。

![](../../.gitbook/assets/image%20%2814%29.png)



