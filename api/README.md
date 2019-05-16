# Golang-Blog-Server

### go-swagger api 文档
  >> 注意: 只是简单的定义了api便于展示(缺少Req,Resp model等), 而且api的定义和划分都不太理想
```bash
1. mac安装swagger
	brew tap go-swagger/go-swagger
	brew install go-swagger

2. 学习例子 
    https://goswagger.io/generate/spec/params.html
    https://studygolang.com/articles/12354?fr=sidebar

// 在api 目录下执行
3.  生成api配置文件 
    // model暂时不定义,需要可以自己独立出req,resp model
    swagger generate spec -o  ../docs/swaggerui/swagger.json --scan-models

4. 运行ui服务器
    swagger serve -F=swagger ../docs/swaggerui/swagger.json
```

