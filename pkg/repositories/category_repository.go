package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
	"time"
)

type CateGoryRepository struct {
	*BaseRepository
}

//根据分类名获取分类 (这是程序初始化生成固定的两个分类('个人分类',和 '归档'))
func (cr *CateGoryRepository) GetCategoryByName(name string) (*models.Category, bool) {
	category := &models.Category{}
	return category, cr.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}, category)
}
func (cr *CateGoryRepository) GetCategoryById(id uint) (*models.Category, bool) {
	category := &models.Category{}
	return category, cr.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}, category)
}

//获取所有分类以及子分类条目
func (cr *CateGoryRepository) GetCategorysByUid(uid int) ([]models.Category, bool) {
	categorys := make([]models.Category, 0)
	return categorys, cr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Preload("CategoryItem", "user_id =?", uid)
	}, &categorys)
}

//根据分类名获取所有子条目
func (cr *CateGoryRepository) GetCategoryItemsByCid(uid, cid uint) ([]models.CategoryItem, bool) {
	categoryItems := make([]models.CategoryItem, 0)
	return categoryItems, cr.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Where("category_id = ? and user_id = ?", cid, uid)
	}, &categoryItems)
}

//根据子分类条目名获取子分类条目
func (cr *CateGoryRepository) GetCategoryItemByNameWithUid(name string, uid int) (*models.CategoryItem, bool) {
	categoryItem := &models.CategoryItem{}
	return categoryItem, cr.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ? and user_id = ?", name, uid)
	}, categoryItem)
}

//根据子条目分类Id获取子分类条目
func (cr *CateGoryRepository) GetCategoryItemById(cid int) (*models.CategoryItem, bool) {
	categoryItem := &models.CategoryItem{}
	return categoryItem, cr.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", cid)
	}, categoryItem)
}

//删除所有子条目
func (cr *CateGoryRepository) DelCategoryItems(uid int) error {
	return cr.Del(func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?",uid)
	}, &models.CategoryItem{})
}

//删除指定子条目
func (cr *CateGoryRepository) DelCategoryItemById(cid int) error {
	return cr.Del(func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", cid)
	}, &models.CategoryItem{})
}

//根据子条目id修改子条目 item_size
func (cr *CateGoryRepository) UpdateItemSizeByCid(cid uint, size int) error {
	return cr.Update(&models.CategoryItem{ID: cid},
		map[string]interface{}{"item_size": size})
}

//根据归档条目名称获取指定归档或者创建某个归档条目
func (cr *CateGoryRepository) GetFirstOnCreateArchive(uid uint, cname string) (*models.CategoryItem, error) {
	if category, found := cr.GetCategoryByName("归档"); found {
		categoryItem := models.CategoryItem{
			Name:       cname,
			CategoryID: category.ID,
			CreateTime: time.Now(),
			UserID:     uid,
		}
		return &categoryItem, cr.db().Where("name = ? and user_id = ?", cname, uid).FirstOrCreate(&categoryItem).Error
	} else {
		return nil, common.ErrCategoryNoExist
	}
}

//根据子条目 修改子条目 item_size
func (cr *CateGoryRepository) UpdateItemSizeBy(citem *models.CategoryItem) error {
	return cr.Update(citem, map[string]interface{}{"item_size": citem.ItemSize})
}
