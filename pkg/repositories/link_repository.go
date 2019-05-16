package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/track/blogserver/pkg/models"
)

//友链管理,对应数据中表: friendly_links
type LinkRepository struct {
	*BaseRepository
}

//根据用户id,友链Id获取指定友链
func (lr *LinkRepository) GetById(uid, id int) (*models.FriendlyLink, bool) {
	link := &models.FriendlyLink{}
	return link, lr.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ? and user_id = ?", id, uid)
	}, link)

}

//获取指定用户的所有友链
func (lr *LinkRepository) GetAllByUserId(uid int) ([]models.FriendlyLink, bool) {
	links := make([]models.FriendlyLink, 0)
	return links, lr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", uid)
	}, &links)
}

//删除根据(查找传入的)指定的友链
func (lr *LinkRepository) DelByInstance(link *models.FriendlyLink) error {
	return lr.Del(func(db *gorm.DB) *gorm.DB {
		return db
	}, link)
}

//喊出某个用户的所有友链
func (lr *LinkRepository) DelAllByUid(uid int) error {
	link := &models.FriendlyLink{}
	return lr.Del(func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", uid)
	}, link)
}
