package main

import (
	"github.com/ACking-you/quickstart_project"
	"github.com/ACking-you/quickstart_project/util"
)

func genChainStyleConfig() {
	_ = util.GenerateMethodByChainStyle(&quickstart.Config{}, "config.go", "config.go")
}

func testQuickStart() {

	err := quickstart.Run(quickstart.DefaultConfig("github.com/ACking-you/quickstart_project/example", "root", "123", "127.0.0.1", 3306, "my_chat").
		EnableDebug(true).BasePath("./example"))
	if err != nil {
		panic(err)
	}
}

func main() {
	testQuickStart()
}
