# 读取文件对比

Go 语言在进行文件操作的时候，可以有多种方法。最常见的比如直接对文件本身进行`Read`和`Write`；除此之外，还可以使用`bufio`库的流式处理以及分片式处理；如果文件较小，使用`ioutil`也不失为一种方法。

面对这么多的文件处理的方式，那么初学者可能就会有困惑：我到底该用那种？它们之间有什么区别？笔者试着从文件读取来对 go 语言的几种文件处理方式进行分析。

## os.File、bufio、ioutil 比较

### 效率测试

文件的读取效率是所有开发者都会关心的话题，尤其是当文件特别大的时候。为了尽可能的展示这三者对文件读取的性能，我准备了三个文件，分别为`small.txt`，`midium.txt`、`large.txt`，分别对应 KB 级别、MB 级别和 GB 级别。

这三个文件大小分别为 4KB、21MB、1GB。其中内容是比较常规的`json`格式的文本。测试代码如下:

```go
//使用File自带的Read
func read1(filename string) int {
    fi, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer fi.Close()
    buf := make([]byte, 4096)
    var nbytes int
    for {
        n, err := fi.Read(buf)
        if err != nil && err != io.EOF {
            panic(err)
        }
        if n == 0 {
            break
        }
        nbytes += n    }
    return nbytes
}
```

`read1`函数使用的是`os`库对文件进行直接操作，为了确定确实都到了文件内容，并将读到的大小字节数返回。

```go
//使用bufio
func read2(filename string) int {
    fi, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer fi.Close()
    buf := make([]byte, 4096)
    var nbytes int
    rd := bufio.NewReader(fi)
    for {
        n, err := rd.Read(buf)
        if err != nil && err != io.EOF {
            panic(err)
        }
        if n == 0 {
            break
        }
        nbytes += n
    }
    return nbytes
}
```

`read2`函数使用的是`bufio`库，操作`NewReader`对文件进行流式处理，和前面一样，为了确定确实都到了文件内容，并将读到的大小字节数返回。

```go
//使用ioutil
func read3(filename string) int {
    fi, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer fi.Close()
    fd, err := ioutil.ReadAll(fi)
    nbytes := len(fd)
    return nbytes
}
```

`read3`函数是使用`ioutil`库进行文件读取，这种方式比较暴力，直接将文件内容一次性全部读到内存中，然后对内存中的文件内容进行相关的操作。我们使用如下的测试代码进行测试：

```go
func testfile1(filename string) {
    fmt.Printf("============test1 %s ===========\n", filename)
    start := time.Now()
    size1 := read1(filename)
    t1 := time.Now()
    fmt.Printf("Read 1 cost: %v, size: %d\n", t1.Sub(start), size1)
    size2 := read2(filename)
    t2 := time.Now()
    fmt.Printf("Read 2 cost: %v, size: %d\n", t2.Sub(t1), size2)
    size3 := read3(filename)
    t3 := time.Now()
    fmt.Printf("Read 3 cost: %v, size: %d\n", t3.Sub(t2), size3)
}
```

在`main`函数中调用如下：

```go
func main() {
    testfile1("small.txt")
    testfile1("midium.txt")
    testfile1("large.txt")
    // testfile2("small.txt")
    // testfile2("midium.txt")
    // testfile2("large.txt")
}
```

测试结果如下所示：

从以上结果可知：

* 当文件较小（KB 级别）时，ioutil &gt; bufio &gt; os。
* 当文件大小比较常规（MB 级别）时，三者差别不大，但 bufio 又是已经显现出来。
* 当文件较大（GB 级别）时，bufio &gt; os &gt; ioutil。

### 原因分析

为什么会出现上面的不同结果？其实`ioutil`最好理解，当文件较小时，`ioutil`使用`ReadAll`函数将文件中所有内容直接读入内存，只进行了一次 io 操作，但是`os`和`bufio`都是进行了多次读取，才将文件处理完，所以`ioutil`肯定要快于`os`和`bufio`的。

但是随着文件的增大，达到接近 GB 级别时，`ioutil`直接读入内存的弊端就显现出来，要将 GB 级别的文件内容全部读入内存，也就意味着要开辟一块 GB 大小的内存用来存放文件数据，这对内存的消耗是非常大的，因此效率就慢了下来。

