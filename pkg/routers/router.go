package routers

import (
	"github.com/devfeel/dotweb"
	"github.com/track/blogserver/pkg/config"
	"github.com/jinzhu/gorm"
	"github.com/track/blogserver/pkg/controllers"
	"github.com/track/blogserver/pkg/middleware"
)

var DBConn *gorm.DB

type Router struct {
	server *dotweb.HttpServer
	group  dotweb.Group
}

//路由配置
func NewApiRouter(server *dotweb.HttpServer) *Router {
	router := &Router{server: server, group: server.Group("/api")}
	if config.Config().EnvProd {
		router.group.Use(middleware.NewCROSMiddleware())
	}
	return router
}

func (r *Router) V1() {
	v1 := r.group.Group("/v1")
	//正式环境开启 签名验证中间件
	if config.Config().EnvProd {
		v1.Use(&middleware.ApiSignMiddleware{})
	}
	v1User(v1.Group("/users"))
	v1Upload(v1.Group("/upload"))
}

func v1User(users dotweb.Group) {
	uc := controllers.NewUserController()
	users.OPTIONS("",uc.Options)
	users.OPTIONS("/:uid",uc.Options)
	//获取客户端初始化参数
	users.GET("", uc.InitBlog)
	//获取单个用户信息
	users.GET("/:uid", uc.GetUserByUid)

	v1UserLink(users.Group("/:uid/links"))
	v1UserCategory(users.Group("/:uid/categorys"))
	v1UserArticle(users.Group("/:uid/articles"))
	v1UserOther(users.Group("/:uid/other"))
}

//资源上传
func v1Upload(upload dotweb.Group) {
	uc := controllers.NewUploadController()
	upload.POST("", uc.UploadImage)
}
//热门最新文章
func v1UserOther(other dotweb.Group) {
	ac := controllers.NewArticleController()
	other.OPTIONS("",ac.Options)
	other.GET("", ac.GetHotAndNewArticles)
}
// 友链管理
func v1UserLink(links dotweb.Group) {
	uc := controllers.NewLinkController()
	links.OPTIONS("",uc.Options)
	links.OPTIONS("/:lid",uc.Options)
	//获取用户全部友链
	links.GET("", uc.GetLinks)
	//获取用户单个友链
	links.GET("/:lid", uc.GetLinkById)
}
func v1UserArticle(article dotweb.Group) {
	ac := controllers.NewArticleController()
	article.OPTIONS("",ac.Options)
	article.OPTIONS("/:aid",ac.Options)
	//获取用户全部文章
	article.GET("", ac.GetUserArticles)
	//获取某篇文章
	article.GET("/:aid", ac.GetArticleById)
}
func v1UserCategory(category dotweb.Group) {
	cc := controllers.NewCategoryController()
	category.OPTIONS("",cc.Options)
	//获取全部分类
	category.GET("", cc.GetCategorys)
	v1CategoryArticle(category.Group("/:cid/articles"))
}

