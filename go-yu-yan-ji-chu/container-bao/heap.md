# heap

转载自：[https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter03/03.3.html](https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter03/03.3.html)

这里的堆使用的数据结构是最小二叉树，即根节点比左边子树和右边子树的所有值都小。 go 的堆包只是实现了一个接口，我们看下它的定义：

```go
type Interface interface {
    sort.Interface
    Push(x interface{}) // add x as element Len()
    Pop() interface{}   // remove and return element Len() - 1.
}
```

可以看出，这个堆结构继承自 sort.Interface, 回顾下 sort.Interface，它需要实现三个方法

* Len\(\) int
* Less\(i, j int\) bool
* Swap\(i, j int\)

加上堆接口定义的两个方法

* Push\(x interface{}\)
* Pop\(\) interface{}

就是说你定义了一个堆，就要实现五个方法，直接拿 package doc 中的 example 做例子：

```go
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
    *h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}
```

那么 IntHeap 就实现了这个堆结构，我们就可以使用堆的方法来对它进行操作：

```go
h := &IntHeap{2, 1, 5}
heap.Init(h)
heap.Push(h, 3)
heap.Pop(h)
```

具体说下内部实现，是使用最小堆，索引排序从根节点开始，然后左子树，右子树的顺序方式。索引布局如下：

> ```text
>                   0
>             1            2
>          3    4       5      6
>         7 8  9 10   11   
> ```
>
> 假设 \(heap\[1\]== 小明 \) 它的左子树 \(heap\[3\]== 小黑 \) 和右子树 \(heap\[4\]== 大黄 \) 且 小明 &gt; 小黑 &gt; 大黄 ;

堆内部实现了 down 和 up 函数 : down 函数用于将索引 i 处存储的值 \( 设 i=1, 即小明 \) 与它的左子树 \( 小黑 \) 和右子树 \( 大黄 \) 相比 , 将三者最小的值大黄与小明的位置交换，交换后小明继续与交换后的子树 \(heap\[9\]和 heap\[10\]\) 相比，重复以上步骤，直到小明位置不变。



```text
up 函数用于将索引 i 处的值 ( 设 i=3, 即小黑 ) 与他的父节点 ( 小明 ) 比较，将两者较小的值放到父节点上，本例中即交换小黑和小明的位置，之后小黑继续与其父节点比较，重复以上步骤，直到小黑位置不变。
```

假设 heap\[11\]== 阿花 当从堆中 Pop 一个元素的时候，先把元素和最后一个节点的值 \( 阿花 \) 交换，然后弹出，然后对阿花调用 down，向下保证最小堆。

当往堆中 Push 一个元素的时候，这个元素插入到最后一个节点，本例中为 heap\[12\]，即作为 heap\[5\]的右子树，然后调用 up 函数向上比较。

