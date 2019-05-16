package repositories

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
)

type ArticleRepository struct {
	*BaseRepository
}

func (ar *ArticleRepository) GetArticlesWithState(uid, page, state int, sort models.SortValue) ([]models.Article, bool) {
	articles := make([]models.Article, 0)
	return articles, ar.Exec(func(db *gorm.DB) *gorm.DB {
		if sort.Key == "create_time" {
			return db.Where("user_id = ? and state = ? ", uid, state).Order(fmt.Sprintf("%s %s", sort.Key, sort.Value))
		} else {
			return db.Where("user_id = ? and state = ? ", uid, state)
		}
	}, &articles, page*common.LenLimit, common.LenLimit)
}
func (ar *ArticleRepository) GetCategoryItemArticlesPage(uid, cid, page int, sort models.SortValue) ([]models.Article, bool) {
	articles := make([]models.Article, 0)
	return articles, ar.Exec(func(db *gorm.DB) *gorm.DB {
		if sort.Key == "create_time" {
			return db.Where("user_id = ? and category_item_id = ? and state = 0", uid, cid).Order(fmt.Sprintf("%s %s", sort.Key, sort.Value))
		} else {
			return db.Where("user_id = ? and category_item_id = ? and state = 0", uid, cid)
		}
	}, &articles, page*common.LenLimit, common.LenLimit)
}
func (ar *ArticleRepository) GetArchizeArticlesPage(uid, page int, cname string, sort models.SortValue) ([]models.Article, bool) {
	articles := make([]models.Article, 0)
	return articles, ar.Exec(func(db *gorm.DB) *gorm.DB {
		if sort.Key == "create_time" {
			return db.Where("user_id = ? and date_format(create_time ,'%Y年%m月' )  = ? and state = 0", uid, cname).Order(fmt.Sprintf("%s %s", sort.Key, sort.Value))
		} else {
			return db.Where("user_id = ? and date_format(create_time ,'%Y年%m月' )  = ? and state = 0", uid, cname)
		}
	}, &articles, page*common.LenLimit, common.LenLimit)
}

//根据用户的文章状态统计
func (ar *ArticleRepository) CountByUserIdWithState(uid, state int) int {
	_, count := ar.Count(&models.Article{}, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ? and state = ?", uid, state)
	})
	return count
}
func (ar *ArticleRepository) GetCountByUserId(uid int) int {
	return ar.CountByUserIdWithState(uid, 0)
}
func (ar *ArticleRepository) GetDraftCountByUserId(uid int) int {
	return ar.CountByUserIdWithState(uid, 2)
}

//根据分类id统计已发布
func (ar *ArticleRepository) GetCountByCid(uid ,cid int) int {
	_, count := ar.Count(&models.Article{}, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ? and category_item_id = ? and state = 0",uid, cid)
	})
	return count
}

//根据分类名称统计已发布
func (ar *ArticleRepository) GetCountByCname(uid int,cname string) int {
	_, count := ar.Count(&models.Article{}, func(db *gorm.DB) *gorm.DB {
		return db.Where("date_format(create_time ,'%Y年%m月' ) = ? and state = 0 and user_id = ?", cname,uid)
	})
	return count
}

// 根据状态查文章
func (ar *ArticleRepository) GetByUidAndAidWithState(uid, aid, state int) (*models.Article, bool) {
	article := &models.Article{}
	return article, ar.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ? and id = ?  and state = ?", uid, aid, state)
	}, &article)
}
func (ar *ArticleRepository) GetByUidAndCidAndAidWithState(uid, cid, aid, state int) (*models.Article, bool) {
	article := &models.Article{}
	return article, ar.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("category_item_id = ? and id = ? and user_id = ? and state = ?", cid, aid, uid, state)
	}, &article)
}

//不依赖状态
func (ar *ArticleRepository) GetByUidAndAid(uid, aid int) (*models.Article, bool) {
	article := &models.Article{}
	return article, ar.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where(" user_id = ? and id = ? ", uid, aid)
	}, &article)
}
func (ar *ArticleRepository) DelAllByUidAndCid(uid, cid int) error {
	return ar.Del(func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ? and category_item_id =?", uid, cid)
	}, &models.Article{})
}
func (ar *ArticleRepository) DelFor(article *models.Article) error {
	return ar.Del(func(db *gorm.DB) *gorm.DB {
		return db
	}, article)
}
