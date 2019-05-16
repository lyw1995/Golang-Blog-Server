package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/services"
	"github.com/track/blogserver/pkg/utils"
	"net/http"
	"strconv"

)

type ExtendController struct {
	Service *services.ExtendSerivce
}

func NewExtendController() *ExtendController {
	return &ExtendController{Service: services.NewExtendSerivce()}
}

func (ec *ExtendController) Options(ctx dotweb.Context) error {
	return ctx.WriteJsonC(http.StatusNoContent, nil)
}

//获取日活ip,独立ip,文章总访问量,文章总篇数
func (ec *ExtendController) ExtByQueryParams(ctx dotweb.Context) error {
	op := ctx.QueryString("op")
	uid := ctx.QueryString("uid")
	if len(op) > 0 {
		switch op {
		case "all":
			if id, err := strconv.Atoi(uid); err != nil {
				return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
			} else {
				//判断响应码  不为0,
				if resp := ec.Service.GetAll(id); resp.Code == 0 {
					return ctx.WriteJsonC(http.StatusOK, resp)
				} else {
					return ctx.WriteJsonC(http.StatusBadRequest, resp)
				}
			}
		}
	}
	return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
}

//采集文章 refer: 来源, 目前仅支持csdn
func (ec *ExtendController) ArticleCollection(ctx dotweb.Context) error {
	uid := ctx.FormValue("uid")
	origin := ctx.FormValue("origin")
	cid := ctx.FormValue("cid")
	refer := ctx.FormValue("refer")
	action := ctx.FormValue("action")
	collectUrl := ctx.FormValue("url")

	v := utils.Validation{}
	if v.Required(uid) || v.Required(cid) || v.Required(origin) || v.Required(refer) || v.Required(action) || v.Required(collectUrl) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	//检测必填值
	if action != "publish" && action != "draft" {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if refer != "csdn" {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	article := &models.Article{}
	//直接发布还是存草稿
	if action == "publish" {
		article.State = 0
	} else {
		article.State = 2
	}
	if uid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		article.UserID = uint(uid)
	}
	if origin, err := strconv.Atoi(origin); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		article.Origin = origin
	}
	if cid, err := strconv.Atoi(cid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		article.CategoryItemID = uint(cid)
	}

	if resp := ec.Service.CollectByCsdn(collectUrl, article); resp.Code == 0 {
		return ctx.WriteJsonC(http.StatusCreated, resp)
	} else {
		return ctx.WriteJsonC(http.StatusBadRequest, resp)
	}
}
