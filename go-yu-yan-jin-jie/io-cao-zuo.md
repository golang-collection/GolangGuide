# I/O操作

转载自：[http://c.biancheng.net/view/22.html](http://c.biancheng.net/view/22.html)

在Go语言中，几乎所有的[数据结构](http://c.biancheng.net/data_structure/)都围绕接口展开，接口是Go语言中所有数据结构的核心。在实际开发过程中，无论是实现 web 应用程序，还是控制台输入输出，又或者是网络操作，都不可避免的会遇到 I/O 操作。

 Go语言标准库的 bufio 包中，实现了对数据 I/O 接口的缓冲功能。这些功能封装于接口 io.ReadWriter、io.Reader 和 io.Writer 中，并对应创建了 ReadWriter、Reader 或 Writer 对象，在提供缓冲的同时实现了一些文本基本 I/O 操作功能。

##  ReadWriter 对象

 ReadWriter 对象可以对数据 I/O 接口 io.ReadWriter 进行输入输出缓冲操作，ReadWriter 结构定义如下：

 type ReadWriter struct {  
     \*Reader  
     \*Writer  
 }

 默认情况下，ReadWriter 对象中存放了一对 Reader 和 Writer 指针，它同时提供了对数据 I/O 对象的读写缓冲功能。

 可以使用 NewReadWriter\(\) 函数创建 ReadWriter 对象，该函数的功能是根据指定的 Reader 和 Writer 创建一个 ReadWriter 对象，ReadWriter 对象将会向底层 io.ReadWriter 接口写入数据，或者从 io.ReadWriter 接口读取数据。该函数原型声明如下：

 func NewReadWriter\(r \*Reader, w \*Writer\) \*ReadWriter

 在函数 NewReadWriter\(\) 中，参数 r 是要读取的来源 Reader 对象，参数 w 是要写入的目的 Writer 对象。

##  Reader 对象

 Reader 对象可以对数据 I/O 接口 io.Reader 进行输入缓冲操作，Reader 结构定义如下：

 type Reader struct {  
     //contains filtered or unexported fields  
 \)

 默认情况下 Reader 对象没有定义初始值，输入缓冲区最小值为 16。当超出限制时，另创建一个二倍的存储空间。

###  创建 Reader 对象

 可以创建 Reader 对象的函数一共有两个，分别是 NewReader\(\) 和 NewReaderSize\(\)，下面分别介绍。

####  1\) NewReader\(\) 函数

 NewReader\(\) 函数的功能是按照缓冲区默认长度创建 Reader 对象，Reader 对象会从底层 io.Reader 接口读取尽量多的数据进行缓存。该函数原型如下：

 func NewReader\(rd io.Reader\) \*Reader

 其中，参数 rd 是 io.Reader 接口，Reader 对象将从该接口读取数据。

####  2\) NewReaderSize\(\) 函数

 NewReaderSize\(\) 函数的功能是按照指定的缓冲区长度创建 Reader 对象，Reader 对象会从底层 io.Reader 接口读取尽量多的数据进行缓存。该函数原型如下：

 func NewReaderSize\(rd io.Reader, size int\) \*Reader

 其中，参数 rd 是 io.Reader 接口，参数 size 是指定的缓冲区字节长度。

###  操作 Reader 对象

 操作 Reader 对象的方法共有 11 个，分别是 Read\(\)、ReadByte\(\)、ReadBytes\(\)、ReadLine\(\)、ReadRune \(\)、ReadSlice\(\)、ReadString\(\)、UnreadByte\(\)、UnreadRune\(\)、Buffered\(\)、Peek\(\)，下面分别介绍。

####  1\) Read\(\) 方法

 Read\(\) 方法的功能是读取数据，并存放到字节切片 p 中。Read\(\) 执行结束会返回已读取的字节数，因为最多只调用底层的 io.Reader 一次，所以返回的 n 可能小于 len\(p\)，当字节流结束时，n 为 0，err 为 io. EOF。该方法原型如下：

 func \(b \*Reader\) Read\(p \[\]byte\) \(n int, err error\)

 在方法 Read\(\) 中，参数 p 是用于存放读取数据的字节切片。示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    data := []byte("C语言中文网")
    rd := bytes.NewReader(data)
    r := bufio.NewReader(rd)
    var buf [128]byte
    n, err := r.Read(buf[:])
    fmt.Println(string(buf[:n]), n, err)
}
```

 运行结果如下：

 C语言中文网 16

####  2\) ReadByte\(\) 方法

 ReadByte\(\) 方法的功能是读取并返回一个字节，如果没有字节可读，则返回错误信息。该方法原型如下：

 func \(b \*Reader\) ReadByte\(\) \(c byte,err error\)

 示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    data := []byte("Go语言入门教程")
    rd := bytes.NewReader(data)
    r := bufio.NewReader(rd)
    c, err := r.ReadByte()
    fmt.Println(string(c), err)
}
```

 运行结果如下：

 G

####  3\) ReadBytes\(\) 方法

 ReadBytes\(\) 方法的功能是读取数据直到遇到第一个分隔符“delim”，并返回读取的字节序列（包括“delim”）。如果 ReadBytes 在读到第一个“delim”之前出错，它返回已读取的数据和那个错误（通常是 io.EOF）。只有当返回的数据不以“delim”结尾时，返回的 err 才不为空值。该方法原型如下：

 func \(b \*Reader\) ReadBytes\(delim byte\) \(line \[\]byte, err error\)

 其中，参数 delim 用于指定分割字节。示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    data := []byte("C语言中文网, Go语言入门教程")
    rd := bytes.NewReader(data)
    r := bufio.NewReader(rd)
    var delim byte = ','
    line, err := r.ReadBytes(delim)
    fmt.Println(string(line), err)
}
```

 运行结果如下：

 C语言中文网,

####  4\) ReadLine\(\) 方法

 ReadLine\(\) 是一个低级的用于读取一行数据的方法，大多数调用者应该使用 ReadBytes\('\n'\) 或者 ReadString\('\n'\)。ReadLine 返回一行，不包括结尾的回车字符，如果一行太长（超过缓冲区长度），参数 isPrefix 会设置为 true 并且只返回前面的数据，剩余的数据会在以后的调用中返回。

 当返回最后一行数据时，参数 isPrefix 会置为 false。返回的字节切片只在下一次调用 ReadLine 前有效。ReadLine 会返回一个非空的字节切片或一个错误，方法原型如下：

 func \(b \*Reader\) ReadLine\(\) \(line \[\]byte, isPrefix bool, err error\)

 示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    data := []byte("Golang is a beautiful language. \r\n I like it!")
    rd := bytes.NewReader(data)
    r := bufio.NewReader(rd)
    line, prefix, err := r.ReadLine()
    fmt.Println(string(line), prefix, err)
}
```

 运行结果如下：

 Golang is a beautiful language.  false

