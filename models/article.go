package models

import (
	"time"
)

// 文章: 文章ID 用户ID 分类ID 文章标题 封面 文章内容 创建时间 修改时间 阅读数 评论数
type Article struct {
	ID             uint       `gorm:"primary_key" json:"aid"`
	CreatedAt      time.Time  `json:"-"`
	UpdatedAt      time.Time  `json:"-"`
	DeletedAt      *time.Time `sql:"index" json:"-"`
	UserID         uint       `gorm:"index;not null" json:"-"`
	CategoryItemID uint       `gorm:"index;not null" json:"cid"`
	Title          string     `gorm:"type:varchar(100);not null" json:"title"`
	Content        string     `gorm:"type:text;not null;" json:"content"`
	Cover          string     `gorm:"not null" json:"cover"`
	CreateTime     time.Time  `gorm:"not null" json:"create_time"`
	Views          int        `gorm:"default:0" json:"views"`
	Origin         int        `gorm:"not null" json:"origin"` //是否原创 1原创 0转载
	State          int        `gorm:"default:0" json:"-"`     //0正常发布 2并未发布(草稿箱)
	//	Comments int `gorm:"default:0"` //暂时先不做评论
}

func (ar Article) IsValid() bool {
	return ar.ID > 0
}
