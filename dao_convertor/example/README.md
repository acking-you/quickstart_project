## dao_convertor example

```go
package main

import (
	"github.com/ACking-you/quickstart_project/dao_convertor"
	"github.com/ACking-you/quickstart_project/dao_convertor/example/models"
	"github.com/ACking-you/quickstart_project/util"
)

func genChainStyleConfig() {
	_ = util.GenerateMethodByChainStyle(&dao_convertor.Config{}, "config.go", "config.go")
}

func testDAO() {
	config := dao_convertor.DefaultConfig("root", "123", "127.0.0.1", 3306, "my_chat")
	config.SavePath("./example/dao/test.go")

	convert := dao_convertor.NewStruct2DAO(config)
	err := convert.AutoMigrate(&models.Score{}, &models.Teacher{}, &models.Course{}, &models.TeachCourse{}, &models.Student{}).Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	testDAO()
}

```