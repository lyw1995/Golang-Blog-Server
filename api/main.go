// Golang-Blog-Server - 博客API文档
//
// API Docs for Golang-Blog-Server v1
//
// 	 Terms Of Service:  N/A
//     Schemes: [http]
//     Version: 1.0.0
//     License: N/A
//     Contact: Track <24245@163.com> http://blog.yinguiw.com
//     Host:
//	   basePath: /api
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer: []
//
//     SecurityDefinitions:
//     bearer:
//          type: token
//          name: Authorization
//          in: header
//
// swagger:meta
package api


// 输出文档json到控制台
//go:generate swagger generate spec
// 输出文档json到指定文件
//go:generate swagger generate spec -o  ../docs/swaggerui/swagger.json --scan-models
// 打开本地api文档服务
//go:generate swagger serve -F=swagger ../docs/swaggerui/swagger.json