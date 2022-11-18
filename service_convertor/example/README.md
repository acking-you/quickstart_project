## service_convertor example
### tag的支持与使用
由于service还需要生成对应的vo和to结构体，而vo是面向前端的接收数据，而且由于使用gin框架，默认是支持打tag来进行信息验证的。所以service层提供tag来自动生成也很有必要。

#### 通过tag筛选字段以及自定义添加tag

具体的tag形式非常单一，没有dao层那么复杂，这个tag仅仅只有一个作用就是识别key和value来填补vo的tag。
形如：`service:"binding(required,email)"`

举个例子，如果在user的username字段添加如下tag：`service:"binding(required,email);form(user_name)"`，
那么将会在userVO结构体中添加username字段作为成员，并生成tag：`binding:"required,email" form:"user_name"`。
所以只是起到VO结构体字段的取舍以及tag的内容生成的作用，没有其他作用。

### config配置
> 具体的配置信息请查看源代码注释：[code](../config.go)

### 简单示例

> 原始配置
```go
package main

import (
	"github.com/ACking-you/quickstart_project/service_convertor"
)

type User struct {
	Id       int
	Username string `service:"binding(required,email)"`
}

func autoService() {
	config := service_convertor.DefaultConfig().
		//是否输出内容到控制台
		EnableDebug(true)
	convert := service_convertor.NewStruct2Service(config)

	err := convert.AutoMigrate(&User{}).Run()
	if err != nil {
		panic(err)
	}
}


func main() {
	autoService()
}

```

> 代码生成

```go
package vo

// UserVO 默认会自动生成json和form的tag，可以在config中取消
type UserVO struct {
	Username string `binding:"required,email" json:"username" form:"username"`
}
```