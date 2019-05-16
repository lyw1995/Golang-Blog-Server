package services

import (
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/repositories"
	"github.com/track/blogserver/pkg/utils"
)

type LinkSerivce struct {
	repo *repositories.LinkRepository
}

func NewLinkSerivce() *LinkSerivce {
	return &LinkSerivce{repo: new(repositories.LinkRepository)}
}
func (ls *LinkSerivce) CheckUserExist(uid int) (models.Response, bool) {
	userRepo := new(repositories.UserRepository)
	if _, found := userRepo.GetByUserId(uid); !found {
		return models.Response{Err: common.ErrUserNoExist, Data: nil}, found
	} else {
		return models.Response{}, found
	}
}
func (ls *LinkSerivce) GetBy(userId, id int) models.Response {
	if resp, found := ls.CheckUserExist(userId); !found {
		return resp
	}
	if link, ok := ls.repo.GetById(userId, id); ok {
		return models.Response{Err: common.Err{Msg: common.MsgGetLinksSucc}, Data: link}
	} else {
		return models.Response{Err: common.ErrUserLinksNoExist, Data: nil}
	}
}

func (ls *LinkSerivce) GetAllBy(userId int) models.Response {
	if resp, found := ls.CheckUserExist(userId); !found {
		return resp
	}
	if links, ok := ls.repo.GetAllByUserId(userId); ok {
		return models.Response{Err: common.Err{Msg: common.MsgGetLinksSucc}, Data: links}
	} else {
		return models.Response{Err: common.ErrUserLinksNoExist, Data: links}
	}
}

func (ls *LinkSerivce) UpdateBy(userId int, id int, params map[string]interface{}) models.Response {
	if resp, found := ls.CheckUserExist(userId); !found {
		return resp
	}
	if link, found := ls.repo.GetById(userId, id); !found {
		return models.Response{Err: common.ErrUserLinksNoExist, Data: nil}
	} else {
		if err := ls.repo.Update(link, params); err != nil {
			return models.Response{Err: common.ErrInternal, Data: nil}
		} else {
			return models.Response{Err: common.Err{Msg: common.MsgUpdateLinksSucc}, Data: nil}
		}
	}
}

func (ls *LinkSerivce) Create(v utils.Validation, link *models.FriendlyLink) models.Response {
	if resp, found := ls.CheckUserExist(int(link.UserID)); !found {
		return resp
	}
	if err := ls.repo.Insert(link); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgCreateLinksSucc}, Data: nil}

	}
}
func (ls *LinkSerivce) DelUserLinkById(userId, id int) models.Response {
	if resp, found := ls.CheckUserExist(userId); !found {
		return resp
	}
	if link, found := ls.repo.GetById(userId, id); !found {
		return models.Response{Err: common.ErrUserLinksNoExist, Data: nil}
	} else {
		if err := ls.repo.DelByInstance(link); err != nil {
			return models.Response{Err: common.ErrInternal, Data: nil}
		} else {
			return models.Response{Err: common.Err{Msg: common.MsgDelLinksSucc}, Data: nil}
		}
	}
}
func (ls *LinkSerivce) DelAllByUserId(userId int) models.Response {
	if resp, found := ls.CheckUserExist(userId); !found {
		return resp
	}
	if err := ls.repo.DelAllByUid(userId); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgDelLinksSucc}, Data: nil}
	}
}
