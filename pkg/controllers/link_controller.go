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

type LinkController struct {
	Service *services.LinkSerivce
}

func NewLinkController() *LinkController {
	return &LinkController{Service: services.NewLinkSerivce()}
}

func (lc *LinkController) Options(ctx dotweb.Context) error {
	return ctx.WriteJsonC(http.StatusNoContent, nil)
}

//创建某个用户的链接
func (lc *LinkController) CreateFriendlyLink(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	avator := ctx.FormValue("avator") //必填
	linkname := ctx.FormValue("name") //必填
	linkurl := ctx.FormValue("url")   //必填
	//参数效验
	v := utils.Validation{}
	if v.Required(uid) || v.Required(avator) || v.Required(linkname) || v.Required(linkurl) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}

	if id, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		//判断响应码  不为0, 则创建失败
		if resp := lc.Service.Create(v, &models.FriendlyLink{UserID: uint(id), LinkName: linkname, Avator: avator, LinkUrl: linkurl}); resp.Code == 0 {
			return ctx.WriteJsonC(http.StatusCreated, resp)
		} else {
			return ctx.WriteJsonC(http.StatusBadRequest, resp)
		}
	}
}

//获取某个用户所有友链
func (lc *LinkController) GetLinks(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	//参数效验
	v := utils.Validation{}
	if v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	//判断参数是否为int
	if tid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		//判断响应码  不为0, 则数据获取失败
		if resp := lc.Service.GetAllBy(tid); resp.Code == 0 {
			return ctx.WriteJsonC(http.StatusOK, resp)
		} else {
			return ctx.WriteJsonC(http.StatusBadRequest, resp)
		}
	}
}

//删除某个用户某个友链
func (lc *LinkController) DelLinkById(ctx dotweb.Context) error {
	id := ctx.RouterParams().ByName("lid")
	uid := ctx.RouterParams().ByName("uid")
	//参数效验
	v := utils.Validation{}
	if v.Required(id) || v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	//判断参数是否为int
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if lid, err := strconv.Atoi(id); err != nil {
			return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
		} else {
			//判断响应码  不为0, 则删除失败
			if resp := lc.Service.DelUserLinkById(tuid, lid); resp.Code == 0 {
				return ctx.WriteJsonC(http.StatusNoContent, resp)
			} else {
				return ctx.WriteJsonC(http.StatusBadRequest, resp)
			}
		}
	}
}

//获取某个用户某个友链
func (lc *LinkController) GetLinkById(ctx dotweb.Context) error {
	id := ctx.RouterParams().ByName("lid")
	uid := ctx.RouterParams().ByName("uid")
	//参数效验
	v := utils.Validation{}
	if v.Required(id) || v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	//判断参数是否为int
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if lid, err := strconv.Atoi(id); err != nil {
			return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
		} else {
			//判断响应码  不为0, 则数据获取失败
			if resp := lc.Service.GetBy(tuid, lid); resp.Code == 0 {
				return ctx.WriteJsonC(http.StatusOK, resp)
			} else {
				return ctx.WriteJsonC(http.StatusBadRequest, resp)
			}
		}
	}
}

//删除该用户所有友链
func (lc *LinkController) DelLinks(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	//参数效验
	v := utils.Validation{}
	if v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	//判断参数是否为int
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		//判断响应码  不为0, 则删除失败
		if resp := lc.Service.DelAllByUserId(tuid); resp.Code == 0 {
			return ctx.WriteJsonC(http.StatusOK, resp)
		} else {
			return ctx.WriteJsonC(http.StatusBadRequest, resp)
		}
	}
}

//修改某个链接
func (lc *LinkController) UpdateLinkById(ctx dotweb.Context) error {
	id := ctx.RouterParams().ByName("lid")
	uid := ctx.RouterParams().ByName("uid")
	//参数效验
	v := utils.Validation{}
	if v.Required(id) || v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}

	params := make(map[string]interface{})

	if nickname := ctx.FormValue("url"); !v.Required(nickname) {
		params["link_url"] = nickname
	}

	if nickname := ctx.FormValue("name"); !v.Required(nickname) {
		params["link_name"] = nickname
	}

	if nickname := ctx.FormValue("avator"); !v.Required(nickname) {
		params["avator"] = nickname
	}
	if ok := len(params) <= 0; ok {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrUpdateParams, Data: nil})
	}
	//判断参数是否为int
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if lid, err := strconv.Atoi(id); err != nil {
			return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
		} else {
			//判断响应码  不为0, 则修改获取失败
			if resp := lc.Service.UpdateBy(tuid, lid, params); resp.Code == 0 {
				return ctx.WriteJsonC(http.StatusOK, resp)
			} else {
				return ctx.WriteJsonC(http.StatusBadRequest, resp)
			}
		}
	}
}
