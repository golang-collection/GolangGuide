# Once

有的时候，我们多个 goroutine 都要过一个操作，但是这个操作我只希望被执行一次，这个时候 Once 就上场了。比如下面的例子 :

```text
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    var once sync.Once
    onceBody := func() {
        fmt.Println("Only once")
    }
    for i := 0; i < 10; i++ {
        go func() {
            once.Do(onceBody)
        }()
    }
    time.Sleep(3e9)
}
```

只会打出一次 "Only once"。

