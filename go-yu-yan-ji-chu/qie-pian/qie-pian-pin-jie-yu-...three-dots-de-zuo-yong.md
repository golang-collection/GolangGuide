# 切片拼接与...\(three dots\)的作用

先看一段代码

```go
func main() {
    nums1 := []int{1,2,3,4}
    nums2 := []int{2,6,7,8}
    nums1 = append(nums1, nums2...)
    fmt.Println(nums1)
}

输出结果
[1 2 3 4 2 6 7 8]
```

要了解为什么可以这样传递，我们就需要深入了解一下golang中的`...`作用是什么

## golang中`...`的作用

### 用法一

我们都知道在golang中可以使用`...`来作为函数的可变参数，使用这种写法,其中`arg`作为可变参数项可以接受多个`int`型数据。

```go
func test(arg ...int){
    fmt.Printf("arg类型: %T\n", arg)
    fmt.Printf("arg地址: %p\n", &arg)
    fmt.Println(arg)
    arg[1] = 12
}

func main() {
    nums1 := []int{1,2,3,4}
    fmt.Printf("nums1地址: %p\n", &nums1)
    test(nums1...)
    fmt.Println(nums1)
}

输出结果:
nums1地址: 0xc00000c060
arg类型: []int
arg地址: 0xc00000c080
[1 2 3 4]
[1 12 3 4]
```

但是我们还有几个细节需要注意

* 将切片传递给可变参数，并不会新创建一个切片，相当于引用传递，在`test`函数中修改切片的值，会对原来的切片有影响
* 从程序的效果来看确实是引用传递，但是`nums1`与`arg`的地址却不一样，这一点可以探究一下
* 从`arg`的类型可以看出，`arg`确实也为切片类型
* **如果我们直接将切片传入可变参数是会报错的，需要在切片后加`...`，因为加入`...`之后会将切片unpacking\(我的理解是将切片拆箱，将其中的值传递给可变参数，传递后再将其包装成切片\)**

### 用法二

```go
func main(){
    arr := [...]int{1,2,3,4}
}
```

在定义数组时也可以使用`...`, 使用这种方法可以不用显示的声明数组长度。
