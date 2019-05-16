package services

import (
	"database/sql"
	"github.com/devfeel/middleware/jwt"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/config"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/repositories"
	"github.com/track/blogserver/pkg/utils"
	"time"

	"unicode/utf8"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{repo: new(repositories.UserRepository)}
}
func (us *UserService) InitBlog() models.Response {
	//用配置默认用户 作为客户端初始化用户
	if user, ok := us.repo.GetByUserName(config.Config().DefaultClientUser); !ok {
		return models.Response{Err: common.ErrUserNoExist, Data: nil}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgGetUserInfoSucc}, Data: user.ToMap()}
	}
}

//
func (us *UserService) GetUserByUid(uid int) models.Response {
	if user, ok := us.repo.GetByUserId(uid); !ok {
		return models.Response{Err: common.ErrUserNoExist, Data: nil}
	} else {
		// 如果账号>0, 用户被禁用了
		if active := user.IsValidAccout(); active {
			return models.Response{Err: common.ErrAccoutDeny, Data: nil}
		}
		return models.Response{Err: common.Err{Msg: common.MsgGetUserInfoSucc}, Data: user.ToMap()}
	}
}

func (us *UserService) GetUsers() models.Response {
	if users, ok := us.repo.GetAll(); ok {
		return models.Response{Err: common.Err{Msg: common.MsgGetUserInfoSucc}, Data: users}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgGetUserInfoSucc}, Data: users}
	}
}
func (us *UserService) Create(v utils.Validation, user *models.User) models.Response {

	//检查用户名是否合法
	if ok := v.Range(len(user.UserName), common.LenUserNameMin, common.LenUserNameMax); ok {
		return models.Response{Err: common.ErrUserNameFormat, Data: nil}
	}
	if ok := v.NumberAndLetter(user.UserName); !ok {
		return models.Response{Err: common.ErrUserNameFormat, Data: nil}
	}
	//检查昵称长度
	if ok := v.Range(len(user.UserInfo.NickName), common.LenUserNameMin, common.LenUserNameMax); ok {
		return models.Response{Err: common.ErrUserNickNameFormat, Data: nil}
	}
	// 检查密码长度
	if ok := v.Range(len(user.UserPassWord), common.LenUserNameMin, common.LenPasswordMax); ok {
		return models.Response{Err: common.ErrUserPwdFormat, Data: nil}
	}
	if ok := v.NumberAndLetter(user.UserPassWord); !ok {
		return models.Response{Err: common.ErrUserPwdFormat, Data: nil}
	}

	//可选只要检查不超过范围就行
	if ok := v.Length(user.UserInfo.UserDesc, common.LenDesc); !ok {
		return models.Response{Err: common.ErrUserDescLen, Data: nil}
	}
	if ok := v.Length(user.UserInfo.UserAddr, common.LenAddr); !ok {
		return models.Response{Err: common.ErrUserAddrLen, Data: nil}
	}

	if len(user.UserInfo.UserEmail) > 0 {
		//检查邮箱是否合法
		if ok := v.Email(user.UserInfo.UserEmail); !ok {
			return models.Response{Err: common.ErrUserEmailFormat, Data: nil}
		}
	}

	//检查用户是否重复
	if _, ok := us.repo.GetByUserName(user.UserName); ok {
		return models.Response{Err: common.ErrUserExist, Data: nil}
	}
	//密码进行sha1加密
	password := user.UserPassWord
	user.UserPassWord = utils.Sha1(password)

	if err := us.repo.Insert(user); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	}
	return models.Response{Err: common.Err{Msg: common.MsgResistSucc}, Data: user.ToMap()}
}
func (us *UserService) DelAllUsers() models.Response {
	if err := us.repo.DelAll(); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgDelUserSucc}, Data: nil}
	}

}
func (us *UserService) DelUserByUid(uid int) models.Response {
	if err := us.repo.DelByUid(uid); err != nil {
		if _, ok := err.(common.Err); ok {
			return models.Response{Err: common.ErrUserNoExist, Data: nil}
		}
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgDelUserSucc}, Data: nil}
	}
}

