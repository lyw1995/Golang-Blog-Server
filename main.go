package main

import (
	"blogserver/app"
	"blogserver/relations"
)

func main() {
	//判断表是否存在 存在就自动迁移模式
	relations.InitRelations()
	//启动http服务
	app := app.NewApp()
	defer app.Destory()
	app.Launch()
}
