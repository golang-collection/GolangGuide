# rune

在go语言中rune定义如下

```go
type rune = int32
```

使用rune主要是为了区分字符值和整数值，大家知道计算机在底层处理的都是零和一的数字，那么字符串也不另外。string字符串在转换成byte之后其实都是一个数值，我们知道常规的英文字符是ascii码是通过一个字节\( 2^8 其实还有一位是不用的 \)来存储，中国文字、日本文字常用文字就有4000+，通过2^8肯定表达不了，所有可以通过unicode来存储，占用2个字节。而go采用utf-8编码，在utf-8中会将一个码位编码为1-4个字节，所以用int32表示较为合适。

utf-8编码表如下

```text
U+ 0000 ~ U+  007F: 0XXXXXXX
U+ 0080 ~ U+  07FF: 110XXXXX 10XXXXXX
U+ 0800 ~ U+  FFFF: 1110XXXX 10XXXXXX 10XXXXXX
U+10000 ~ U+10FFFF: 11110XXX 10XXXXXX 10XXXXXX 10XXXXXX
```



