# 控制台输入

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

