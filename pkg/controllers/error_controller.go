package controllers

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/track/blogserver/pkg/common"
	"net/http"
)

type ErrorController struct {
}

func NewErrorController() *ErrorController {
	return &ErrorController{}
}

//配置404
func (ec *ErrorController) NotFound(context dotweb.Context) {
	context.WriteJsonC(http.StatusNotFound, common.ErrNotFound)
}

//配置500
func (ec *ErrorController) Internal(context dotweb.Context, err error) {
	errCust := common.Err{
		Code: common.ErrInternal.Code,
		Msg:  fmt.Sprintf("%s %s", common.ErrInternal.Msg, err.Error()),
	}
	context.WriteJsonC(http.StatusInternalServerError, errCust)
}

//配置405
func (ec *ErrorController) MethodNotAllowed(context dotweb.Context) {
	context.WriteJsonC(http.StatusMethodNotAllowed, common.ErrMethodNotAllow)
}
