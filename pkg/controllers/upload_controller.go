package controllers

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/services"
	"net/http"
	"strings"
)

type UploadController struct {
	Service *services.UploadService
}

func NewUploadController() *UploadController {
	return &UploadController{services.NewUploadService()}
}
func (ur *UploadController) Options(ctx dotweb.Context) error {
	return ctx.WriteJsonC(http.StatusNoContent, nil)
}

//一次只能上传一张图片
func (ur *UploadController) UploadImage(ctx dotweb.Context) error {

	uploadFiles, err := ctx.Request().FormFiles()
	if err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if len(uploadFiles) > 1 {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrUploadLenNotAllow, Data: nil})
	}
	avator, err := ctx.Request().FormFile("avator")
	cover, err := ctx.Request().FormFile("cover")
	if avator != nil {
		if resp := ur.Service.SaveAvator(avator); resp.Code == 0 {
			return ctx.WriteJsonC(http.StatusCreated, resp)
		} else {
			return ctx.WriteJsonC(http.StatusBadRequest, resp)
		}
	}
	if cover != nil {
		// https://www.froala.com/wysiwyg-editor/docs/concepts/image/upload 编辑器图片上传返回格式
		if resp := ur.Service.SaveCover(cover); resp.Code == 0 {
			host := ctx.Request().Host
			if strings.Index(host, "127.0.0.1") != -1 || strings.Index(host, "localhost") != -1 {
				host = "http://" + host
			}
			return ctx.WriteJsonC(http.StatusCreated, map[string]interface{}{"link": fmt.Sprintf("%s/%s", host, resp.Data)})
		} else {
			return ctx.WriteJsonC(http.StatusBadRequest, resp)
		}
	}
	return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
}
