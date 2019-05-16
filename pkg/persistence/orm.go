package persistence

import (

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/track/blogserver/pkg/config"
)

// 持久化 之 gorm  mysql
var orm *gorm.DB

//根据配置初始化gorm 打开数据库连接
func init() {
	conf := config.Config().DBCfg
	var err error
	orm, err = gorm.Open(conf.Dtype,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			conf.User, conf.Password, conf.Addr, conf.Port, conf.Name))
	if err != nil {
		panic(err)
	}
	orm.LogMode(conf.Debug)
}

// 获取gorm全局实例
func GetOrm() *gorm.DB {
	return orm
}
