# 切片与数组

## 数组与切片

数组和切片都属于集合类型，他们的相同点是都能存储某一类型的元素。不同点是数组的长度是固定的而切片的长度是不固定的，或者说数组的长度是其类型的一部分。

因为指向数组长度为三的指针不能再指向数组长度为4的数组。

![](../../.gitbook/assets/image%20%2815%29.png)

## 切片与数组的关系

对于任何一个切片其底层都有对应的数组，而切片是对底层数组的引用。我们可以想象有一个窗口，你可以通过这个窗口看到一个数组，但是不一定能看到该数组中的所有元素，有时候只能看到连续的一部分元素。这就是切片与数组的关系。

![&#x56FE;&#x7247;&#x6765;&#x6E90; https://time.geekbang.org/column/article/14106](../../.gitbook/assets/image%20%2816%29.png)

> 一个切片的容量可以被看作是透过这个窗口最多可以看到的底层数组中元素的个数。由于slice7是通过在array7上施加切片操作得来的，所以slice7的底层数组就是array7。又因为，在底层数组不变的情况下，切片代表的窗口可以向右扩展，直至其底层数组的末尾。所以，slice7的容量就是其底层数组的长度8, 减去上述切片表达式中的那个起始索引3，即5。
>
> 来源：[https://time.geekbang.org/column/article/14106](https://time.geekbang.org/column/article/14106)

## 多个切片和数组

还需要注意的一种情况是多个切片引用相同数组。此时修改其中一个切片会同步修改其他切片。

不仅多个切片可以引用同一个数组，在切片上还可以再做切片，此时这两个切片都指向相同的底层数组。

```go
func main() {
	array := [10]int{1,2,3,4,5,6,7,8,9,0}
	fmt.Println(array)

	fmt.Println("s1:")
	s1 := array[3:8]
	fmt.Println(s1)
	fmt.Println(reflect.TypeOf(s1))
	fmt.Println(len(s1), cap(s1))

	fmt.Println("s2:")
	s2 := array[5:8]
	fmt.Println(s2)
	fmt.Println(reflect.TypeOf(s2))
	fmt.Println(len(s2), cap(s2))

	fmt.Println("s3:")
	s3 := s1[2:4]
	fmt.Println(s3)
	fmt.Println(reflect.TypeOf(s3))
	fmt.Println(len(s3), cap(s3))
	
	s3[0] = 5
	fmt.Println(array)
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)
}
```

输出结果

```go
[1 2 3 4 5 6 7 8 9 0]
s1:
[4 5 6 7 8]
[]int
5 7
s2:
[6 7 8]
[]int
3 5
s3:
[6 7]
[]int
2 5
```

## 扩容

Go 语言源代码 [runtime/slice.go](https://golang.org/src/runtime/slice.go) 中是这么实现的

```go
// go 1.9.5 src/runtime/slice.go:82
func growslice(et *_type, old slice, cap int) slice {
    // ……
    newcap := old.cap
	doublecap := newcap + newcap
	if cap > doublecap {
		newcap = cap
	} else {
		if old.len < 1024 {
			newcap = doublecap
		} else {
			for newcap < cap {
				newcap += newcap / 4
			}
		}
	}
	// ……
	
	capmem = roundupsize(uintptr(newcap) * ptrSize)
	newcap = int(capmem / ptrSize)
}
```

### 扩容过程

#### 预估容量

如果 oldCap\*2 &lt; cap，那么 newCap=cap

否则 oldLen&lt;1024 那么 cap=oldCap\*2，oldCap&gt;1024=oldCap\*1.25

#### 所占内存

扩容时会将预估容量与变量类型所占字节相乘，然后进行下一步申请内存

#### 申请内存

根据所占内存向提前申请好的内存块进行匹配，匹配到最合适的，然后将拿到的内存大小除以数据类型所占字节，就是最终的cap大小

例子：

![](../../.gitbook/assets/image%20%2843%29.png)

注：这里新申请的string按照其底层实现占16字节

```go
type stringStruct struct {
	str unsafe.Pointer
	len int
}
```

## 参考资料

\[1\] [Go语言核心36讲：07数组和切片](https://time.geekbang.org/column/article/14106)

{% embed url="https://www.bilibili.com/video/BV1CV411d7W8" %}



