# json操作

## 什么是json

JSON\(JavaScript Object Notation\) 是一种轻量级的数据交换格式。它使得人们很容易的进行阅读和编写。同时也方便了机器进行解析和生成。它是基于 [JavaScript Programming Language](http://www.crockford.com/javascript) , [Standard ECMA-262 3rd Edition - December 1999](http://www.ecma-international.org/publications/files/ecma-st/ECMA-262.pdf) 的一个子集。 JSON采用完全独立于程序语言的文本格式，但是也使用了类C语言的习惯（包括C, C++, C\#, Java, JavaScript, Perl, Python等）。这些特性使JSON成为理想的数据交换语言。

### **JSON基于两种结构：**

JSON有两种结构，json简单说就是javascript中的对象和数组，所以这两种结构就是对象和数组两种结构，通过这两种结构可以表示各种复杂的结构

* 1、对象：对象在js中表示为“{}”括起来的内容，数据结构为 {key：value,key：value,...}的键值对的结构，在面向对象的语言中，key为对象的属性，value为对应的属性值，所以很容易理解，取值方法为 对象.key 获取属性值，这个属性值的类型可以是 数字、字符串、数组、对象几种。
* 2、数组：数组在js中是中括号“\[\]”括起来的内容，数据结构为 \["java","javascript","vb",...\]，取值方式和所有语言中一样，使用索引获取，字段值的类型可以是 数字、字符串、数组、对象几种。

```javascript

{
    "animals": {
        "dog": [
            {
                "name": "Rufus",
                "age":15
            },
            {
                "name": "Marty",
                "age": null
            }
        ]
    }
}
```

经过对象、数组2种结构就可以组合成复杂的数据结构了。

“名称/值”对的集合（A collection of name/value pairs）。不同的编程语言中，它被理解为对象（object），纪录（record），结构（struct），字典（dictionary），哈希表（hash table），有键列表（keyed list），或者关联数组 （associative array）。 值的有序列表（An ordered list of values）。在大部分语言中，它被实现为数组（array），矢量（vector），列表（list），序列（sequence）。这些都是常见的数据结构。目前，绝大部分编程语言都以某种形式支持它们。这使得在各种编程语言之间交换同样格式的数据成为可能。

### **JSON具有以下这些形式：**

对象（object） 是一个无序的“‘名称/值’对”集合。一个对象以“{”（左括号）开始，“}”（右括号）结束。每个“名称”后跟一个“:”（冒号）；“‘名称/值’ 对”之间使用“,”（逗号）分隔。数组（array） 是值（value）的有序集合。一个数组以“\[”（左中括号）开始，“\]”（右中括号）结束。值之间使用“,”（逗号）分隔。值（value） 可以是双引号括起来的字符串（`string`）、数值\(`number`\)、`true`、`false`、 `null`、对象（object）或者数组（array）。这些结构可以嵌套。字符串（string） 是由双引号包围的任意数量Unicode字符的集合，使用反斜线转义。一个字符（character）即一个单独的字符串（character string）。JSON的字符串（string）与C或者Java的字符串非常相似。数值（number） 也与C或者Java的数值非常相似。只是JSON的数值没有使用八进制与十六进制格式。同时，可以在任意标记之间添加空白。

### **访问数据**

尽管看起来不明显，但是上面的长字符串实际上只是一个数组；将这个数组放进 JavaScript变量之后，就可以很轻松地访问它。实际上，只需用点号表示法来表示数组元素。所以，要想访问 programmers 列表的第一个条目的姓氏，只需在 JavaScript 中使用下面这样的代码：`people.programmers[0].lastName;`  
注意，数组索引是从零开始的。所以，这行代码首先访问 people变量中的数据；然后移动到称为 programmers的条目，再移动到第一个记录（\[0\]）；最后，访问 lastName键的值。结果是字符串值 “McLaughlin”。 下面是使用同一变量的几个示例。

`people.authors[1].genre//Valueis"fantasy"  
people.musicians[3].lastName//Undefined.Thisreferstothefourthentry,andthereisn'tone  
people.programmers[2].firstName//Valueis"Elliotte"`  


利用这样的语法，可以处理任何 JSON 格式的数据，而不需要使用任何额外的 JavaScript 工具包或 API。

## **和XML的比较**

### **可读性**

JSON和XML的可读性可谓不相上下，一边是简易的语法，一边是规范的标签形式，很难分出胜负。

### **可扩展性**

XML天生有很好的扩展性，JSON当然也有，没有什么是XML可以扩展而JSON却不能扩展的。不过JSON在Javascript主场作战，可以存储Javascript复合对象，有着xml不可比拟的优势。

### **编码难度**

XML有丰富的编码工具，比如Dom4j、JDom等，JSON也有提供的工具。无工具的情况下，相信熟练的开发人员一样能很快的写出想要的xml文档和JSON字符串，不过，xml文档要多很多结构上的字符。

### **解码难度**

* XML的解析方式有两种：
* 一是通过文档模型解析，也就是通过父标签索引出一组标记。例如：xmlData.getElementsByTagName\("tagName"\)，但是这样是要在预先知道文档结构的情况下使用，无法进行通用的封装。
* 另外一种方法是遍历节点（document 以及 childNodes）。这个可以通过递归来实现，不过解析出来的数据仍旧是形式各异，往往也不能满足预先的要求。

凡是这样可扩展的结构数据解析起来一定都很困难。

JSON也同样如此。如果预先知道JSON结构的情况下，使用JSON进行数据传递简直是太美妙了，可以写出很实用美观可读性强的代码。如果你是纯粹的前台开发人员，一定会非常喜欢JSON。但是如果你是一个应用开发人员，就不是那么喜欢了，毕竟xml才是真正的结构化标记语言，用于进行数据传递。

而如果不知道JSON的结构而去解析JSON的话，那简直是噩梦。费时费力不说，代码也会变得冗余拖沓，得到的结果也不尽人意。但是这样也不影响众多前台开发人员选择JSON。因为json.js中的toJSONString\(\)就可以看到JSON的字符串结构。当然不是使用这个字符串，这样仍旧是噩梦。常用JSON的人看到这个字符串之后，就对JSON的结构很明了了，就更容易的操作JSON。

以上是在Javascript中仅对于数据传递的xml与JSON的解析。在Javascript地盘内，JSON毕竟是主场作战，其优势当然要远远优越于xml。如果JSON中存储Javascript复合对象，而且不知道其结构的话，我相信很多程序员也一样是哭着解析JSON的。

除了上述之外，JSON和XML还有另外一个很大的区别在于有效数据率。JSON作为数据包格式传输的时候具有更高的效率，这是因为JSON不像XML那样需要有严格的闭合标签，这就让有效数据量与总数据包比大大提升，从而减少同等数据流量的情况下，网络的传输压力。

## 结构体定义

下面我们通过解析这个Book结构体来实现json的解析。

```go
type Book struct {
	BookID       int     `gorm:"column:book_id" gorm:"PRIMARY_KEY" json:"book_id"`
	Title        string  `gorm:"column:title" json:"title"`
	SubTitle     string  `gorm:"column:sub_title" json:"sub_title"`
	Img          string  `gorm:"column:img" json:"img"`
	Author       string  `gorm:"column:author" json:"author"`
	Publish      string  `gorm:"column:publish" json:"publish"`
	Producer     string  `gorm:"column:producer" json:"producer"`
	PublishYear  string  `gorm:"column:publish_year" gorm:"type:date" json:"publish_year"`
	Pages        int     `gorm:"column:pages" json:"pages"`
	Price        float64 `gorm:"column:price" json:"price"`
	Layout       string  `gorm:"column:layout" json:"layout"`
	Series       string  `gorm:"column:series" json:"series"`
	ISBN         string  `gorm:"column:isbn" json:"isbn"`
	Score        float64 `gorm:"column:score" json:"score"`
	OriginalName string  `gorm:"column:original_name" json:"original_name"`
	Comments     int     `gorm:"column:comments" json:"comments"`
	CommentUrl   string  `gorm:"column:comment_url" json:"comment_url"`
	Url          string  `gorm:"column:url" json:"url"`
}

func (book Book) String() string {
	return fmt.Sprintf("book_id: %d\n"+
		"title: %s\n"+
		"sub_title: %s\n"+
		"img: %s\n"+
		"author: %s\n"+
		"publish: %s\n"+
		"producer: %s\n"+
		"publish_year: %s\n"+
		"pages: %d\n"+
		"price: %f\n"+
		"layout: %s\n"+
		"series: %s\n"+
		"isbn: %s\n"+
		"score: %f\n"+
		"original_name: %s\n"+
		"comments: %d\n"+
		"comment_url: %s\n"+
		"url: %s",
		book.BookID, book.Title, book.SubTitle, book.Img,
		book.Author, book.Publish, book.Producer, book.PublishYear,
		book.Pages, book.Price, book.Layout, book.Series, book.ISBN,
		book.Score, book.OriginalName, book.Comments, book.CommentUrl, book.Url)
}
```

