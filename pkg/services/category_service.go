package services

import (
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/repositories"
	"sync"
)

type CategoryService struct {
	repo *repositories.CateGoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{repo: new(repositories.CateGoryRepository)}
}
func (cs *CategoryService) CheckUserExist(uid int) (models.Response, bool) {
	userRepo := new(repositories.UserRepository)
	if user, found := userRepo.GetByUserId(uid); !found {
		return models.Response{Err: common.ErrUserNoExist, Data: nil}, found
	} else {
		return models.Response{Data: user}, found
	}
}

//获取所有分类
func (cs *CategoryService) GetCategorys(uid int) models.Response {
	if resp, found := cs.CheckUserExist(uid); !found {
		return resp
	} else {
		//个人描述
		user := resp.Data.(*models.User)
		response := make(map[string]interface{})
		response["user_desc"] = user.UserInfo.UserDesc
		//分类列表
		categorys, _ := cs.repo.GetCategorysByUid(uid)
		var wg sync.WaitGroup
		ar := new(repositories.ArticleRepository)
		for index := range categorys {
			categoryItems := categorys[index].CategoryItem
			for index := range categoryItems {
				wg.Add(1)
				go cs.duplicateCategoryItemSize(ar, &wg, &categoryItems[index])
			}
			wg.Wait()
		}
		response["category_list"] = categorys
		es := new(repositories.RedisRepository)
		//其他信息
		if count, err := es.CountIps(); err != nil {
			return models.Response{Err: common.ErrInternal, Data: nil}
		} else {
			response["total_ip_count"] = count
		}
		if count, err := es.CountArticls(uid); err != nil {
			return models.Response{Err: common.ErrInternal, Data: nil}
		} else {
			response["articles_count"] = count
		}
		if count, err := es.TotalPV(uid); err != nil {
			return models.Response{Err: common.ErrInternal, Data: nil}
		} else {
			response["total_pv"] = count
		}
		return models.Response{Err: common.Err{Msg: common.MsgGetCategorySucc}, Data: response}
	}

	return models.Response{Err: common.ErrInternal, Data: nil}
}

//获取个人分类所有子条目(遍历数据库查询)
func (cs *CategoryService) GetPersonalCategorys(uid int) models.Response {
	if resp, found := cs.CheckUserExist(uid); !found {
		return resp
	}
	if temp, found := cs.repo.GetCategoryByName("个人分类"); found {
		var wg sync.WaitGroup
		ar := new(repositories.ArticleRepository)
		categoryItems, _ := cs.repo.GetCategoryItemsByCid(uint(uid), temp.ID)
		for index := range categoryItems {
			wg.Add(1)
			go cs.duplicateCategoryItemSize(ar, &wg, &categoryItems[index])
		}
		wg.Wait()
		return models.Response{Err: common.Err{Msg: common.MsgGetCategorySucc}, Data: categoryItems}
	} else {
		return models.Response{Err: common.ErrCategoryNoExist, Data: nil}
	}
}

//遍历查询复制分类条目大小
func (cs *CategoryService) duplicateCategoryItemSize(ar *repositories.ArticleRepository, wg *sync.WaitGroup, item *models.CategoryItem) {
	if temp, found := cs.repo.GetCategoryByName("个人分类"); found {
		if temp.ID == item.CategoryID {
			//个人分类去redis获取
			item.ItemSize = ar.GetCountByCid(int(item.UserID),int(item.ID))
		} else {
			//归档根据年月获取
			item.ItemSize = ar.GetCountByCname(int(item.UserID),item.Name)
		}
	}
	wg.Done()
}

//创建个人分类子条目
func (cs *CategoryService) CreateCategoryItem(category *models.CategoryItem) models.Response {
	if resp, found := cs.CheckUserExist(int(category.UserID)); !found {
		return resp
	}
	//存在并且是当前用户ID,则判断已经存在
	if _, found := cs.repo.GetCategoryItemByNameWithUid(category.Name,int(category.UserID)); found  {
		return models.Response{Err: common.ErrCategoryExist, Data: nil}
	}
	//创建都是个人分类, 归档文章自动生成的
	if temp, found := cs.repo.GetCategoryByName("个人分类"); found {
		category.CategoryID = temp.ID
	} else {
		return models.Response{Err: common.ErrCategoryNoExist, Data: nil}
	}

	if err := cs.repo.Insert(category); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgCreateCategorySucc}, Data: nil}
	}
}

//根据分类子条目ID 修改 ...
func (cs *CategoryService) UpdateCategoryItemBy(uid int, cid int, params map[string]interface{}) models.Response {
	if resp, found := cs.CheckUserExist(uid); !found {
		return resp
	}

	if categoryItem, found := cs.repo.GetCategoryItemById(cid); !found {
		return models.Response{Err: common.ErrCategoryNoExist, Data: nil}
	} else {
		//如果修改名称,则检查名称是否已存在
		if cname, ok := params["name"]; ok {
			if _, found := cs.repo.GetCategoryItemByNameWithUid(cname.(string),uid); found {
				return models.Response{Err: common.ErrCategoryExist, Data: nil}
			}
		}
		if err := cs.repo.Update(categoryItem, params); err != nil {
			return models.Response{Err: common.ErrInternal, Data: nil}
		} else {
			return models.Response{Err: common.Err{Msg: common.MsgUpdateCategorySucc}, Data: nil}
		}
	}
}

//软删除所有分类子条目
func (cs *CategoryService) DelAllCategoryItems(uid int) models.Response {
	if resp, found := cs.CheckUserExist(uid); !found {
		return resp
	}
	if err := cs.repo.DelCategoryItems(uid); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgDelCategorySucc}, Data: nil}
	}
}

//软删除某个分类子条目
func (cs *CategoryService) DelAllCategoryItemBy(uid, cid int) models.Response {
	if resp, found := cs.CheckUserExist(uid); !found {
		return resp
	}
	if _, found := cs.repo.GetCategoryItemById(cid); !found {
		return models.Response{Err: common.ErrCategoryNoExist, Data: nil}
	}
	if err := cs.repo.DelCategoryItemById(cid); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgDelCategorySucc}, Data: nil}
	}
}
