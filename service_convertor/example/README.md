## service_convertor example

```go
package main

import (
	"github.com/ACking-you/quickstart_project/service_convertor"
	"github.com/ACking-you/quickstart_project/service_convertor/example/models"
	"github.com/ACking-you/quickstart_project/util"
)

func genChainStyleConfig() {
	_ = util.GenerateMethodByChainStyle(&service_convertor.Config{}, "config.go", "config.go")
}

func serviceTest() {
	config := service_convertor.DefaultConfig().
		EnableDebug(true).
		SavePath("./example/service/").
		DefaultMethodName("Login").
		EnableTOFileSingle(true).
		EnableVOFileSingle(true)
	err := service_convertor.NewStruct2Service(config).
		AutoMigrate(&models.Score{}, &models.Course{}, &models.TeachCourse{}, &models.Teacher{}, &models.Student{}).
		Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	serviceTest()
}
```