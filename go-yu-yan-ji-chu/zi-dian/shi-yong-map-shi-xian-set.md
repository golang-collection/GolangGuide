# 使用map实现set

go本身并没有实现set集合，但是set集合在日常开发中也是一种很常见的数据结构，在go中我们可以通过map来实现一个set集合，就像java语言set底层也是通过map实现的一样。

```go
func main() {
	set := map[int]bool{}
	set[1] = true
	fmt.Println(set[1])
	fmt.Println(set[2])
}
//输出
true
false
```

通过拓展上面这段代码，将各种操作封装为函数就可以实现一个基本的set。

```go
type MySet struct {
	data   map[interface{}]bool
	length int
}

func InitSet() *MySet {
	return &MySet{
		data:   make(map[interface{}]bool),
		length: 0,
	}
}

func (set *MySet) AddElementSet(v interface{}) {
	if !set.data[v] {
		set.data[v] = true
		set.length++
	}
}

func (set *MySet) IsExistSet(v interface{}) bool {
	return set.data[v]
}

//return true 删除成功
//return false 删除失败
func (set *MySet) DeleteElementSet(v interface{}) bool {
	if set.data[v] {
		delete(set.data, v)
		return true
	}
	return false
}

func (set *MySet) ShowSet() {
	for k := range set.data{
		fmt.Println(k)
	}
}
```

