package test

import (
	"github.com/ACking-you/quickstart_project/service_convertor"
	"testing"
)

//注意!!dao层的key要以`:`和具体操作隔开。以`;`分割对不同key的操作
//而service层则不需要`:`，因为对key值对应的操作只有一种，就是填充对应的vo标签
type Users struct {
	Id int `gorm:"id" dao:"c:omit;r:omit;d:(id>?)" service:"binding(required,email);"`
}

func TestService(t *testing.T) {
	config := service_convertor.DefaultConfig().EnableDebug(true)
	err := service_convertor.NewStruct2Service(config).AutoMigrate(&Users{}).Run()
	if err != nil {
		panic(err)
	}
}