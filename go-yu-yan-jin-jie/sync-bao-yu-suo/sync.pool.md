# Sync.Pool

> 转载自：[https://github.com/polaris1119/The-Golang-Standard-Library-by-Example/blob/master/chapter16/16.01.md](https://github.com/polaris1119/The-Golang-Standard-Library-by-Example/blob/master/chapter16/16.01.md)

sync.Pool 其实不适合用来做持久保存的对象池（比如连接池）。它更适合用来做临时对象池，目的是为了降低 GC 的压力。

```go
import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestSyncPool(t *testing.T) {
	pool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Create a new object.")
			return 100
		},
	}

	v := pool.Get().(int)
	fmt.Println(v)
	pool.Put(3)
	runtime.GC() //GC 会清除sync.pool中缓存的对象，所以不适合做为对象池
	v1, _ := pool.Get().(int)
	fmt.Println(v1)
	//新创建的v1并不会自动放入pool，通过put才可以
	v2, _ := pool.Get().(int)
	fmt.Println(v2)
}

func TestSyncPoolInMultiGroutine(t *testing.T) {
	pool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Create a new object.")
			return 10
		},
	}

	pool.Put(100)
	pool.Put(100)
	pool.Put(100)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			fmt.Println(pool.Get())
			wg.Done()
		}(i)
	}
	wg.Wait()
}

```

