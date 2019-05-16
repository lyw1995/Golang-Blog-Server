package models

type UserInfo struct {
	BaseModel
	UserAvator string `gorm:"not null;default:'http://dwz.cn/WtaRcB72'"`
	UserDesc   string `gorm:"not null;"`
	UserEmail  string `gorm:"not null;"`
	UserAddr   string `gorm:"not null;"`
	NickName   string `gorm:"not null;size:12"`
}
