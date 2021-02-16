# 控制台输入

go语言中控制台输入可以通过fmt中的sacn相关的函数来实现，`Scanln` 扫描来自标准输入的文本，将空格分隔的值依次存放到后续的参数内，直到碰到换行。`Scanf` 与其类似，除了 `Scanf` 的第一个参数用作格式字符串，用来决定如何读取。`Sscan` 和以 `Sscan` 开头的函数则是从字符串读取，除此之外，与 `Scanf` 相同。

```go
func main() {
	var str string
	var a int
	var b float64
	fmt.Println("请输入内容")
	fmt.Scanf("%s", &str)
	fmt.Scanf("%d", &a)
	fmt.Scanf("%f", &b)
	fmt.Printf("str: %s, a: %d, b: %f", str, a, b)
}
```

输入与输出

```go
请输入内容
str
1
2.5
str: str, a: 1, b: 2.500000
```

## 循环输入

```go
func main() {
	var a, b int
	for {
		n, _ := fmt.Scanf("%d %d", &a, &b)
		if n != 2 {
			break
		}
		fmt.Println(a + b)
	}
}
```

