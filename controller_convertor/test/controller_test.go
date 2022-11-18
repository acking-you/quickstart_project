package test

import (
	"github.com/ACking-you/quickstart_project/controller_convertor"
	"testing"
)

type User struct {
	Id       int
	Username string
	Password string
}

func autoController() {
	config := controller_convertor.DefaultCConfig().
		//支持配置是否生成CRUD
		EnableQuery(false).
		//是否生成gin的棋手例子
		EnableGinExample(true).
		//是否生成封装好的Response类
		EnableResponse(true).
		//是否生成VO和前端数据的绑定代码
		EnableVOBind(true)

	err := controller_convertor.NewStruct2Controller(config).
		AutoMigrate(&User{}).
		Run()
	if err != nil {
		panic(err)
	}
}

func TestController(t *testing.T) {
	autoController()
}
