# slice底层

切片本质是一个数组片段的描述，包括了数组的指针，这个片段的长度和容量\(不改变内存分配情况下的最大长度\)。

```go
// runtime/slice.go
// slice结构体定义
type slice struct {
    array unsafe.Pointer // 元素指针
    len   int // 长度 
    cap   int // 容量
}
```

切片操作并不复制切片指向的元素，创建一个新的切片会复用原来切片的底层数组，因此切片操作是非常高效的。

![](../../.gitbook/assets/image%20%2840%29.png)

```go
nums := make([]int, 0, 8)
nums = append(nums, 1, 2, 3, 4, 5)
nums2 := nums[2:4]
printLenCap(nums)  // len: 5, cap: 8 [1 2 3 4 5]
printLenCap(nums2) // len: 2, cap: 6 [3 4]

nums2 = append(nums2, 50, 60)
printLenCap(nums)  // len: 5, cap: 8 [1 2 3 4 50]
printLenCap(nums2) // len: 4, cap: 6 [3 4 50 60]
```

* nums2 执行了一个切片操作 `[2, 4)`，此时 nums 和 nums2 指向的是同一个数组。
* nums2 增加 2 个元素 50 和 60 后，将底层数组下标 \[4\] 的值改为了 50，下标\[5\] 的值置为 60。
* 因为 nums 和 nums2 指向的是同一个数组，因此 nums 被修改为 \[1, 2, 3, 4, 50\]。

## 推荐资源

{% embed url="https://ueokande.github.io/go-slice-tricks/" %}

{% embed url="https://geektutu.com/post/hpg-slice.html\#1-2-%E5%88%87%E7%89%87" %}

{% embed url="https://www.cnblogs.com/qcrao-2018/p/10631989.html" %}



