package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/repositories"
	"github.com/track/blogserver/pkg/services"
	"github.com/track/blogserver/pkg/utils"
	"net/http"
	"strconv"
	"time"

)

type ArticleController struct {
	Service *services.ArticleService
}

func NewArticleController() *ArticleController {
	return &ArticleController{services.NewArticleService()}
}
func (ac *ArticleController) Options(ctx dotweb.Context) error {
	return ctx.WriteJsonC(http.StatusNoContent, nil)
}

//热门文章和最新文章
func (ac *ArticleController) GetHotAndNewArticles(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	//参数效验
	v := utils.Validation{}
	if v.Required(uid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		r := new(repositories.RedisRepository)
		return ctx.WriteJsonC(http.StatusOK, models.Response{Err: common.Err{Msg: common.MsgGetArticleSucc}, Data: r.SelectHotAndNewArticle(tuid)})
	}
}

//获取全部文章(已发布,草稿)
func (ac *ArticleController) GetUserArticles(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	page := ctx.QueryString("page") //分页选填
	published := ctx.QueryString("published")
	sort := ctx.QueryString("sort")
	ipage := 0 //第0页开始 默认返回common.LenLimit = 10 条
	//参数效验
	v := utils.Validation{}
	if v.Required(uid) || v.Required(published) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	//检测必填值
	if published != "true" && published != "false" {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}

	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		//如果有page,并且为正整数 则赋值
		if !v.Required(page) {
			if tpage, err := strconv.Atoi(page); err != nil {
				return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
			} else {
				ipage = tpage
			}
		}
		if published == "true" {
			return ctx.WriteJsonC(http.StatusOK, ac.Service.GetPageBy(tuid, ipage, sort))
		} else {
			return ctx.WriteJsonC(http.StatusOK, ac.Service.GetDraftPageBy(tuid, ipage, sort))
		}
	}
}

//根据状态获取文章(已发布的 草稿箱的)
func (ac *ArticleController) UpdateArticleWithState(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	aid := ctx.RouterParams().ByName("aid")
	title := ctx.FormValue("title")     //必填
	content := ctx.FormValue("content") //必填
	origin := ctx.FormValue("origin")   //必填 1 原创 0 转载
	cid := ctx.FormValue("cid")         //必填 分类id
	state := ctx.QueryString("state")
	v := utils.Validation{}
	if v.Required(uid) || v.Required(aid) || v.Required(cid) || v.Required(title) || v.Required(content) || v.Required(origin) || v.Required(state) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	//检测state值规范
	if state != "publish" && state != "draft" {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	article := &models.Article{
		CreateTime: time.Now(),
		CreatedAt:  time.Now(),
	}
	//效验用户id 合法性
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		article.UserID = uint(tuid)
	}
	//效验文章id 合法性
	if taid, err := strconv.Atoi(aid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		article.ID = uint(taid)
	}
	//效验分类id 合法性
	if tcid, err := strconv.Atoi(cid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		article.CategoryItemID = uint(tcid)
	}
	//重新解码一下
	if value, err := utils.DecodeUriCompontent(title); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		article.Title = value
	}
	if value, err := utils.DecodeUriCompontent(content); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		article.Content = value
	}
	//效验 来源
	if torigin, err := strconv.Atoi(origin); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if torigin != 0 && torigin != 1 {
			return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
		}
		article.Origin = torigin
	}
	if resp := ac.Service.UpdateWithState(article, state); resp.Code == 0 {
		return ctx.WriteJsonC(http.StatusCreated, resp)
	} else {
		return ctx.WriteJsonC(http.StatusBadRequest, resp)
	}
}

//根据状态获取文章(已发布的 草稿箱的)
func (ac *ArticleController) GetArticleWithState(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	aid := ctx.RouterParams().ByName("aid")
	state := ctx.QueryString("state")
	v := utils.Validation{}
	if v.Required(uid) || v.Required(aid) || v.Required(state) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	//检测state值规范
	if state != "publish" && state != "draft" {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if taid, err := strconv.Atoi(aid); err != nil {
			return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
		} else {
			if resp := ac.Service.GetArticleWithState(tuid, taid, state); resp.Code == 0 {
				return ctx.WriteJsonC(http.StatusOK, resp)
			} else {
				return ctx.WriteJsonC(http.StatusBadRequest, resp)
			}
		}
	}
	return nil
}

