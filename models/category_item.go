package models

import (
	"time"
)

// 分类条目: 条目ID  分类ID 条目名称 创建时间  修改时间
// 如果是归档,Name就是YYYY年-MM月
type CategoryItem struct {
	ID         uint       `gorm:"primary_key" json:"cid"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `sql:"index" json:"-"`
	Name       string     `gorm:"not null;size:12" json:"cname"`
	CreateTime time.Time  `gorm:"not null" json:"create_time"`
	CategoryID uint       `gorm:"index;not null"json:"-"`
	Article    []Article  `json:"-"`
	ItemSize   int        `gorm:"default:0" json:"size"`
	UserID     uint       `gorm:"index;not null" json:"-"`
}
