package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

//分类: 分类ID  分类名称 创建时间  修改时间
type Category struct {
	gorm.Model   `json:"-"`
	Name         string         `gorm:"not null;size:12" json:"label"`
	CreateTime   time.Time      `gorm:"not null" json:"-"`
	CategoryItem []CategoryItem `json:"categorys"`
}