####  5\) ReadRune\(\) 方法

 ReadRune\(\) 方法的功能是读取一个 UTF-8 编码的字符，并返回其 Unicode 编码和字节数。如果编码错误，ReadRune 只读取一个字节并返回 unicode.ReplacementChar\(U+FFFD\) 和长度 1。该方法原型如下：

 func \(b \*Reader\) ReadRune\(\) \(r rune, size int, err error\)

 示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    data := []byte("C语言中文网")
    rd := bytes.NewReader(data)
    r := bufio.NewReader(rd)
    ch, size, err := r.ReadRune()
    fmt.Println(string(ch), size, err)
}
```

 运行结果如下：

 C 1

####  6\) ReadSlice\(\) 方法

 ReadSlice\(\) 方法的功能是读取数据直到分隔符“delim”处，并返回读取数据的字节切片，下次读取数据时返回的切片会失效。如果 ReadSlice 在查找到“delim”之前遇到错误，它返回读取的所有数据和那个错误（通常是 io.EOF）。

 如果缓冲区满时也没有查找到“delim”，则返回 ErrBufferFull 错误。ReadSlice 返回的数据会在下次 I/O 操作时被覆盖，大多数调用者应该使用 ReadBytes 或者 ReadString。只有当 line 不以“delim”结尾时，ReadSlice 才会返回非空 err。该方法原型如下：

 func \(b \*Reader\) ReadSlice\(delim byte\) \(line \[\]byte, err error\)

 其中，参数 delim 用于指定分割字节。示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    data := []byte("C语言中文网, Go语言入门教程")
    rd := bytes.NewReader(data)
    r := bufio.NewReader(rd)
    var delim byte = ','
    line, err := r.ReadSlice(delim)
    fmt.Println(string(line), err)
    line, err = r.ReadSlice(delim)
    fmt.Println(string(line), err)
    line, err = r.ReadSlice(delim)
    fmt.Println(string(line), err)
}
```

 运行结果如下：

 C语言中文网,  
 Go语言入门教程 EOF  
 EOF

