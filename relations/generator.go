package relations

import (
	"blogserver/models"
	"blogserver/persistence"
	"github.com/jinzhu/gorm"
	"time"
)

/*
	generator.go 根据结构体(Models) 统一创建数据库关系
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
