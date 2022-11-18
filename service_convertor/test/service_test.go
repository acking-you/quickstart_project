package test

import (
	"github.com/ACking-you/quickstart_project/service_convertor"
	"testing"
)

//注意!!dao层的key要以`:`和具体操作隔开。以`;`分割对不同key的操作
//而service层则不需要`:`，因为对key值对应的操作只有一种，就是填充对应的vo标签
type Users struct {
	Id int `gorm:"id" dao:"c:omit;r:omit;d:(id>?)" service:"binding(required,email);form(test)"`
}

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
func TestService(t *testing.T) {
	autoService()
}
