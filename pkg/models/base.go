package models

import (
	"github.com/track/blogserver/pkg/common"
	"time"
)

// 不带软删除的Model
type BaseModel struct {
	ID        uint      `gorm:"primary_key"` // 自增
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Api响应体结构
type Response struct {
	common.Err
	Data interface{} `json:"data"`
}
type SortValue struct {
	Key   string
	Value string
}

func (sv SortValue) IsValidValue() bool {
	return sv.Value == "desc" || sv.Value == "asc"
}
func (sv SortValue) IsValidKey() bool {
	return sv.Key == "create_time" || sv.Key == "views"
}