////根据用户名 密码查找用户
func (us *UserService) FindUserByNameAndPwd(jwtConf *jwt.Config, v utils.Validation, ip, username, password string) models.Response {
	if ok := v.Range(len(username), common.LenUserNameMin, common.LenUserNameMax); ok {
		return models.Response{Err: common.ErrUserNameFormat, Data: nil}
	}
	if ok := v.Range(len(password), common.LenUserNameMin, common.LenPasswordMax); ok {
		return models.Response{Err: common.ErrUserPwdFormat, Data: nil}
	}
	if ok := v.NumberAndLetter(username); !ok {
		return models.Response{Err: common.ErrUserNameFormat, Data: nil}
	}
	if ok := v.NumberAndLetter(password); !ok {
		return models.Response{Err: common.ErrUserPwdFormat, Data: nil}
	}
	if user, ok := us.repo.GetByUserNameAndPwd(username, utils.Sha1(password)); !ok {
		return models.Response{Err: common.ErrUserLogin, Data: nil}
	} else {
		if valid := user.IsValid(); valid {
			//登录成功修改用户登录信息
			user.LoginIP = sql.NullString{String: ip, Valid: true}
			user.LoginTime = time.Now()
			if err := us.repo.Update(user, map[string]interface{}{"login_ip": user.LoginIP, "login_time": user.LoginTime}); err != nil {
				return models.Response{Err: common.ErrInternal, Data: nil}
			}
		} else {
			return models.Response{Err: common.ErrInternal, Data: nil}
		}
		token, _ := jwt.GeneratorToken(jwtConf, map[string]interface{}{"user_ip": user.LoginIP.String, "user_id": int(user.ID)})
		return models.Response{Err: common.Err{Msg: common.MsgLoginSucc}, Data: user.ToMapSimple(token)}
	}
}

//
////根据用户id修改用户信息
func (us *UserService) UpdateUserInfoByUid(v utils.Validation, uid int, params map[string]interface{}) models.Response {
	//根据map参数,检查要修改的密码,头像,昵称,邮箱,描述,地址
	user, ok := us.repo.GetByUserId(uid)
	if !ok {
		return models.Response{Err: common.ErrUserNoExist, Data: nil}
	}
	//TODO 密码实际修改需要原密码, 或者手机短信验证码验证, 此处不做处理
	if password, ok := params["user_pass_word"]; ok {
		// 检查密码长度
		if ok := v.Range(len(password.(string)), common.LenUserNameMin, common.LenPasswordMax); ok {
			return models.Response{Err: common.ErrUserPwdFormat, Data: nil}
		}
		params["user_pass_word"] = utils.Sha1(password.(string))
		if err := us.repo.Update(&user, params); err != nil {
			return models.Response{Err: common.ErrInternal, Data: nil}
		} else {
			return models.Response{Err: common.Err{Msg: common.MsgUpdateUserInfoSucc}, Data: nil}
		}
	}

	//头像不检查, 可以在上传处做 大小,格式等限制

	if nickname, ok := params["nick_name"]; ok {
		//检查昵称长度
		if ok := v.Range(utf8.RuneCountInString(nickname.(string)), common.LenUserNameMin, common.LenUserNameMax); ok {
			return models.Response{Err: common.ErrUserNickNameFormat, Data: nil}
		}
	}
	if desc, ok := params["user_desc"]; ok {
		if ok := v.Length(desc.(string), common.LenDesc); !ok {
			return models.Response{Err: common.ErrUserDescLen, Data: nil}
		}
	}
	if addr, ok := params["user_addr"]; ok {
		if ok := v.Length(addr.(string), common.LenAddr); !ok {
			return models.Response{Err: common.ErrUserAddrLen, Data: nil}
		}
	}
	if email, ok := params["user_email"]; ok {
		if ok := v.Email(email.(string)); !ok {
			return models.Response{Err: common.ErrUserEmailFormat, Data: nil}
		}
	}
	if err := us.repo.Update(&user.UserInfo, params); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	}
	return models.Response{Err: common.Err{Msg: common.MsgUpdateUserInfoSucc}, Data: nil}
}
