# Golang-Blog-Server

![blogadmin](https://github.com/lyw1995/Angular5-Blog-Admin/raw/master/snapshot/blog_admin.png)
![blogfront](https://github.com/lyw1995/Angular5-Blog-Front/raw/master/snapshot/blog_front.png)

Go语言编写的简易版博客服务端

+ 项目创建基于 `go version go1.10.3 darwin/amd64`
+ 测试运行 `go run main.go` 打开`http://localhost:8888`
+ Api配置在`./routers`目录下,自行测试
+ 正式环境部署使用docker-compose,[yml文件参考](http://blog.yinguiw.com/article/detail/1)

## 项目依赖
* [dotweb](https://github.com/devfeel/dotweb)
* jwt
* redis
* goquery
* gorm
* govendor
* toml

## 相关项目
* [BlogFront博客客户端](https://github.com/lyw1995/Angular5-Blog-Front)
* [Admin后台管理系统](https://github.com/lyw1995/Angular5-Blog-Admin)

