# 方法

转载自：[https://www.flysnow.org/2017/03/31/go-in-action-go-method.html](https://www.flysnow.org/2017/03/31/go-in-action-go-method.html) by 飞雪无情

## 方法 <a id="&#x65B9;&#x6CD5;"></a>

方法的声明和函数类似，他们的区别是：方法在定义的时候，会在`func`和方法名之间增加一个参数，这个参数就是接收者，这样我们定义的这个方法就和接收者绑定在了一起，称之为这个接收者的方法。

```go
type person struct {
	name string
}

func (p person) String() string{
	return "the person name is "+p.name
}
```

留意例子中，`func`和方法名之间增加的参数`(p person)`,这个就是接收者。现在我们说，类型`person`有了一个`String`方法，现在我们看下如何使用它。

```go
func main() {
	p:=person{name:"张三"}
	fmt.Println(p.String())
}
```

调用的方法非常简单，使用类型的变量进行调用即可，类型变量和方法之前是一个`.`操作符，表示要调用这个类型变量的某个方法的意思。

Go语言里有两种类型的接收者：值接收者和指针接收者。我们上面的例子中，就是使用值类型接收者的示例。

使用值类型接收者定义的方法，在调用的时候，使用的其实是值接收者的一个副本，所以对该值的任何操作，不会影响原来的类型变量。

```go
func main() {
	p:=person{name:"张三"}
	p.modify() //值接收者，修改无效
	fmt.Println(p.String())
}

type person struct {
	name string
}

func (p person) String() string{
	return "the person name is "+p.name
}

func (p person) modify(){
	p.name = "李四"
}
```

以上的例子，打印出来的值还是`张三`，对其进行的修改无效。如果我们使用一个指针作为接收者，那么就会其作用了，因为指针接收者传递的是一个指向原值指针的副本，指针的副本，指向的还是原来类型的值，所以修改时，同时也会影响原来类型变量的值。

```go
func main() {
	p:=person{name:"张三"}
	p.modify() //指针接收者，修改有效
	fmt.Println(p.String())
}

type person struct {
	name string
}

func (p person) String() string{
	return "the person name is "+p.name
}

func (p *person) modify(){
	p.name = "李四"
}
```

只需要改动一下，变成指针的接收者，就可以完成了修改。

> 在调用方法的时候，传递的接收者本质上都是副本，只不过一个是这个值副本，一是指向这个值指针的副本。指针具有指向原有值的特性，所以修改了指针指向的值，也就修改了原有的值。我们可以简单的理解为值接收者使用的是值的副本来调用方法，而指针接收者使用实际的值来调用方法。

在上面的例子中，有没有发现，我们在调用指针接收者方法的时候，使用的也是一个值的变量，并不是一个指针，如果我们使用下面的也是可以的。

```go
p:=person{name:"张三"}
(&p).modify() //指针接收者，修改有效
```

这样也是可以的。如果我们没有这么强制使用指针进行调用，Go的编译器自动会帮我们取指针，以满足接收者的要求。

**同样的，如果是一个值接收者的方法，使用指针也是可以调用的，Go编译器自动会解引用，以满足接收者的要求，比如例子中定义的`String()`方法，也可以这么调用：**

```go
p:=person{name:"张三"}
fmt.Println((&p).String())
```

总之，方法的调用，既可以使用值，也可以使用指针，我们不必要严格的遵守这些，Go语言编译器会帮我们进行自动转义的，这大大方便了我们开发者。

不管是使用值接收者，还是指针接收者，一定要搞清楚类型的本质：对类型进行操作的时候，是要改变当前值，还是要创建一个新值进行返回？这些就可以决定我们是采用值传递，还是指针传递。

