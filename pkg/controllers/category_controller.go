package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/services"
	"github.com/track/blogserver/pkg/utils"
	"net/http"
	"strconv"
	"time"
)

type CategoryController struct {
	Service *services.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{services.NewCategoryService()}
}
func (cc *CategoryController) Options(ctx dotweb.Context) error {

	return ctx.WriteJsonC(http.StatusNoContent, nil)
}

//获取全部分类
func (cc *CategoryController) GetCategorys(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	//参数效验
	v := utils.Validation{}
	if v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		return ctx.WriteJsonC(http.StatusOK, cc.Service.GetCategorys(tuid))
	}

}

//获取个人分类所有子条目
func (cc *CategoryController) GetPersonalCategorys(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	//参数效验
	v := utils.Validation{}
	if v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		return ctx.WriteJsonC(http.StatusOK, cc.Service.GetPersonalCategorys(tuid))
	}

}

//创建分类
func (cc *CategoryController) CreateCategory(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	name := ctx.FormValue("name") //必填
	//参数效验
	v := utils.Validation{}
	if v.Required(uid) || v.Required(name) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	//重新解码一下
	if value, err := utils.DecodeUriCompontent(name); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		name = value
	}
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if resp := cc.Service.CreateCategoryItem(&models.CategoryItem{
			Name:       name,
			CreateTime: time.Now(),
			UserID:     uint(tuid),
		}); resp.Code == 0 {
			return ctx.WriteJsonC(http.StatusCreated, resp)
		} else {
			return ctx.WriteJsonC(http.StatusBadRequest, resp)
		}
	}
}

//修改某个分类名称
func (cc *CategoryController) UpdateCategoryById(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	cid := ctx.RouterParams().ByName("cid")
	cname := ctx.FormValue("name") //必填 要先编码一下
	if value, err := utils.DecodeUriCompontent(cname); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		cname = value
	}
	//参数效验
	v := utils.Validation{}
	if v.Required(uid) || v.Required(cname) || v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	//重新解码一下
	if value, err := utils.DecodeUriCompontent(cname); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		cname = value
	}
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if tcid, err := strconv.Atoi(cid); err != nil {
			return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
		} else {
			if resp := cc.Service.UpdateCategoryItemBy(tuid, tcid, map[string]interface{}{"name": cname}); resp.Code == 0 {
				return ctx.WriteJsonC(http.StatusOK, resp)
			} else {
				return ctx.WriteJsonC(http.StatusBadRequest, resp)
			}
		}
	}
}

//删除某个分类
func (cc *CategoryController) DelCategoryById(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	cid := ctx.RouterParams().ByName("cid")
	//参数效验
	v := utils.Validation{}
	if v.Required(uid) || v.Required(cid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if tcid, err := strconv.Atoi(cid); err != nil {
			return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
		} else {
			if resp := cc.Service.DelAllCategoryItemBy(tuid, tcid); resp.Code == 0 {
				return ctx.WriteJsonC(http.StatusNoContent, resp)
			} else {
				return ctx.WriteJsonC(http.StatusBadRequest, resp)
			}
		}
	}
}

//删除全部分类
func (cc *CategoryController) DelCategorys(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	//参数效验
	v := utils.Validation{}
	if v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if resp := cc.Service.DelAllCategoryItems(tuid); resp.Code == 0 {
			return ctx.WriteJsonC(http.StatusNoContent, resp)
		} else {
			return ctx.WriteJsonC(http.StatusBadRequest, resp)
		}
	}

}
