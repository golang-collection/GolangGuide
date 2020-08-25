# BDD

## 安装

{% embed url="https://github.com/smartystreets/goconvey" %}

## 使用

```go
import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given 2 even numbers", t, func() {
		a := 3
		b := 4

		Convey("When add the two numbers", func() {
			c := a + b

			Convey("Then the result is still even", func() {
				So(c%2, ShouldEqual, 0)
			})
		})
	})
}
```

## 浏览器观察

```go
$GOPATH/bin/goconvey
```

打开浏览器：[http://localhost:8080](http://localhost:8080)

![](../../.gitbook/assets/image%20%2821%29.png)

## BDD

{% embed url="https://medium.com/javascript-scene/behavior-driven-development-bdd-and-functional-testing-62084ad7f1f2" %}

{% embed url="https://blog.csdn.net/chancein007/article/details/77543494" %}



