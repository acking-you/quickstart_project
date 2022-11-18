## dao_convertor example

### tag的支持与使用
由于dao层的代码基于 [gorm](https://github.com/go-gorm/gorm) 框架生成对应的CURD代码，所以 `dao` tag的作用就是控制CRUD的gorm代码生成。
语法形式如：`dao:"c:操作1;r:操作2;u:操作3;d:操作4"`，这里的c、r、u、d控制的就是gorm代码对于该字段的CRUD生成。
#### Delete的Where条件
gorm 中对删除的操作，如果传入的是一个struct，它只会看这个struct的主键然后删除对应的一行，如果要根据条件删除，还是要是添加一些where条件来进行删除。
代码如下：
```go
db.Where("age = 20").Delete(&user)
```
你可以通过在User的Age属性的tag中加入 `dao:"d:(age=20)"`来实现上述代码的生成。你可以在config中配置where条件的连接方式，是以or还是and。

注意：目前where操作只支持delete，因为其他操作可以直接Model进行替代。
#### 其他操作的Omit
在gorm中，由于传入数据后，默认是以他的所有字段进行CRUD，但是业务上有时候由于struct的复杂性，可能有其他struct有嵌套的情况，所以需要忽略某些字段来进行CRUD。
对应的gorm代码：
```go
db.Omit("Name", "Age", "CreatedAt").Create(&user)
```
你可以在User的Name和Age字段上加入 `dao:"c:omit"`来生成上述代码。
此操作支持所有CRUD。
### config配置
> 具体的配置信息请查看源代码注释：[code](../config.go)

### 简单示例

> 原始配置
```go
package main

import (
	"github.com/ACking-you/quickstart_project/dao_convertor"
)

type User struct {
	Id   int 
	Name string `dao:"c:omit;r:omit;u:omit;d:(name <> ?)"`
}

func autoDAO() {
	config := dao_convertor.DefaultConfig("root", "123", "127.0.0.1", 3306, "test").
		//是否产生CRUD代码，默认为true
		EnableCreate(true).
		EnableQuery(true).
		EnableUpdate(true).
		EnableDelete(true)
	convert := dao_convertor.NewStruct2DAO(config)

	err := convert.AutoMigrate(&User{}).Run()
	if err != nil {
		panic(err)
	}
}


func main() {
	autoDAO()
}

```

> 代码生成
```go
package dao

type UserDAO struct {
}

var userDAO = UserDAO{}

func NewUserDAO() *UserDAO {
	return &userDAO
}

func (UserDAO) CreateUser(info *test.User) error {
	return db.Omit("name", "password").Create(info).Error
}

func (UserDAO) CreateAllUser(infos *[]test.User) error {
	return db.Omit("name", "password").Create(infos).Error
}

func (UserDAO) QueryUser(cond, result *test.User) error {
	return db.Model(cond).Omit("name", "password").First(result).Error
}

func (UserDAO) QueryAllUser(cond, results *[]test.User) error {
	return db.Model(cond).Omit("name", "password").Find(results).Error
}

func (UserDAO) UpdateUser(old, new *test.User) error {
	return db.Model(old).Omit("name", "password").Updates(*new).Error
}

func (UserDAO) DeleteUser(cond *test.User) error {
	return db.Where(" name <> ? AND password = ?", cond.Name, cond.Password).Delete(cond).Error
}

```