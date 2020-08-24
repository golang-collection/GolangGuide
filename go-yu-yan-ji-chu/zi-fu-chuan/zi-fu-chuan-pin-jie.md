# 字符串拼接

## SPrintf

```go
const numbers = 100

func BenchmarkSprintf(b *testing.B) {
	b.ResetTimer()
	for idx := 0; idx < b.N; idx++ {
		var s string
		for i := 0; i < numbers; i++ {
			s = fmt.Sprintf("%v%v", s, i)
		}
	}
	b.StopTimer()
}
```

## +拼接

```go
func BenchmarkStringAdd(b *testing.B) {
	b.ResetTimer()
	for idx := 0; idx < b.N; idx++ {
		var s string
		for i := 0; i < numbers; i++ {
			s += strconv.Itoa(i)
		}

	}
	b.StopTimer()
}
```

## bytes.Buffer

```go
func BenchmarkBytesBuf(b *testing.B) {
	b.ResetTimer()
	for idx := 0; idx < b.N; idx++ {
		var buf bytes.Buffer
		for i := 0; i < numbers; i++ {
			buf.WriteString(strconv.Itoa(i))
		}
		_ = buf.String()
	}
	b.StopTimer()
}
```

## strings.Builder拼接

```go
func BenchmarkStringBuilder(b *testing.B) {
	b.ResetTimer()
	for idx := 0; idx < b.N; idx++ {
		var builder strings.Builder
		for i := 0; i < numbers; i++ {
			builder.WriteString(strconv.Itoa(i))

		}
		_ = builder.String()
	}
	b.StopTimer()
}
```

## 对比

```go
BenchmarkSprintf-8         	   68277	     18431 ns/op
BenchmarkStringBuilder-8   	 1302448	       922 ns/op
BenchmarkBytesBuf-8        	  884354	      1264 ns/op
BenchmarkStringAdd-8       	  208486	      5703 ns/op
```

可以看到通过strings.Builder拼接字符串是最高效的。