如果文件继续增大，达到 3GB 甚至以上，`ioutil`这种读取方式就完全无能为力了。（一个单独的进程空间为 4GB，真正存放数据的堆区和栈区更是远远小于 4GB）。

而`os`为什么在面对大文件时，效率会低于`bufio`？通过查看`bufio`的`NewReader`源码不难发现，在`NewReader`里，默认为我们提供了一个大小为 4096 的缓冲区，所以系统调用会每次先读取 4096 字节到缓冲区，然后`rd.Read`会从缓冲区去读取。

```go
const (
    defaultBufSize = 4096
)

func NewReader(rd io.Reader) *Reader {
    return NewReaderSize(rd, defaultBufSize)
}

func NewReaderSize(rd io.Reader, size int) *Reader {
    // Is it already a Reader?
    b, ok := rd.(*Reader)
    if ok && len(b.buf) >= size {
        return b
    }
    if size < minReadBufferSize {
        size = minReadBufferSize
    }
    r := new(Reader)
    r.reset(make([]byte, size), rd)
    return r
}
```

而`os`因为少了这一层缓冲区，每次读取，都会执行系统调用，因此内核频繁的在用户态和内核态之间切换，而这种切换，也是需要消耗的，故而会慢于`bufio`的读取方式。

笔者翻阅网上资料，关于缓冲，有**内核中的缓冲**和**进程中的缓冲**两种，其中，内核中的缓冲是内核提供的，即系统对磁盘提供一个缓冲区，不管有没有提供进程中的缓冲，内核缓冲都是存在的。

而进程中的缓冲是对输入输出流做了一定的改进，提供的一种流缓冲，它在读写操作发生时，先将数据存入流缓冲中，只有当流缓冲区满了或者刷新（如调用`flush`函数）时，才将数据取出，送往内核缓冲区，它起到了一定的保护内核的作用。

因此，我们不难发现，`os`是典型的内核中的缓冲，而`bufio`和`ioutil`都属于进程中的缓冲。

### 总结

**当读取小文件时，使用`ioutil`效率明显优于`os`和`bufio`，但如果是大文件，`bufio`读取会更快。**

## 读取一行数据

前面简要分析了 go 语言三种不同读取文件方式之间的区别。但实际的开发中，我们对文件的读取往往是以行为单位的，即每次读取一行进行处理。

go 语言并没有像 C 语言一样给我们提供好了类似于`fgets`这样的函数可以正好读取一行内容，因此，需要自己去实现。

从前面的对比分析可以知道，无论是处理大文件还是小文件，`bufio`始终是最为平滑和高效的，因此我们考虑使用`bufio`库进行处理。翻阅`bufio`库的源码，发现可以使用如下几种方式进行读取一行文件的处理：

* `ReadBytes`
* `ReadString`
* `ReadSlice`
* `ReadLine`

## 效率测试

在讨论这四种读取一行文件操作的函数之前，仍然做一下效率测试。测试代码如下：

```go
func readline1(filename string) {
    fi, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer fi.Close()
    rd := bufio.NewReader(fi)
    for {
        _, err := rd.ReadBytes('\n')
        if err != nil || err == io.EOF {
            break
        }
    }
}

func readline2(filename string) {
    fi, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer fi.Close()
    rd := bufio.NewReader(fi)
    for {
        _, err := rd.ReadString('\n')
        if err != nil || err == io.EOF {
            break
        }
    }
}

func readline3(filename string) {
    fi, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer fi.Close()
    rd := bufio.NewReader(fi)
    for {
        _, err := rd.ReadSlice('\n')
        if err != nil || err == io.EOF {
            break
        }
    }
}

func readline4(filename string) {
    fi, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer fi.Close()
    rd := bufio.NewReader(fi)
    for {
        _, _, err := rd.ReadLine()
        if err != nil || err == io.EOF {
            break
        }
    }
}
```

可以看到，这四种操作方式，无论是函数调用，还是函数返回值的处理，其实都是大同小异的。但通过测试效率，则可以看出它们之间的区别。我们使用下面的测试代码：

