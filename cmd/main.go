package main

import (
	"github.com/track/blogserver/pkg/app"
	"github.com/track/blogserver/pkg/relations"
	"log"
)

func main() {

	// 判断表是否存在 存在就自动迁移模式
	// generator.go 根据结构体(Models) 统一创建数据库关系
	relations.InitRelations()

	// 初始化app
	app := app.NewApp()

	defer app.Destory()

	// 启动
	log.Fatal(app.Launch())

}
