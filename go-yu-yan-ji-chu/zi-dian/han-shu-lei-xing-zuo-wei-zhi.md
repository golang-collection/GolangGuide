# 函数类型做为值

在go语言中，函数是一等公民。go语言是支持头等函数（First Class Function）的编程语言，可以把函数赋值给变量，也可以把函数作为其它函数的参数或者返回值。

在go中我们也可以将函数做为值传递给map，例如下面这段代码，map键1，2，3对应的返回原值，平方和三次方。

```go
func main() {
	m := map[int]func(int) int{}
	m[1] = func(i int) int {
		return i
	}
	m[2] = func(i int) int {
		return i * i
	}
	m[3] = func(i int) int {
		return i * i * i
	}
	fmt.Println(m[1](2))
	fmt.Println(m[2](2))
	fmt.Println(m[3](2))
}
```

