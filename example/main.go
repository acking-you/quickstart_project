package main

import (
	"github.com/ACking-you/quickstart_project"
	"github.com/ACking-you/quickstart_project/util"
)

func genChainStyleConfig() {
	_ = util.GenerateMethodByChainStyle(&quickstart.Config{}, "config.go", "config.go")
}

func testQuickStart() {

	config := quickstart.DefaultConfig("项目名称", "root", "123", "127.0.0.1", 3306, "数据库名称").
		//打印出生成结果
		EnableDebug(true).
		//改变基本路径（默认为项目根目录）
		BasePath("./example")
	err := quickstart.Run(config)
	if err != nil {
		panic(err)
	}
}

func main() {
	testQuickStart()
}
