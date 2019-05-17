# Golang-Blog-Server
Go语言编写的简易版博客服务端
+ 开启go mod` export GO111MODULE=on`
+ god mod 代理设置`export GOPROXY=https://athens.azurefd.net`
+ `make api` 重新生成文档
+ `make build` 打包docker image镜像
    >> ### terminal run
    `go run cmd/main.go` 启动服务(redis,mysql关键依赖需要安装,查看配置文件)
 
    >> ### docker run (安装docker,docker-compose具体百度)
    - `make up` docker-compose启动服务
    - `make down` docker-compose停止并删除服务
+ `curl http://localhost:8888/api/v1/users` 获取初始化插入用户数据
+ ` http://localhost:8888/docs/index.html` 查看api文档
+ ` docker方式,浏览器访问 http://localhost:24245/` 博客客户端
+ ` docker方式,浏览器访问 http://localhost:24246/` 博客后台管理

![apidoc](/screenshot/apidoc.png)
![blogadmin](https://github.com/lyw1995/Angular5-Blog-Admin/raw/master/snapshot/blog_admin.png)
![blogfront](https://github.com/lyw1995/Angular5-Blog-Front/raw/master/snapshot/blog_front.png)
## 项目依赖
* [dotweb](https://github.com/devfeel/dotweb)
* jwt
* redis
* goquery
* gorm
* toml

## 相关项目
* [BlogFront博客客户端](https://github.com/lyw1995/Angular5-Blog-Front)
* [Admin后台管理系统](https://github.com/lyw1995/Angular5-Blog-Admin)

