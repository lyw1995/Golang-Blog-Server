package controllers

import (
	"database/sql"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/middleware/jwt"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/config"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/services"
	"github.com/track/blogserver/pkg/utils"
	"net/http"
	"strconv"
	"time"
)
//用户控制中心
type UserController struct {
	service *services.UserService
}

func NewUserController() *UserController {
	return &UserController{service: services.NewUserService()}
}

func (uc *UserController) Options(ctx dotweb.Context) error {
	return ctx.WriteJsonC(http.StatusNoContent, nil)
}

//TODO 暂不具体实现,没有超级管理员
//获取全部用户
func (uc *UserController) GetUsers(ctx dotweb.Context) error {
	return ctx.WriteJsonC(http.StatusOK, uc.service.GetUsers())
}
//博客客户端初始化
func (uc *UserController) InitBlog(ctx dotweb.Context) error {
	if resp := uc.service.InitBlog(); resp.Code == 0 {
		return ctx.WriteJsonC(http.StatusOK, resp)
	} else {
		return ctx.WriteJsonC(http.StatusBadRequest, resp)
	}
}
//获取某个用户
func (uc *UserController) GetUserByUid(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	//效验uid 不可空并且为数字
	v := utils.Validation{}
	if v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if uid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		//判断响应码
		if resp := uc.service.GetUserByUid(uid); resp.Code == 0 {
			return ctx.WriteJsonC(http.StatusOK, resp)
		} else {
			return ctx.WriteJsonC(http.StatusBadRequest, resp)
		}
	}
}

//创建用户
func (uc *UserController) CreateUser(ctx dotweb.Context) error {
	//如果是正式环境.将关闭新注册用户
	if config.Config().EnvProd {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrApiClose, Data: nil})
	}
	username := ctx.FormValue("username") //必填
	password := ctx.FormValue("password") //必填
	nickname := ctx.FormValue("nickname") //必填
	avator := ctx.FormValue("avator")     //可选
	desc := ctx.FormValue("desc")         //可选
	addr := ctx.FormValue("addr")         //可选
	email := ctx.FormValue("email")       //可选

	v := utils.Validation{}
	if v.Required(username) || v.Required(password) || v.Required(nickname) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}

	user := models.User{
		UserName:     username,
		UserPassWord: password,
		LoginTime:    time.Now(),
		LoginIP:      sql.NullString{String: ctx.RemoteIP(), Valid: true},
		UserInfo: models.UserInfo{
			UserAvator: avator,
			UserDesc:   desc,
			UserAddr:   addr,
			UserEmail:  email,
			NickName:   nickname,
		},
	}

	if resp := uc.service.Create(v, &user); resp.Code == 0 {
		return ctx.WriteJsonC(http.StatusCreated, resp)
	} else {
		return ctx.WriteJsonC(http.StatusBadRequest, resp)
	}

}

//删除全部用户
func (uc *UserController) DelAllUser(ctx dotweb.Context) error {
	if resp := uc.service.DelAllUsers(); resp.Code == 0 {
		return ctx.WriteJsonC(http.StatusNoContent, resp)
	} else {
		return ctx.WriteJsonC(http.StatusBadRequest, resp)
	}
}

//删除某个用户
func (uc *UserController) DelByUid(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	v := utils.Validation{}
	if v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if uid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if resp := uc.service.DelUserByUid(uid); resp.Code == 0 {
			return ctx.WriteJsonC(http.StatusNoContent, resp)
		} else {
			return ctx.WriteJsonC(http.StatusBadRequest, resp)
		}
	}
}

//修改某个用户资料
func (uc *UserController) UpdateUserInfoByUid(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	v := utils.Validation{}
	if v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	//密码 头像 昵称 描述 地址 邮箱  传啥改啥
	params := make(map[string]interface{})

	if password := ctx.FormValue("password"); !v.Required(password) {
		params["user_pass_word"] = password
	}

	if avator := ctx.FormValue("avator"); !v.Required(avator) {
		params["user_avator"] = avator
	}

	if nickname := ctx.FormValue("nickname"); !v.Required(nickname) {
		params["nick_name"] = nickname
	}

	if desc := ctx.FormValue("desc"); !v.Required(desc) {
		params["user_desc"] = desc
	}

	if addr := ctx.FormValue("addr"); !v.Required(addr) {
		params["user_addr"] = addr
	}

	if email := ctx.FormValue("email"); !v.Required(email) {
		params["user_email"] = email
	}
	if ok := len(params) <= 0; ok {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrUpdateParams, Data: nil})
	}
	if uid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if resp := uc.service.UpdateUserInfoByUid(v, uid, params); resp.Code == 0 {
			return ctx.WriteJsonC(http.StatusOK, resp)
		} else {
			return ctx.WriteJsonC(http.StatusBadRequest, resp)
		}
	}
}
func parseJwtConfig(c interface{}, exists bool) (config *jwt.Config) {
	if c == nil || !exists {
		return nil
	}
	return c.(*jwt.Config)
}

//用户登录
func (uc *UserController) UserLogin(ctx dotweb.Context) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	v := utils.Validation{}
	if v.Required(username) || v.Required(password) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if jwtConf := parseJwtConfig(ctx.AppItems().Get(config.Config().SecretKey)); jwtConf == nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrInternal, Data: nil})
	} else {
		realIp := ctx.Request().QueryHeader("X-Real-IP")
		if resp := uc.service.FindUserByNameAndPwd(jwtConf, v, realIp, username, password); resp.Code == 0 {
			return ctx.WriteJsonC(http.StatusOK, resp)
		} else {
			return ctx.WriteJsonC(http.StatusBadRequest, resp)
		}
	}
}

//用户登出
func (uc *UserController) UserLogout(ctx dotweb.Context) error {
	return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrApiClose, Data: nil})
}
