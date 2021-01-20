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

### b.N

基准测试框架默认会在持续 1 秒的时间内，反复调用需要测试的函数。测试框架每次调用测试函数时，都会增加 b.N 的值。第一次调用时， b.N 的值为 1 。需要注意，一定要将所有要进行基准测试的代码都放到循环里，并且循环要使用 b.N 的值。 否则，测试的结果是不可靠的。

benchmark 用例的参数 `b *testing.B`，有个属性 `b.N` 表示这个用例需要运行的次数。`b.N` 对于每个用例都是不一样的。

那这个值是如何决定的呢？`b.N` 从 1 开始，如果该用例能够在 1s 内完成，`b.N` 的值便会增加，再次执行。`b.N` 的值大概以 1, 2, 3, 5, 10, 20, 30, 50, 100 这样的序列递增，越到后面，增加得越快。

### 调用cpu核数

在基准测试输出中的-8就代表使用8核，默认就等于cpu核数GOMAXPROCS，可以通过-cpu改变

```go
go test -bench='rse$' -cpu=2,4 .
```

### 提升测试准确度

通过benchtime来使测试时间更长

```go
go test -bench='rse$' -benchtime=5s .
```

除了制定制定多长时间还可以制定执行多少次

```go
go test -bench='rse$' -benchtime=50x .
```

`count` 参数可以用来设置 benchmark 的轮数

```go
go test -bench='rse$' -benchtime=5s -count=3 .
```

### 内存分配情况

转载自：[https://geektutu.com/post/hpg-benchmark.html](https://geektutu.com/post/hpg-benchmark.html)

`-benchmem` 参数可以度量内存分配的次数。内存分配次数也性能也是息息相关的，例如不合理的切片容量，将导致内存重新分配，带来不必要的开销。

在下面的例子中，`generateWithCap` 和 `generate` 的作用是一致的，生成一组长度为 n 的随机序列。唯一的不同在于，`generateWithCap` 创建切片时，将切片的容量\(capacity\)设置为 n，这样切片就会一次性申请 n 个整数所需的内存。

```go
package main

import (
	"math/rand"
	"testing"
	"time"
)

func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func generate(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func BenchmarkGenerateWithCap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generateWithCap(1000000)
	}
}

func BenchmarkGenerate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generate(1000000)
	}
}
```

程序输出结果

```go
goos: darwin
goarch: amd64
pkg: example
BenchmarkGenerateWithCap-8            54          21915768 ns/op
BenchmarkGenerate-8                   43          27031462 ns/op
PASS
ok      example   6.790s
```

可以使用 `-benchmem` 参数看到内存分配的情况：

```go
goos: darwin
goarch: amd64
pkg: example
BenchmarkGenerateWithCap-8            50          22020888 ns/op         8003665 B/op          1 allocs/op
BenchmarkGenerate-8                   42          26996118 ns/op        45188388 B/op         40 allocs/op
PASS
ok      example   4.690s

```

`Generate` 分配的内存是 `GenerateWithCap` 的 6 倍，设置了切片容量，内存只分配一次，而不设置切片容量，内存分配了 40 次。

### 测试不同输入

不同的函数复杂度不同，O\(1\)，O\(n\)，O\(n^2\) 等，利用 benchmark 验证复杂度一个简单的方式，是构造不同的输入。

```go
import (
	"math/rand"
	"testing"
	"time"
)

func generate(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}
func benchmarkGenerate(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		generate(i)
	}
}

func BenchmarkGenerate1000(b *testing.B)    { benchmarkGenerate(1000, b) }
func BenchmarkGenerate10000(b *testing.B)   { benchmarkGenerate(10000, b) }
func BenchmarkGenerate100000(b *testing.B)  { benchmarkGenerate(100000, b) }
func BenchmarkGenerate1000000(b *testing.B) { benchmarkGenerate(1000000, b) }
```

这里，我们实现一个辅助函数 `benchmarkGenerate` 允许传入参数 i，并构造了 4 个不同输入的 benchmark 用例。运行结果如下：

```go
goos: darwin
goarch: amd64
pkg: example
BenchmarkGenerate1000-8            38001             31278 ns/op
BenchmarkGenerate10000-8            4448            250563 ns/op
BenchmarkGenerate100000-8            480           2420767 ns/op
BenchmarkGenerate1000000-8            46          25209769 ns/op
PASS
```

通过测试结果可以发现，输入变为原来的 10 倍，函数每次调用的时长也差不多是原来的 10 倍，这说明复杂度是线性的。

### 注意事项

#### ResetTimer

如果在 benchmark 开始前，需要一些准备工作，如果准备工作比较耗时，则需要将这部分代码的耗时忽略掉。比如下面的例子：

```go
func BenchmarkFib(b *testing.B) {
	time.Sleep(time.Second * 3) // 模拟耗时准备任务
	for n := 0; n < b.N; n++ {
		fib(30) // run fib(30) b.N times
	}
}
```

运行结果是：

```go
$ go test -bench='Fib$' -benchtime=50x .
goos: darwin
goarch: amd64
pkg: example
BenchmarkFib-8                50          65912552 ns/op
PASS
ok      example 6.319s
```

50次调用，每次调用约 0.66s，是之前的 0.06s 的 11 倍。究其原因，受到了耗时准备任务的干扰。我们需要用 `ResetTimer` 屏蔽掉：

```go
func BenchmarkFib(b *testing.B) {
	time.Sleep(time.Second * 3) // 模拟耗时准备任务
	b.ResetTimer() // 重置定时器
	for n := 0; n < b.N; n++ {
		fib(30) // run fib(30) b.N times
	}
}
```

运行结果恢复正常，每次调用约 0.06s。

```go
$ go test -bench='Fib$' -benchtime=50x .
goos: darwin
goarch: amd64
pkg: example
BenchmarkFib-8                50           6187485 ns/op
PASS
ok      example 6.330s
```

#### StopTimer & StartTimer <a id="3-2-StopTimer-amp-StartTimer"></a>

还有一种情况，每次函数调用前后需要一些准备工作和清理工作，我们可以使用 `StopTimer` 暂停计时以及使用 `StartTimer` 开始计时。

例如，如果测试一个冒泡函数的性能，每次调用冒泡函数前，需要随机生成一个数字序列，这是非常耗时的操作，这种场景下，就需要使用 `StopTimer` 和 `StartTimer` 避免将这部分时间计算在内。例如：

```go
package main

import (
	"math/rand"
	"testing"
	"time"
)

func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func bubbleSort(nums []int) {
	for i := 0; i < len(nums); i++ {
		for j := 1; j < len(nums)-i; j++ {
			if nums[j] < nums[j-1] {
				nums[j], nums[j-1] = nums[j-1], nums[j]
			}
		}
	}
}

func BenchmarkBubbleSort(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		nums := generateWithCap(10000)
		b.StartTimer()
		bubbleSort(nums)
	}
}
```

执行该用例，每次排序耗时约 0.1s。

```go
$ go test -bench='Sort$' .
goos: darwin
goarch: amd64
pkg: example
BenchmarkBubbleSort-8                  9         113280509 ns/op
PASS
ok      example 1.146s
```

## 代码覆盖率

单元测试和性能测试可以看作是测试代码的质量，那么代码覆盖率就是测试代码的编写质量，检查测试代码是否覆盖了足够多的分支语句等。

但是测试还是要以质量为主，而不能过分的追求代码的覆盖率。

![](../../.gitbook/assets/image%20%2814%29.png)

## 推荐阅读

{% embed url="https://mp.weixin.qq.com/s/LDYJMZ72k9PCSiBwJAMKlQ" %}



