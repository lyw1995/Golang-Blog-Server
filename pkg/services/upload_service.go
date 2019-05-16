package services

import (
	"github.com/devfeel/dotweb"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/repositories"
)

type UploadService struct {
	repo *repositories.UploadRepository
}

func NewUploadService() *UploadService {
	return &UploadService{repo: repositories.NewUploadRepository()}
}

func (us *UploadService) SaveAvator(file *dotweb.UploadFile) models.Response {
	if file.GetFileExt() != ".jpg" && file.GetFileExt() != ".png" && file.GetFileExt() != ".jpeg" {
		return models.Response{Err: common.ErrUploadExtNotAllow, Data: nil}
	}
	if filepath, err := us.repo.Save(file.RandomFileName(), file.File); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgSaveImageSucc}, Data: filepath}
	}
}
func (us *UploadService) SaveCover(file *dotweb.UploadFile) models.Response {
	if file.GetFileExt() != ".jpg" && file.GetFileExt() != ".png" && file.GetFileExt() != ".jpeg" {
		return models.Response{Err: common.ErrUploadExtNotAllow, Data: nil}
	}
	if filepath, err := us.repo.Save(file.RandomFileName(), file.File); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgSaveImageSucc}, Data: filepath}
	}
}
