package app

import (
	"github.com/devfeel/dotweb"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/track/blogserver/pkg/config"
	"github.com/track/blogserver/pkg/controllers"
	"github.com/track/blogserver/pkg/persistence"
	"github.com/track/blogserver/pkg/routers"
	"os"
)

//博客应用服务器
type App struct {
	DB     *gorm.DB
	RPool  *redis.Pool
	Conf   *config.BlogConfig
	Server *dotweb.DotWeb
}

func NewApp() *App {
	return &App{}
}

//启动服务器
func (app *App) Launch() error {
	app.Conf = config.Config()
	app.initDB()
	app.initRedis()
	app.initServer()
	app.initRouter()
	app.initImgServer()
	app.initApiDocs()
	return app.Server.StartServer(app.Conf.ServerPort)
}

//关闭操作
func (app *App) Destory() {
	if app.DB != nil {
		app.DB.Close()
	}
	if app.RPool != nil {
		app.RPool.Close()
	}
	if app.Server != nil {
		app.Server.Close()
	}
}

//根据配置文件初始化数据库
func (app *App) initDB() {
	app.DB = persistence.GetOrm()
}

//根据配置文件初始化Redis
func (app *App) initRedis() {
	app.RPool = persistence.GetRedisPool()
}

//根据配置初始化服务器
func (app *App) initServer() {
	app.Server = dotweb.New()

	//配置Log
	app.Server.SetEnabledLog(app.Conf.LogEnable)
	app.Server.SetLogPath(app.Conf.LogPath)

	//配置自定义error
	app.initError()

	//配置环境模式
	if app.Conf.EnvProd {
		app.Server.SetProductionMode()
	} else {
		app.Server.SetDevelopmentMode()
	}
	//开启Gzip压缩
	app.Server.HttpServer.SetEnabledGzip(false)
}

//初始化路由配置
func (app *App) initRouter() {
	r := routers.NewApiRouter(app.Server.HttpServer)
	// api/v1/xx api
	r.V1()
	// api/admin/xx api
	r.Admin()
}

//配置自定义error
func (app *App) initError() {
	ec := controllers.NewErrorController()
	app.Server.SetNotFoundHandle(ec.NotFound)
	app.Server.SetExceptionHandle(ec.Internal)
	app.Server.SetMethodNotAllowedHandle(ec.MethodNotAllowed)
}

//配置图片文件服务器
func (app *App) initImgServer() {
	if _, err := os.Stat(app.Conf.ImgPath); err != nil {
		if err = os.MkdirAll(app.Conf.ImgPath, os.ModePerm); err != nil {
			panic("Create ImagePath Error!")
		}
	}
	//开启图片文件服务器访问,否则无法访问
	app.Server.HttpServer.ServerFile("/image/*filepath", app.Conf.ImgPath)
}
// 初始化文档服务器配置
func (app *App) initApiDocs() {
	path := os.Getenv("APP_CONFIG_PATH")
	if len(path) <= 0 {
		path = "docs/swaggerui"
	}else{
		path = "home/blogserver/api_docs"
	}
	app.Server.HttpServer.ServerFile("/docs/*filepath", path)
}