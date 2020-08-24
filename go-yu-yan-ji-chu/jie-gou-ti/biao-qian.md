# 标签

结构体中的字段除了有名字和类型外，还可以有一个可选的标签（tag）：它是一个附属于字段的字符串，可以是文档或其他的重要标记。标签的内容不可以在一般的编程中使用，一般使用在orm或者json中。只有包 `reflect` 能获取它。它可以在运行时自省类型、属性和方法，比如：在一个变量上调用 `reflect.TypeOf()` 可以获取变量的正确类型，如果变量是一个结构体类型，就可以通过 Field 来索引结构体的字段，然后就可以使用 Tag 属性。

给结构体打tag例如下面这样。

```go
type Activity struct {
	ActivityId int `gorm:"column:activity_id" gorm:"PRIMARY_KEY" json:"activity_id"`
	PlayerId int `gorm:"column:player_id" json:"player_id"`
	DeviceId int `gorm:"column:device_id" json:"device_id"`
	EventDate string `gorm:"column:event_date" json:"event_date"`
	GamesPlayed int `gorm:"column:games_played" json:"games_played"`
}
```

更多内容可以参考：[https://www.flysnow.org/2017/06/25/go-in-action-struct-tag.html](https://www.flysnow.org/2017/06/25/go-in-action-struct-tag.html)

