# rand包

在go语言中生成随机数也非常方便，下面我们就通过代码来演示在go语言中如何生成随机数。

```go
func main() {
	//返回0<=n<=参数的值
	fmt.Println(rand.Intn(10000))
	//返回一个64位浮点数f，0.0 <= f <= 1.0。
	fmt.Println(rand.Float64())
	//生成随机数f，范围在5.0 <= f <= 10.0
	fmt.Println((rand.Float64() * 5) + 5)
	//默认情况下，给定的种子是确定的，每次都会产生相同的随机数数字序列。
	// 要产生变化的序列，需要给定一个变化的种子。
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	fmt.Print(r1.Intn(100))
}
```

程序输出

```go
8081
0.9405090880450124
8.322800266092452
12
```

math/rand包生成随机数的算法不够安全，我们可以使用**crypto**的rand包，所谓的安全就是通过crypto/rand生成的随机内容更加的不可预测。

