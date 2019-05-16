package relations

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/track/blogserver/pkg/config"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/persistence"
	"github.com/track/blogserver/pkg/utils"
	"time"
)

/*
	generator.go 根据结构体(Models) 统一创建数据库关系

        1. 初始化分类
		2. 初始化默认用户

*/
func InitRelations() {
	db := persistence.GetOrm()
	// 判断存不存在表, 不存在就新建, 否则就是自动迁移(其他修改)
	if !db.HasTable("users") {
		db.CreateTable(&models.User{}, &models.UserInfo{}, &models.Category{}, &models.CategoryItem{},
			&models.Article{}, models.FriendlyLink{})
	} else {
		db.AutoMigrate(&models.User{}, &models.UserInfo{}, &models.Category{}, &models.CategoryItem{},
			&models.Article{}, models.FriendlyLink{})
	}
	initCategory(db)
	initDefaultUser(db)
}
// 插入默认用户 数据库记录
func initDefaultUser(db *gorm.DB) {
	defaultUser := config.Config().DefaultClientUser
	if found := db.Find(&models.User{},"user_name = ?",defaultUser).RecordNotFound(); found {
		// 使用UserService插入
		// 使用UserRepo插入
		// 直接db插入
		user := models.User{
			UserName:     defaultUser,
			UserPassWord: utils.Sha1(defaultUser+"123"),
			LoginTime:    time.Now(),
			LoginIP:      sql.NullString{String: "127.0.0.1", Valid: true},
			UserInfo: models.UserInfo{
				UserAvator: "https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=866585511,3203326197&fm=27&gp=0.jpg",
				UserDesc:   "无法描述了.",
				UserAddr:   "我迷路了.",
				UserEmail:  "24245@163.com",
				NickName:   defaultUser,
			},
		}
		db.Create(&user)
	}
}

//系统默认插入主分类("个人分类","归档")
func initCategory(db *gorm.DB) {
	if found := db.Find(&models.Category{}).RecordNotFound(); found {
		db.Create(&models.Category{
			Name:       "个人分类",
			CreateTime: time.Now(),
		})
		db.Create(&models.Category{
			Name:       "归档",
			CreateTime: time.Now(),
		})
	}
}
