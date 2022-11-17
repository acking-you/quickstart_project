## controller_convertor example

```go
package main

import (
	controller_convertor "github.com/ACking-you/quickstart_project/controller_convertor"
	"github.com/ACking-you/quickstart_project/controller_convertor/example/models"
	"github.com/ACking-you/quickstart_project/service_convertor"
)

func CreateService() {
	config := service_convertor.DefaultConfig().SavePath("./example/service").EnableDebug(true)
	err := service_convertor.NewStruct2Service(config).AutoMigrate(&models.Student{}, &models.Score{}, &models.Teacher{}).Run()
	if err != nil {
		panic(err)
	}
}

func testController() {
	config := controller_convertor.DefaultCConfig().
		SavePath("./example/controller").
		EnableDebug(true).
		EnableQuery(false)

	err := controller_convertor.NewStruct2Controller(config).
		AutoMigrate(&models.Score{}, &models.Student{}, &models.Teacher{}).
		Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	testController()
}
```