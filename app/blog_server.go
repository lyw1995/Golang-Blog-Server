package app

import (
	"blogserver/config"
	"blogserver/controllers"
	"blogserver/persistence"
	"blogserver/routers"
	"bufio"
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/jinzhu/gorm"
	"os"
	"github.com/garyburd/redigo/redis"
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
	app.bless()
	return app.Server.StartServer(app.Conf.ServerPort)
}

//打印 佛陀logo
func (app *App) bless() {
	if !app.Conf.EnvProd {
		return
	}
	file, err := os.Open("gless.txt")
	defer file.Close()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
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
	app.Server.HttpServer.SetEnabledGzip(true)
}

//初始化路由配置
func (app *App) initRouter() {
	r := routers.NewApiRouter(app.Server.HttpServer)

	// /api/v1/ 版本接口
	r.V1()

	// 后台接口
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
