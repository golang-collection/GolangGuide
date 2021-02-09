# if赋值语句

```go
type A struct {
}

func newObject() (interface{}, error) {
	return A{}, errors.New("error")
}

func main() {
	var a interface{}
	if a, err := newObject(); err != nil {
		fmt.Println(a)
		fmt.Println(err)
	}
	fmt.Println(a)
}
```

当前程序输出

```go
{}
error
<nil>
```

要注意此种错误，所以一般的变量尽量在函数开头统一通过var定义。