//用户在获取某篇文章(分类cid可选)
func (ac *ArticleController) GetArticleById(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	cid := ctx.RouterParams().ByName("cid") //可选
	aid := ctx.RouterParams().ByName("aid")

	v := utils.Validation{}
	if v.Required(uid) || v.Required(aid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if taid, err := strconv.Atoi(aid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {

		if tuid, err := strconv.Atoi(uid); err != nil {
			return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
		} else {
			//判断是否是分类路径进来的
			if !v.Required(cid) {
				if tcid, err := strconv.Atoi(cid); err != nil {
					return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
				} else {
					if resp := ac.Service.GetArticleByUidAndCIdAndAid(tuid, tcid, taid); resp.Code == 0 {
						return ctx.WriteJsonC(http.StatusOK, resp)
					} else {
						return ctx.WriteJsonC(http.StatusBadRequest, resp)
					}
				}
			}
			if resp := ac.Service.GetArticleByUidAndAid(tuid, taid, ctx.Request().RemoteIP()); resp.Code == 0 {
				return ctx.WriteJsonC(http.StatusOK, resp)
			} else {
				return ctx.WriteJsonC(http.StatusBadRequest, resp)
			}
		}
	}
}

//用户在获取分类下全部文章
func (ac *ArticleController) GetCategoryArticles(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	cid := ctx.RouterParams().ByName("cid")
	page := ctx.QueryString("page") //分页选填
	sort := ctx.QueryString("sort")
	ipage := 0 //第0页开始 默认返回common.LenLimit = 10 条
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
			//如果有page,并且为正整数 则赋值
			if !v.Required(page) {
				if tpage, err := strconv.Atoi(page); err != nil {
					return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
				} else {
					ipage = tpage
				}
			}
			if resp := ac.Service.GetCategoryItemArticlesPageBy(tuid, tcid, ipage, sort); resp.Code == 0 {
				return ctx.WriteJsonC(http.StatusOK, resp)
			} else {
				return ctx.WriteJsonC(http.StatusBadRequest, resp)
			}
		}
	}
}

//用户在分类下创建文章
func (ac *ArticleController) CreateArticleByCid(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	cid := ctx.RouterParams().ByName("cid")
	title := ctx.FormValue("title")     //必填
	content := ctx.FormValue("content") //必填
	origin := ctx.FormValue("origin")   //必填 1 原创 0 转载
	publish := ctx.QueryString("publish")

	v := utils.Validation{}
	if v.Required(uid) || v.Required(cid) || v.Required(title) || v.Required(content) || v.Required(origin) || v.Required(publish) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	article := &models.Article{
		CreateTime: time.Now(),
	}
	//重新解码一下
	if value, err := utils.DecodeUriCompontent(title); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		article.Title = value
	}
	if value, err := utils.DecodeUriCompontent(content); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		article.Content = value
	}

	//效验y用户id 合法性
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		article.UserID = uint(tuid)
	}
	//效验分类id 合法性
	if tcid, err := strconv.Atoi(cid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		article.CategoryItemID = uint(tcid)
	}
	//效验 来源
	if torigin, err := strconv.Atoi(origin); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if torigin != 0 && torigin != 1 {
			return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
		}
		article.Origin = torigin
	}

	//检测必填值
	if publish != "true" && publish != "false" {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}

	//效验是发布还是保存草稿
	if publish == "true" {
		article.State = 0
	} else {
		article.State = 2
	}
	if resp := ac.Service.Create(article); resp.Code == 0 {
		return ctx.WriteJsonC(http.StatusCreated, resp)
	} else {
		return ctx.WriteJsonC(http.StatusBadRequest, resp)
	}
}

//用户在分类下修改文章
func (ac *ArticleController) UpdateArticleByCid(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	cid := ctx.RouterParams().ByName("cid")
	aid := ctx.RouterParams().ByName("aid")

	//参数效验
	v := utils.Validation{}
	if v.Required(uid) || v.Required(cid) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}

	params := make(map[string]interface{})

	if title := ctx.FormValue("title"); !v.Required(title) {
		params["title"] = title
	}

	if content := ctx.FormValue("content"); !v.Required(content) {
		params["content"] = content
	}

	if ok := len(params) <= 0; ok {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}
	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if tcid, err := strconv.Atoi(cid); err != nil {
			return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
		} else {
			if taid, err := strconv.Atoi(aid); err != nil {
				return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
			} else {
				if resp := ac.Service.UpdateArticle(tuid, tcid, taid, params); resp.Code == 0 {
					return ctx.WriteJsonC(http.StatusNoContent, resp)
				} else {
					return ctx.WriteJsonC(http.StatusBadRequest, resp)
				}
			}

		}
	}
}

//用户在分类删除全部文章
func (ac *ArticleController) DelArticlesByCid(ctx dotweb.Context) error {
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
			if resp := ac.Service.DelArticlesByUidAndCid(tuid, tcid); resp.Code == 0 {
				return ctx.WriteJsonC(http.StatusNoContent, resp)
			} else {
				return ctx.WriteJsonC(http.StatusBadRequest, resp)
			}
		}
	}
}

//用户在分类删除单篇文章
func (ac *ArticleController) DelArticleByCid(ctx dotweb.Context) error {
	uid := ctx.RouterParams().ByName("uid")
	cid := ctx.RouterParams().ByName("cid")
	aid := ctx.RouterParams().ByName("aid")
	state := ctx.QueryString("state")
	//参数效验
	v := utils.Validation{}
	if v.Required(uid) || v.Required(cid) || v.Required(aid) || v.Required(state) {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}

	//检测必填值
	if state != "publish" && state != "draft" {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	}

	if tuid, err := strconv.Atoi(uid); err != nil {
		return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
	} else {
		if tcid, err := strconv.Atoi(cid); err != nil {
			return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
		} else {
			if taid, err := strconv.Atoi(aid); err != nil {
				return ctx.WriteJsonC(http.StatusBadRequest, models.Response{Err: common.ErrClientParams, Data: nil})
			} else {
				if resp := ac.Service.DelArticleByUidAndCidAndAid(tuid, tcid, taid, state); resp.Code == 0 {
					return ctx.WriteJsonC(http.StatusNoContent, resp)
				} else {
					return ctx.WriteJsonC(http.StatusBadRequest, resp)
				}
			}

		}
	}
}
