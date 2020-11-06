# go自带json包

go自带的json解析操作是通过反射进行，故会影响性能。

```go
func BookToJson(book *model.Book) (string, error) {
	str, err := json.Marshal(book)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func JsonToBook(str string) (*model.Book, error) {
	book := &model.Book{}
	err := json.Unmarshal([]byte(str), book)
	if err != nil {
		return nil, err
	}
	return book, nil
}
```

对上面的方法做性能测试

```go
func BenchmarkBookToJson(b *testing.B) {
	for i:=0; i<b.N; i++{
		_, err := BookToJson(&book)
		if err != nil{
			fmt.Println(err)
		}
	}
}
```

运行结果为`3800ns/op`

## Unmarshal与Decoer

json的反序列化方式有两种：

1. Use json.Unmarshal passing the entire response string

   ```text
   // func Unmarshal(data []byte, v interface{}) error
   data, err := ioutil.ReadAll(resp.Body)
   if err == nil && data != nil {
       err = json.Unmarshal(data, value)
   }
   12345
   ```

2. using json.NewDecoder.Decode

   ```text
   // func NewDecoder(r io.Reader) *Decoder
   // func (dec *Decoder) Decode(v interface{}) error
   err = json.NewDecoder(resp.Body).Decode(value)
   123
   ```

这两种方法看似差不多，但有不同的应用场景

* Use json.Decoder if your data is coming from an io.Reader stream, or you need to decode multiple values from a stream of data.

  For the case of reading from an HTTP request, I’d pick json.Decoder since you’re obviously reading from a stream.

* Use json.Unmarshal if you already have the JSON data in memory.

