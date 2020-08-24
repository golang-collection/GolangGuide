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

//TODO 添加图片

## 扩容

> 一旦一个切片无法容纳更多的元素，Go 语言就会想办法扩容。但它并不会改变原来的切片，而是会生成一个容量更大的切片，然后将把原有的元素和新元素一并拷贝到新切片中。在一般的情况下，你可以简单地认为新切片的容量（以下简称新容量）将会是原切片容量（以下简称原容量）的 2 倍。但是，当原切片的长度（以下简称原长度）大于或等于1024时，Go 语言将会以原容量的1.25倍作为新容量的基准（以下新容量基准）。新容量基准会被调整（不断地与1.25相乘），直到结果不小于原长度与要追加的元素数量之和（以下简称新长度）。最终，新容量往往会比新长度大一些，当然，相等也是可能的。另外，如果我们一次追加的元素过多，以至于使新长度比原容量的 2 倍还要大，那么新容量就会以新长度为基准。注意，与前面那种情况一样，最终的新容量在很多时候都要比新容量基准更大一些。
>
> 来源：[https://time.geekbang.org/column/article/14106](https://time.geekbang.org/column/article/14106)

而且只要新长度不超过切片的原容量，那么使用append函数对其追加元素的时候就不会引起扩容。这只会使紧邻切片窗口右边的（底层数组中的）元素被新的元素替换掉。

## 参考资料

\[1\] [Go语言核心36讲：07数组和切片](https://time.geekbang.org/column/article/14106)





