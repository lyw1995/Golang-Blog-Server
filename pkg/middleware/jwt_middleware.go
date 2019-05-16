package middleware

import (
	"errors"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/middleware/jwt"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/config"
	"github.com/track/blogserver/pkg/models"
	"net/http"
	"strconv"
	"time"
)

//对于admin接口 使用jwt auth
func NewJwtMiddleware(app *dotweb.DotWeb) dotweb.Middleware {
	option := &jwt.Config{
		TTL:           time.Hour * 24, // token 24小时有效期
		ContextKey:    config.Config().ServerName,
		SigningKey:    []byte(config.Config().SecretKey),
		SigningMethod: jwt.SigningMethodHS256,
		ExceptionHandler: func(ctx dotweb.Context, err error) {
			ctx.WriteJsonC(http.StatusUnauthorized, models.Response{Err: common.ErrAuthorized, Data: nil})
		},
		AddonValidator: func(jwtCon *jwt.Config, ctx dotweb.Context) error {
			//payload camis 用户id ip判断,
			uid := ctx.RouterParams().ByName("uid")
			//仅测试 获取用户列表 不需要判断uid
			if len(uid) <= 0 {
				return nil
			}
			if tuid, err := strconv.Atoi(uid); err != nil {
				return errors.New("uid input format is not allow")
			} else {
				jwtobj, exists := ctx.Items().Get(config.Config().ServerName)
				if !exists {
					return errors.New("no token exists")
				}
				//不判断ip了..ip变化有点快
				//realIp := ctx.Request().QueryHeader("X-Real-IP")
				jwtmap := jwtobj.(map[string]interface{})
				//jwtUserIp := jwtmap["user_ip"].(string)
				jwtUserId := jwtmap["user_id"].(float64)
				//if jwtUserIp != realIp {
				//	return errors.New("user_ip is not match")
				//}
				if int(jwtUserId) != tuid {
					return errors.New("user_id is not match")
				}
			}
			return nil
		},
		//use header
		Extractor: jwt.ExtractorFromHeader,
	}

	app.Items.Set(config.Config().SecretKey, option)

	return jwt.Middleware(option)
}