func v1CategoryArticle(article dotweb.Group) {
	ac := controllers.NewArticleController()
	article.OPTIONS("",ac.Options)
	article.OPTIONS("/:aid",ac.Options)
	//获取分类全部文章
	article.GET("", ac.GetCategoryArticles)
	//获取分类下某篇文章
	article.GET("/:aid", ac.GetArticleById)
}
func (r *Router) Admin() {
	admin := r.group.Group("/admin")
	//用户会话管理
	adminSession(admin.Group("/sessions"))
	//后台其他扩展管理
	adminExtend(admin.Group("/extends"))
	// users使用jwt
	adminUser(admin.Group("/users", middleware.NewJwtMiddleware(r.server.DotApp)))
}
func adminExtend(ext dotweb.Group) {
	ec := controllers.NewExtendController()
	ext.OPTIONS("", ec.Options)
	//根据请求参数处理
	ext.GET("", ec.ExtByQueryParams)
	//处理文章采集
	ext.POST("", ec.ArticleCollection)
}
func adminUser(users dotweb.Group) {
	uc := controllers.NewUserController()
	//探测
	users.OPTIONS("/:uid", uc.Options)
	users.OPTIONS("", uc.Options)

	//获取全部用户信息
	users.GET("", uc.GetUsers)
	//获取单个用户信息
	users.GET("/:uid", uc.GetUserByUid)
	//创建用户
	users.POST("", uc.CreateUser)
	//删除全部用户
	users.DELETE("", uc.DelAllUser)
	//删除单个用户
	users.DELETE("/:uid", uc.DelByUid)
	//修改用户资料
	users.PUT("/:uid", uc.UpdateUserInfoByUid)

	adminUserLink(users.Group("/:uid/links"))
	adminUserCategory(users.Group("/:uid/categorys"))
	adminUserArticles(users.Group("/:uid/articles"))
}
func adminUserArticles(articles dotweb.Group) {
	ac := controllers.NewArticleController()
	articles.OPTIONS("", ac.Options)
	articles.OPTIONS("/:aid", ac.Options)

	articles.GET("", ac.GetUserArticles)
	articles.GET("/:aid", ac.GetArticleWithState)
	articles.PUT("/:aid", ac.UpdateArticleWithState)
}
func adminSession(session dotweb.Group) {
	uc := controllers.NewUserController()
	//session.OPTIONS("", uc.Options)
	//用户登录
	session.POST("", uc.UserLogin)

	//TODO 新版本dotweb , panic: 'upload' in new path '/api/admin/sessions/upload'
	// conflicts with existing wildcard ':uid' in existing prefix '/api/admin/sessions/:uid'
	// 先注释 用户退出api, 也没实现
	//用户退出
	//session.DELETE("/:uid", uc.UserLogout)

	//管理员登录后上传
	adminSessionUpload(session.Group("/upload"))
}
func adminSessionUpload(upload dotweb.Group) {
	uc := controllers.NewUploadController()
	upload.OPTIONS("", uc.Options)
	upload.POST("", uc.UploadImage)
}
func adminUserLink(links dotweb.Group) {
	lc := controllers.NewLinkController()

	links.OPTIONS("", lc.Options)
	links.OPTIONS("/:lid", lc.Options)

	//获取用户所有友链
	links.GET("", lc.GetLinks)
	//用户创建友链
	links.POST("", lc.CreateFriendlyLink)
	//用户删除全部友链
	links.DELETE("", lc.DelLinks)
	//用户删除某个友链
	links.DELETE("/:lid", lc.DelLinkById)
	//用户修改某个友链
	links.PUT("/:lid", lc.UpdateLinkById)

}
func adminUserCategory(category dotweb.Group) {
	cc := controllers.NewCategoryController()

	category.OPTIONS("", cc.Options)
	category.OPTIONS("/:cid", cc.Options)
	//获取个人分类
	category.GET("", cc.GetPersonalCategorys)
	//创建分类
	category.POST("", cc.CreateCategory)
	//修改某个分类
	category.PUT("/:cid", cc.UpdateCategoryById)
	//删除某个分类
	category.DELETE("/:cid", cc.DelCategoryById)
	//删除全部分类
	category.DELETE("", cc.DelCategorys)

	adminCategoryArticle(category.Group("/:cid/articles"))
}
func adminCategoryArticle(article dotweb.Group) {
	ac := controllers.NewArticleController()
	article.OPTIONS("", ac.Options)
	article.OPTIONS("/:aid", ac.Options)

	//分类下创建文章
	article.POST("", ac.CreateArticleByCid)
	//分类下修改文章
	article.PUT("/:aid", ac.UpdateArticleByCid)
	//分类删除全部文章
	article.DELETE("", ac.DelArticlesByCid)
	//分类删除单篇文章
	article.DELETE("/:aid", ac.DelArticleByCid)
}
