# easyjson

easyjson库是更高效的json解析方式

## 安装

{% embed url="https://github.com/mailru/easyjson" %}

```go
go get -u github.com/mailru/easyjson/...
```

## 使用

安装结束之后，若已经将GOPATH配置到系统的环境变量中在终端直接执行下面这条命令。

```go
easyjson -all <file>.go
```

其中file是对应的结构体所在的文件。执行完命令后会自动生成一个文件。

![](../../.gitbook/assets/image%20%2819%29.png)

```go
import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson67646b7bDecodeGoJsonModel(in *jlexer.Lexer, out *Book) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "book_id":
			out.BookID = int(in.Int())
		case "title":
			out.Title = string(in.String())
		case "sub_title":
			out.SubTitle = string(in.String())
		case "img":
			out.Img = string(in.String())
		case "author":
			out.Author = string(in.String())
		case "publish":
			out.Publish = string(in.String())
		case "producer":
			out.Producer = string(in.String())
		case "publish_year":
			out.PublishYear = string(in.String())
		case "pages":
			out.Pages = int(in.Int())
		case "price":
			out.Price = float64(in.Float64())
		case "layout":
			out.Layout = string(in.String())
		case "series":
			out.Series = string(in.String())
		case "isbn":
			out.ISBN = string(in.String())
		case "score":
			out.Score = float64(in.Float64())
		case "original_name":
			out.OriginalName = string(in.String())
		case "comments":
			out.Comments = int(in.Int())
		case "comment_url":
			out.CommentUrl = string(in.String())
		case "url":
			out.Url = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson67646b7bEncodeGoJsonModel(out *jwriter.Writer, in Book) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"book_id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.BookID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"sub_title\":"
		out.RawString(prefix)
		out.String(string(in.SubTitle))
	}
	{
		const prefix string = ",\"img\":"
		out.RawString(prefix)
		out.String(string(in.Img))
	}
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		out.String(string(in.Author))
	}
	{
		const prefix string = ",\"publish\":"
		out.RawString(prefix)
		out.String(string(in.Publish))
	}
	{
		const prefix string = ",\"producer\":"
		out.RawString(prefix)
		out.String(string(in.Producer))
	}
	{
		const prefix string = ",\"publish_year\":"
		out.RawString(prefix)
		out.String(string(in.PublishYear))
	}
	{
		const prefix string = ",\"pages\":"
		out.RawString(prefix)
		out.Int(int(in.Pages))
	}
	{
		const prefix string = ",\"price\":"
		out.RawString(prefix)
		out.Float64(float64(in.Price))
	}
	{
		const prefix string = ",\"layout\":"
		out.RawString(prefix)
		out.String(string(in.Layout))
	}
	{
		const prefix string = ",\"series\":"
		out.RawString(prefix)
		out.String(string(in.Series))
	}
	{
		const prefix string = ",\"isbn\":"
		out.RawString(prefix)
		out.String(string(in.ISBN))
	}
	{
		const prefix string = ",\"score\":"
		out.RawString(prefix)
		out.Float64(float64(in.Score))
	}
	{
		const prefix string = ",\"original_name\":"
		out.RawString(prefix)
		out.String(string(in.OriginalName))
	}
	{
		const prefix string = ",\"comments\":"
		out.RawString(prefix)
		out.Int(int(in.Comments))
	}
	{
		const prefix string = ",\"comment_url\":"
		out.RawString(prefix)
		out.String(string(in.CommentUrl))
	}
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix)
		out.String(string(in.Url))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Book) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson67646b7bEncodeGoJsonModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Book) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson67646b7bEncodeGoJsonModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Book) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson67646b7bDecodeGoJsonModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Book) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson67646b7bDecodeGoJsonModel(l, v)
}
```

之后调用其自动生成的MarshJSON等方法即可。