####  7\) ReadString\(\) 方法

 ReadString\(\) 方法的功能是读取数据直到分隔符“delim”第一次出现，并返回一个包含“delim”的字符串。如果 ReadString 在读取到“delim”前遇到错误，它返回已读字符串和那个错误（通常是 io.EOF）。只有当返回的字符串不以“delim”结尾时，ReadString 才返回非空 err。该方法原型如下：

 func \(b \*Reader\) ReadString\(delim byte\) \(line string, err error\)

 其中，参数 delim 用于指定分割字节。示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    data := []byte("C语言中文网, Go语言入门教程")
    rd := bytes.NewReader(data)
    r := bufio.NewReader(rd)
    var delim byte = ','
    line, err := r.ReadString(delim)
    fmt.Println(line, err)
}
```

 运行结果为：

 C语言中文网,

####  8\) UnreadByte\(\) 方法

 UnreadByte\(\) 方法的功能是取消已读取的最后一个字节（即把字节重新放回读取缓冲区的前部）。只有最近一次读取的单个字节才能取消读取。该方法原型如下：

 func \(b \*Reader\) UnreadByte\(\) error

####  9\) UnreadRune\(\) 方法

 UnreadRune\(\) 方法的功能是取消读取最后一次读取的 Unicode 字符。如果最后一次读取操作不是 ReadRune，UnreadRune 会返回一个错误（在这方面它比 UnreadByte 更严格，因为 UnreadByte 会取消上次任意读操作的最后一个字节）。该方法原型如下：

 func \(b \*Reader\) UnreadRune\(\) error

####  10\) Buffered\(\) 方法

 Buffered\(\) 方法的功能是返回可从缓冲区读出数据的字节数, 示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    data := []byte("Go语言入门教程")
    rd := bytes.NewReader(data)
    r := bufio.NewReader(rd)
    var buf [14]byte
    n, err := r.Read(buf[:])
    fmt.Println(string(buf[:n]), n, err)
    rn := r.Buffered()
    fmt.Println(rn)
    n, err = r.Read(buf[:])
    fmt.Println(string(buf[:n]), n, err)
    rn = r.Buffered()
    fmt.Println(rn)
}
```

 运行结果如下：

 Go语言入门 14  
 6  
 教程 6  
 0

####  11\) Peek\(\) 方法

 Peek\(\) 方法的功能是读取指定字节数的数据，这些被读取的数据不会从缓冲区中清除。在下次读取之后，本次返回的字节切片会失效。如果 Peek 返回的字节数不足 n 字节，则会同时返回一个错误说明原因，如果 n 比缓冲区要大，则错误为 ErrBufferFull。该方法原型如下：

 func \(b \*Reader\) Peek\(n int\) \(\[\]byte, error\)

 在方法 Peek\(\) 中，参数 n 是希望读取的字节数。示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    data := []byte("Go语言入门教程")
    rd := bytes.NewReader(data)
    r := bufio.NewReader(rd)
    bl, err := r.Peek(8)
    fmt.Println(string(bl), err)
    bl, err = r.Peek(14)
    fmt.Println(string(bl), err)
    bl, err = r.Peek(20)
    fmt.Println(string(bl), err)
}
```

 运行结果如下：

 Go语言  
 Go语言入门  
 Go语言入门教程

##  Writer 对象

 Writer 对象可以对数据 I/O 接口 io.Writer 进行输出缓冲操作，Writer 结构定义如下：

 type Writer struct {  
     //contains filtered or unexported fields  
 }

 默认情况下 Writer 对象没有定义初始值，如果输出缓冲过程中发生错误，则数据写入操作立刻被终止，后续的写操作都会返回写入异常错误。

###  创建 Writer 对象

 创建 Writer 对象的函数共有两个分别是 NewWriter\(\) 和 NewWriterSize\(\)，下面分别介绍一下。

####  1\) NewWriter\(\) 函数

 NewWriter\(\) 函数的功能是按照默认缓冲区长度创建 Writer 对象，Writer 对象会将缓存的数据批量写入底层 io.Writer 接口。该函数原型如下：

 func NewWriter\(wr io.Writer\) \*Writer

 其中，参数 wr 是 io.Writer 接口，Writer 对象会将数据写入该接口。

####  2\) NewWriterSize\(\) 函数

 NewWriterSize\(\) 函数的功能是按照指定的缓冲区长度创建 Writer 对象，Writer 对象会将缓存的数据批量写入底层 io.Writer 接口。该函数原型如下：

 func NewWriterSize\(wr io.Writer, size int\) \*Writer

 其中，参数 wr 是 io.Writer 接口，参数 size 是指定的缓冲区字节长度。

###  操作 Writer 对象

 操作 Writer 对象的方法共有 7 个，分别是 Available\(\)、Buffered\(\)、Flush\(\)、Write\(\)、WriteByte\(\)、WriteRune\(\) 和 WriteString\(\) 方法，下面分别介绍。

####  1\) Available\(\) 方法

 Available\(\) 方法的功能是返回缓冲区中未使用的字节数，该方法原型如下：

 func \(b \*Writer\) Available\(\) int

 示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    wr := bytes.NewBuffer(nil)
    w := bufio.NewWriter(wr)
    p := []byte("C语言中文网")
    fmt.Println("写入前未使用的缓冲区为：", w.Available())
    w.Write(p)
    fmt.Printf("写入%q后，未使用的缓冲区为：%d\n", string(p), w.Available())
}
```

 运行结果如下：

 写入前未使用的缓冲区为： 4096  
 写入"C语言中文网"后，未使用的缓冲区为：4080

