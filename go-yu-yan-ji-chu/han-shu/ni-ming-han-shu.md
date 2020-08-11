# 匿名函数

匿名函数就如字面理解的那样，是一个没有名字的函数。除了没有名字外他与普通函数完全相同。匿名函数可以直接调用，保存到变量，作为参数或者返回值。

## 实例1

```go
func main() {
	func(s string){
		fmt.Println(s)
	}("hello")
}
```

## 实例2

```go
func main() {
	add := func(x, y int) int {
		return x + y
	}
	fmt.Println(add(1,2))
}
```

## 实例3

```go
func main() {
	testArg(func() {
		fmt.Println("world")
	})
}

func testArg(f func()){
	f()
}
```

## 实例4

```go
func main() {
	testReturn := testReturn()
	fmt.Println(testReturn(1,2))
}

func testReturn() func(int, int) int {
	return func(x int, y int) int {
		return x - y
	}
}
```

