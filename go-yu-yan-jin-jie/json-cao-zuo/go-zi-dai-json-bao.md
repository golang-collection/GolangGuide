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

