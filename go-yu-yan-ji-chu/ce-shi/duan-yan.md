# 断言

## 安装

{% embed url="https://github.com/stretchr/testify" %}

```go
go get github.com/stretchr/testify
```

## assert包

assert包提供了一些有用的方法，允许您在Go中编写更好的测试代码。

* 打印友好，易于阅读的失败描述
* 允许非常易读的代码
* 可以选择使用消息注释每个断言

```go
package yours

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {

  // assert equality
  assert.Equal(t, 123, 123, "they should be equal")

  // assert inequality
  assert.NotEqual(t, 123, 456, "they should not be equal")

  // assert for nil (good for errors)
  assert.Nil(t, object)

  // assert for not nil (good when you expect something)
  if assert.NotNil(t, object) {

    // now we know that object isn't nil, we are safe to make
    // further assertions without causing any errors
    assert.Equal(t, "Something", object.Value)

  }

}
```



