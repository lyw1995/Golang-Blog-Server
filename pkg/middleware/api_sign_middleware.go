package middleware

import (
	"github.com/devfeel/dotweb"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/config"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/utils"
	"net/http"
	"strings"
)

//对于user接口使用 sha1 sign 验证 ...dotwebmiddleware中间件有问题
type ApiSignMiddleware struct {
	dotweb.BaseMiddlware
}

func (asm *ApiSignMiddleware) Handle(ctx dotweb.Context) error {
	if sign := ctx.Request().QueryHeader("Sign"); len(sign) <= 0 {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrSignParams, Data: nil})
	} else {
		uri := ctx.Request().RequestURI
		if index := strings.Index(uri, "?"); index != -1 {
			uri = uri[:index]
		}
		if ok := checkSign(sign, uri); !ok {
			return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrSignParams, Data: nil})
		}
		return asm.Next(ctx)
	}
}

//验证签名 (requestUri(不含query)+secret)
func checkSign(sign, uri string) bool {
	result := utils.Md5(uri + config.Config().SecretKey)
	return result == sign
}
