package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
)
//对应的数据库中表: users user_infos
type UserRepository struct {
	*BaseRepository
}

//获取全部用户数据
func (ur *UserRepository) GetAll() ([]models.User, bool) {
	users := make([]models.User, 0)
	return users, ur.SelectMany(func(db *gorm.DB) *gorm.DB {
		return db.Preload("UserInfo")

	}, &users)
}

//根据用户Id获取用户
func (ur *UserRepository) GetByUserId(id int) (*models.User, bool) {
	user := &models.User{}
	found := ur.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}, user)
	if found {
		ur.db().Model(user).Related(&user.UserInfo)
	}
	return user, found
}

//根据用户名(唯一)获取用户
func (ur *UserRepository) GetByUserName(username string) (*models.User, bool) {
	user := &models.User{}
	found := ur.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("user_name = ?", username)
	}, user)
	if found {
		ur.db().Model(user).Related(&user.UserInfo)
	}
	return user, found
}

//根据用户名和密码获取用户
func (ur *UserRepository) GetByUserNameAndPwd(username, password string) (*models.User, bool) {
	user := &models.User{}
	found := ur.Select(func(db *gorm.DB) *gorm.DB {
		return db.Where("user_name = ? and user_pass_word = ?", username, password)
	}, user)
	if found {
		ur.db().Model(user).Related(&user.UserInfo)
	}
	return user, found
}

//删除表中全部用户
func (us *UserRepository) DelAll() error {
	if err := us.Del(func(db *gorm.DB) *gorm.DB {
		return db
	}, &models.User{}); err != nil {
		return err
	} else {
		return us.Del(func(db *gorm.DB) *gorm.DB {
			return db
		}, &models.UserInfo{})
	}
}

//根据用户id从表中删除某个用户,与其关联的用户信息
func (us *UserRepository) DelByUid(uid int) error {
	if user, found := us.GetByUserId(uid); found {
		err := us.Del(func(db *gorm.DB) *gorm.DB {
			return db
		}, user)
		err = us.Del(func(db *gorm.DB) *gorm.DB {
			return db
		}, &user.UserInfo)
		return err
	} else {
		return common.ErrUserNoExist
	}
}
