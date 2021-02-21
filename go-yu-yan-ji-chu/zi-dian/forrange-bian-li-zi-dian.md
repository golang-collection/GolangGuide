# for-range遍历字典

转载自：[https://github.com/unknwon/the-way-to-go\_ZH\_CN/blob/master/eBook/08.3.md](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/08.3.md) by unknown

可以使用 for 循环遍历 map：

```go
for key, value := range map1 {
	...
}
```

第一个返回值 key 是 map 中的 key 值，第二个返回值则是该 key 对应的 value 值；这两个都是仅 for 循环内部可见的局部变量。其中第一个返回值key值是一个可选元素。如果你只关心值，可以这么使用：

```go
for _, value := range map1 {
	...
}
```

如果只想获取 key，你可以这么使用：

```go
for key := range map1 {
	fmt.Printf("key is: %d\n", key)
}
```

示例

```go
package main

import "fmt"

func main() {
	map1 := make(map[int]float32)
	map1[1] = 1.0
	map1[2] = 2.0
	map1[3] = 3.0
	map1[4] = 4.0
	for key, value := range map1 {
		fmt.Printf("key is: %d - value is: %f\n", key, value)
	}
}
```

输出结果：

```text
key is: 3 - value is: 3.000000
key is: 1 - value is: 1.000000
key is: 4 - value is: 4.000000
key is: 2 - value is: 2.000000
```

{% hint style="info" %}
**注意 map 不是按照 key 的顺序排列的，也不是按照 value 的序排列的。**
{% endhint %}

