# 正则表达式

转载自：[https://github.com/unknwon/the-way-to-go\_ZH\_CN/blob/master/eBook/09.2.md](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/09.2.md) by unknown

## regexp 包

在下面的程序里，我们将在字符串中对正则表达式模式（pattern）进行匹配。

如果是简单模式，使用 `Match` 方法便可：

```go
ok, _ := regexp.Match(pat, []byte(searchIn))
```

变量 ok 将返回 true 或者 false,我们也可以使用 `MatchString`：

```go
ok, _ := regexp.MatchString(pat, searchIn)
```

更多方法中，必须先将正则模式通过 `Compile` 方法返回一个 Regexp 对象。然后我们将掌握一些匹配，查找，替换相关的功能。

```go
package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func main() {
	//目标字符串
	searchIn := "John: 2578.34 William: 4567.23 Steve: 5632.18"
	pat := "[0-9]+.[0-9]+" //正则

	f := func(s string) string{
    	v, _ := strconv.ParseFloat(s, 32)
    	return strconv.FormatFloat(v * 2, 'f', 2, 32)
	}

	if ok, _ := regexp.Match(pat, []byte(searchIn)); ok {
    fmt.Println("Match Found!")
	}

	re, _ := regexp.Compile(pat)
	//将匹配到的部分替换为"##.#"
	str := re.ReplaceAllString(searchIn, "##.#")
	fmt.Println(str)
	//参数为函数时
	str2 := re.ReplaceAllStringFunc(searchIn, f)
	fmt.Println(str2)
}
```

输出结果：

```go
Match Found!
John: ##.# William: ##.# Steve: ##.#
John: 5156.68 William: 9134.46 Steve: 11264.36
```

`Compile` 函数也可能返回一个错误，我们在使用时忽略对错误的判断是因为我们确信自己正则表达式是有效的。当用户输入或从数据中获取正则表达式的时候，我们有必要去检验它的正确性。另外我们也可以使用 `MustCompile` 方法，它可以像 `Compile` 方法一样检验正则的有效性，但是当正则不合法时程序将 panic。

