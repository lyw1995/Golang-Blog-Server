package models

import "time"

// 友情链接: ID 用户ID 头像 链接名 链接URL
type FriendlyLink struct {
	ID        uint      `gorm:"primary_key" json:"link_id"` // 自增
	CreatedAt time.Time `json:"create_time"`
	UpdatedAt time.Time `json:"-"`
	UserID    uint      `gorm:"index;not null" json:"-"`
	Avator    string    `gorm:"not null" json:"link_icon"`
	LinkName  string    `gorm:"not null;size:12" json:"link_name"`
	LinkUrl   string    `gorm:"not null" json:"link_url"`
}
