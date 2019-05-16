package models

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	BaseModel
	//OwnerID     uuid.UUID    `gorm:"type:uuid" json:"owner_id"` //不用UUID
	UserName     string `gorm:"not null;size:12;unique"` //用户名唯一
	UserPassWord string `gorm:"not null;size:40"`        //sha1
	LoginIP      sql.NullString
	LoginTime    time.Time
	IsActive     int      `gorm:"not null;default:0"` //0可用 1禁用
	UserInfo     UserInfo `gorm:"ForeignKey:UserInfoID;AssociationForeignKey:ID"`
	UserInfoID   uint     `gorm:"index;not null"`
	Article      []Article
	FriendlyLink []FriendlyLink
	CategoryItem []CategoryItem
}

func (user *User) IsValidAccout() bool {
	return user.IsActive > 0
}
func (user *User) IsAdmin() bool {
	return user.UserName == "root"
}
func (user *User) IsValid() bool {
	return user.ID > 0
}
func (user *User) String() string {
	return fmt.Sprintf("用户名: %s , 昵称: %s , 头像: %s , 登录IP: %s , 登录时间: %s ", user.UserName,
		user.UserInfo.NickName, user.UserInfo.UserAvator, user.LoginIP.String, user.LoginTime)
}

//输出精简信息
func (user *User) ToMapSimple(token string) map[string]interface{} {
	return map[string]interface{}{
		"user_id":     user.ID,
		"user_name":   user.UserName,
		"user_avator": user.UserInfo.UserAvator,
		"token":       token,
	}
}
//输出详细信息
func (user *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"user_name":   user.UserName,
		"nick_name":   user.UserInfo.NickName,
		"user_id":     user.ID,
		"is_active":   user.IsActive,
		"user_avator": user.UserInfo.UserAvator,
		"user_email":  user.UserInfo.UserEmail,
		"user_desc":   user.UserInfo.UserDesc,
		"user_addr":   user.UserInfo.UserAddr,
	}
}

//输出详细信息 带token
func (user *User) ToMapHasToken(token string) map[string]interface{} {
	resp := user.ToMap()
	resp["token"] = token
	return resp
}
