# 编码相关

转载自：[http://c.biancheng.net/view/18.html](http://c.biancheng.net/view/18.html)

### UTF-8 和 Unicode 有何区别？

Unicode 与 ASCII 类似，都是一种字符集。  
  
字符集为每个字符分配一个唯一的 ID，我们使用到的所有字符在 Unicode 字符集中都有一个唯一的 ID，例如上面例子中的 a 在 Unicode 与 ASCII 中的编码都是 97。汉字“你”在 Unicode 中的编码为 20320，在不同国家的字符集中，字符所对应的 ID 也会不同。而无论任何情况下，Unicode 中的字符的 ID 都是不会变化的。  
  
UTF-8 是编码规则，将 Unicode 中字符的 ID 以某种方式进行编码，UTF-8 的是一种变长编码规则，从 1 到 4 个字节不等。编码规则如下：

* 0xxxxxx 表示文字符号 0～127，兼容 ASCII 字符集。
* 从 128 到 0x10ffff 表示其他字符。

  
根据这个规则，拉丁文语系的字符编码一般情况下每个字符占用一个字节，而中文每个字符占用 3 个字节。  
  
广义的 Unicode 指的是一个标准，它定义了字符集及编码规则，即 Unicode 字符集和 UTF-8、UTF-16 编码等。