```go
func testfile2(filename string) {
    fmt.Printf("============test2 %s ===========\n", filename)
    start := time.Now()
    readline1(filename)
    t1 := time.Now()
    fmt.Printf("Readline 1 cost: %v\n", t1.Sub(start))
    readline2(filename)
    t2 := time.Now()
    fmt.Printf("Readline 2 cost: %v\n", t2.Sub(t1))
    readline3(filename)
    t3 := time.Now()
    fmt.Printf("Readline 3 cost: %v\n", t3.Sub(t2))
    readline4(filename)
    t4 := time.Now()
    fmt.Printf("Readline 4 cost: %v\n", t4.Sub(t3))
}
```

在`main`函数中调用如下：

```go
func main() {
    // testfile1("small.txt")
    // testfile1("midium.txt")
    // testfile1("large.txt")
    testfile2("small.txt")
    testfile2("midium.txt")
    testfile2("large.txt")
}
```

运行结果如下所示：

通过现象，除了`small.txt`之外，大致可以分为两组：

* `ReadBytes`对小文件处理效率最差
* 在处理大文件时，`ReadLine`和`ReadSlice`效率相近，要明显快于`ReadString`和`ReadBytes`。

### 原因分析

为什么会出现上面的现象，不防从源码层面进行分析。通过阅读源码，我们发现这四个函数之间存在这样一个关系：

* `ReadLine` &lt;- \(调用\) `ReadSlice`
* `ReadString` &lt;- \(调用\)`ReadBytes`&lt;-\(调用\)`ReadSlice`

既然如此，那为什么在处理大文件时，`ReadLine`效率要明显高于`ReadBytes`呢？

首先，我们要知道，`ReadSlice`是切片式读取，即根据分隔符去进行切片。通过源码发现，`ReadLine`只是在切片读取的基础上，对换行符`\n`和`\r\n`做了一些处理：

```go
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error) {
    line, err = b.ReadSlice('\n')
    if err == ErrBufferFull {
        // Handle the case where "\r\n" straddles the buffer.
        if len(line) > 0 && line[len(line)-1] == '\r' {
            // Put the '\r' back on buf and drop it from line.
            // Let the next call to ReadLine check for "\r\n".
            if b.r == 0 {
                // should be unreachable
                panic("bufio: tried to rewind past start of buffer")
            }
            b.r--
            line = line[:len(line)-1]
        }
        return line, true, nil
    }    if len(line) == 0 {
        if err != nil {
            line = nil
        }
        return
    }
    err = nil    if line[len(line)-1] == '\n' {
        drop := 1
        if len(line) > 1 && line[len(line)-2] == '\r' {
            drop = 2
        }
        line = line[:len(line)-drop]
    }
    return
}
```

而`ReadBytes`则是通过`append`先将读取的内容暂存到`full`数组中，最后再`copy`出来，`append`和`copy`都是要消耗内存和 io 的，因此效率自然就慢了。其源码如下所示：

```go
func (b *Reader) ReadBytes(delim byte) ([]byte, error) {
    // Use ReadSlice to look for array,
    // accumulating full buffers.
    var frag []byte
    var full [][]byte
    var err error
    n := 0
    for {
        var e error
        frag, e = b.ReadSlice(delim)
        if e == nil { // got final fragment
            break
        }
        if e != ErrBufferFull { // unexpected error
            err = e
            break
        }        // Make a copy of the buffer.
        buf := make([]byte, len(frag))
        copy(buf, frag)
        full = append(full, buf)
        n += len(buf)
    }    n += len(frag)    // Allocate new buffer to hold the full pieces and the fragment.
    buf := make([]byte, n)
    n = 0
    // Copy full pieces and fragment in.
    for i := range full {
        n += copy(buf[n:], full[i])
    }
    copy(buf[n:], frag)
    return buf, err
}
```

### 总结

读取文件中一行内容时，`ReadSlice`和`ReadLine`性能优于`ReadBytes`和`ReadString`，但由于`ReadLine`对换行的处理更加全面（兼容`\n`和`\r\n`换行），因此，实际开发过程中，建议使用`ReadLine`函数。

> 作者：禹鼎侯 
>
> 原文链接：https://segmentfault.com/a/1190000023691973

