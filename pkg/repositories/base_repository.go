package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/track/blogserver/pkg/persistence"
)

//动态连接查询
type Where func(*gorm.DB) *gorm.DB

//通用输入输出
type Out interface{}

//DAO 接口
type IRepository interface {
	Insert(out Out) error
	Exec(where Where, out Out, offset int, limit int) bool
	Select(where Where, out Out) bool
	SelectMany(where Where, out Out) bool
	Del(where Where, out Out) error
	Update(out Out, params map[string]interface{}) error
}

type BaseRepository struct {
}

//接口实现自检
var _ IRepository = &BaseRepository{}

//获取数据库实例
func (br *BaseRepository) db() *gorm.DB {
	return persistence.GetOrm()
}

//通用查询
func (br *BaseRepository) Exec(where Where, out Out, offset int, limit int) bool {
	return !br.db().Scopes(where).Offset(offset).Limit(limit).Find(out).RecordNotFound()
}

//插入数据
func (br *BaseRepository) Insert(out Out) error {
	return br.db().Create(out).Error
}

//查单个
func (br *BaseRepository) Select(where Where, out Out) bool {
	return br.Exec(where, out, 0, 1)
}

//查全部
func (br *BaseRepository) SelectMany(where Where, out Out) bool {
	return br.Exec(where, out, 0, -1)
}

//更新数据 ,单个 或者 多个
func (br *BaseRepository) Update(out Out, params map[string]interface{}) error {
	return br.db().Model(out).Update(params).Error
}

//完全更新
func (br *BaseRepository) Save(out Out) error {
	return br.db().Save(out).Error
}

//Count
func (br *BaseRepository) Count(out Out, where Where) (err error, count int) {
	return br.db().Model(out).Scopes(where).Count(&count).Error, count
}

//删除数据, 单个或者多个
func (br *BaseRepository) Del(where Where, out Out) error {
	return br.db().Scopes(where).Delete(out).Error
}