####  2\) Buffered\(\) 方法

 Buffered\(\) 方法的功能是返回已写入当前缓冲区中的字节数，该方法原型如下：

 func \(b \*Writer\) Buffered\(\) int

 示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    wr := bytes.NewBuffer(nil)
    w := bufio.NewWriter(wr)
    p := []byte("C语言中文网")
    fmt.Println("写入前未使用的缓冲区为：", w.Buffered())
    w.Write(p)
    fmt.Printf("写入%q后，未使用的缓冲区为：%d\n", string(p), w.Buffered())
    w.Flush()
    fmt.Println("执行 Flush 方法后，写入的字节数为：", w.Buffered())
}
```

 该例测试结果为：

 写入前未使用的缓冲区为： 0  
 写入"C语言中文网"后，未使用的缓冲区为：16  
 执行 Flush 方法后，写入的字节数为： 0

####  3\) Flush\(\) 方法

 Flush\(\) 方法的功能是把缓冲区中的数据写入底层的 io.Writer，并返回错误信息。如果成功写入，error 返回 nil，否则 error 返回错误原因。该方法原型如下：

 func \(b \*Writer\) Flush\(\) error

 示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    wr := bytes.NewBuffer(nil)
    w := bufio.NewWriter(wr)
    p := []byte("C语言中文网")
    w.Write(p)
    fmt.Printf("未执行 Flush 缓冲区输出 %q\n", string(wr.Bytes()))
    w.Flush()
    fmt.Printf("执行 Flush 后缓冲区输出 %q\n", string(wr.Bytes()))
}
```

 运行结果如下：

 未执行 Flush 缓冲区输出 ""  
 执行 Flush 后缓冲区输出 "C语言中文网"

####  4\) Write\(\) 方法

 Write\(\) 方法的功能是把字节切片 p 写入缓冲区，返回已写入的字节数 nn。如果 nn 小于 len\(p\)，则同时返回一个错误原因。该方法原型如下：

 func \(b \*Writer\) Write\(p \[\]byte\) \(nn int, err error\)

 其中，参数 p 是要写入的字节切片。示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    wr := bytes.NewBuffer(nil)
    w := bufio.NewWriter(wr)
    p := []byte("C语言中文网")
    n, err := w.Write(p)
    w.Flush()
    fmt.Println(string(wr.Bytes()), n, err)
}
```

 运行结果如下：

 C语言中文网 16

####  5\) WriteByte\(\) 方法

 WriteByte\(\) 方法的功能是写入一个字节，如果成功写入，error 返回 nil，否则 error 返回错误原因。该方法原型如下：

 func \(b \*Writer\) WriteByte\(c byte\) error

 其中，参数 c 是要写入的字节数据，比如 ASCII 字符。示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    wr := bytes.NewBuffer(nil)
    w := bufio.NewWriter(wr)
    var c byte = 'G'
    err := w.WriteByte(c)
    w.Flush()
    fmt.Println(string(wr.Bytes()), err)
}
```

 运行结果如下：

 G

####  6\) WriteRune\(\) 方法

 WriteRune\(\) 方法的功能是以 UTF-8 编码写入一个 Unicode 字符，返回写入的字节数和错误信息。该方法原型如下：

 func \(b \*Writer\) WriteRune\(r rune\) \(size int,err error\)

 其中，参数 r 是要写入的 Unicode 字符。示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    wr := bytes.NewBuffer(nil)
    w := bufio.NewWriter(wr)
    var r rune = 'G'
    size, err := w.WriteRune(r)
    w.Flush()
    fmt.Println(string(wr.Bytes()), size, err)
}
```

 该例测试结果为：

 G 1

####  7\) WriteString\(\) 方法

 WriteString\(\) 方法的功能是写入一个字符串，并返回写入的字节数和错误信息。如果返回的字节数小于 len\(s\)，则同时返回一个错误说明原因。该方法原型如下：

 func \(b \*Writer\) WriteString\(s string\) \(int, error\)

 其中，参数 s 是要写入的字符串。示例代码如下：

```text
package main

import (
    "bufio"
    "bytes"
    "fmt"
)

func main() {
    wr := bytes.NewBuffer(nil)
    w := bufio.NewWriter(wr)
    s := "C语言中文网"
    n, err := w.WriteString(s)
    w.Flush()
    fmt.Println(string(wr.Bytes()), n, err)
}
```

 运行结果如下：

 C语言中文网 16

