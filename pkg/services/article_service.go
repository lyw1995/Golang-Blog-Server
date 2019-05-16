package services

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/repositories"
	"sort"
	"strings"
	"sync"
	"time"
)

type ArticleService struct {
	repo  *repositories.ArticleRepository
	redis *repositories.RedisRepository
}

func NewArticleService() *ArticleService {
	return &ArticleService{repo: new(repositories.ArticleRepository), redis: new(repositories.RedisRepository)}
}
func (as *ArticleService) CheckUserExist(uid int) (models.Response, bool) {
	userRepo := new(repositories.UserRepository)
	if _, found := userRepo.GetByUserId(uid); !found {
		return models.Response{Err: common.ErrUserNoExist, Data: nil}, found
	} else {
		return models.Response{}, found
	}
}

//遍历查询复制文章浏览量
func (as *ArticleService) duplicatePv(wg *sync.WaitGroup, item *models.Article) {
	if pv, err := as.redis.GetPVByAid(int(item.UserID),int(item.ID)); err == nil {
		item.Views = pv
	}
	wg.Done()
}
func (as *ArticleService) parseSortValue(sort string) (models.SortValue, error) {
	var sortValues models.SortValue
	if strings.Contains(sort, ",") {
		values := strings.Split(sort, ",")
		if len(values) > 2 {
			return models.SortValue{}, common.ErrClientParams
		}
		sortValues = models.SortValue{Key: values[0], Value: values[1]}
		if !(sortValues.IsValidKey() && sortValues.IsValidValue()) {
			return models.SortValue{}, common.ErrClientParams
		}
	}
	return sortValues, nil
}
func (as *ArticleService) GetPageBy(uid, page int, sortValue string) models.Response {
	if resp, found := as.CheckUserExist(uid); !found {
		return resp
	}
	if sortt, err := as.parseSortValue(sortValue); err != nil {
		return models.Response{Err: err.(common.Err), Data: nil}
	} else {
		resp := make(map[string]interface{})
		articles, _ := as.repo.GetArticlesWithState(uid, page, 0, sortt)
		resp["articles"] = articles
		resp["count"] = as.repo.GetCountByUserId(uid)
		var wg sync.WaitGroup
		for index := range articles {
			wg.Add(1)
			go as.duplicatePv(&wg, &articles[index])
		}
		wg.Wait()
		if sortt.Key == "views" {
			sort.Slice(articles, func(i, j int) bool {
				if sortt.Value == "desc" {
					return articles[i].Views > articles[j].Views
				} else {
					return articles[i].Views < articles[j].Views
				}
			})
		}
		return models.Response{Err: common.Err{Msg: common.MsgGetArticleSucc}, Data: resp}
	}
}

func (as *ArticleService) GetDraftPageBy(uid, page int, sort string) models.Response {
	if resp, found := as.CheckUserExist(uid); !found {
		return resp
	}
	if sort, err := as.parseSortValue(sort); err != nil {
		return models.Response{Err: err.(common.Err), Data: nil}
	} else {
		resp := make(map[string]interface{})
		articles, _ := as.repo.GetArticlesWithState(uid, page, 2, sort)
		resp["articles"] = articles
		resp["count"] = as.repo.GetDraftCountByUserId(uid)
		return models.Response{Err: common.Err{Msg: common.MsgGetArticleSucc}, Data: resp}
	}
}

func (as *ArticleService) GetCategoryItemArticlesPageBy(uid, cid, page int, sortValue string) models.Response {
	//判断用户跟分类是否存在
	if resp, found := as.CheckUserExist(uid); !found {
		return resp
	}
	categoryRepo := new(repositories.CateGoryRepository)
	if categoryItem, found := categoryRepo.GetCategoryItemById(cid); !found {
		return models.Response{Err: common.ErrCategoryNoExist, Data: nil}
	} else {
		if category, found := categoryRepo.GetCategoryById(categoryItem.CategoryID); !found {
			return models.Response{Err: common.ErrCategoryNoExist, Data: nil}
		} else {
			if sortt, err := as.parseSortValue(sortValue); err != nil {
				return models.Response{Err: err.(common.Err), Data: nil}
			} else {
				resp := make(map[string]interface{})
				if category.Name == "个人分类" {
					articles, _ := as.repo.GetCategoryItemArticlesPage(uid, cid, page, sortt)
					resp["articles"] = articles
					resp["count"] = as.repo.GetCountByCid(uid,cid)
				} else {
					articles, _ := as.repo.GetArchizeArticlesPage(uid, page, categoryItem.Name, sortt)
					resp["articles"] = articles
					resp["count"] = as.repo.GetCountByCname(uid,categoryItem.Name)
				}
				articles := resp["articles"].([]models.Article)
				var wg sync.WaitGroup
				for index := range articles {
					wg.Add(1)
					go as.duplicatePv(&wg, &articles[index])
				}
				wg.Wait()
				if sortt.Key == "views" {
					sort.Slice(articles, func(i, j int) bool {
						if sortt.Value == "desc" {
							return articles[i].Views > articles[j].Views
						} else {
							return articles[i].Views < articles[j].Views
						}
					})
				}
				return models.Response{Err: common.Err{Msg: common.MsgGetArticleSucc}, Data: resp}
			}
		}
	}

}

func (as *ArticleService) Create(article *models.Article) models.Response {
	//判断用户跟分类是否存在
	if resp, found := as.CheckUserExist(int(article.UserID)); !found {
		return resp
	}
	categoryRepo := new(repositories.CateGoryRepository)
	if _, found := categoryRepo.GetCategoryItemById(int(article.CategoryItemID)); !found {
		return models.Response{Err: common.ErrCategoryNoExist, Data: nil}
	}
	as.FindCover(article)
	if err := as.repo.Insert(article); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		//保存草稿
		if article.State == 2 {
			return models.Response{Err: common.Err{Msg: common.MsgSaveArticleSucc}, Data: nil}
		} else {
			//直接发布, 查询归档并创建,插入redis
			as.UpdateOrInsertArchive(categoryRepo, article)
			as.redis.InsertToRedis(article)
			return models.Response{Err: common.Err{Msg: common.MsgCreateArticleSucc}, Data: nil}
		}
	}
	return models.Response{Err: common.Err{Msg: common.MsgCreateArticleSucc}, Data: nil}
}

//纳入归档
func (as *ArticleService) UpdateOrInsertArchive(cr *repositories.CateGoryRepository, article *models.Article) error {
	categoryName := fmt.Sprintf("%d年%02d月", article.CreateTime.Year(), article.CreateTime.Month())
	_, err := cr.GetFirstOnCreateArchive(article.UserID, categoryName)
	return err
}

//正则取Content的 img 标签第一张作为cover
func (as *ArticleService) FindCover(article *models.Article) {
	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Content)); err != nil {
		return
	} else {
		doc.Find("img").EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if value, ok := selection.Attr("src"); ok {
				article.Cover = value
				return false
			}
			return true
		})
	}
}

//根据状态修改文章
func (as *ArticleService) UpdateWithState(article *models.Article, state string) models.Response {
	if resp, found := as.CheckUserExist(int(article.UserID)); !found {
		return resp
	}
	categoryRepo := new(repositories.CateGoryRepository)
	if _, found := categoryRepo.GetCategoryItemById(int(article.CategoryItemID)); !found {
		return models.Response{Err: common.ErrCategoryNoExist, Data: nil}
	}
	//先去查文章
	if tArticle, found := as.repo.GetByUidAndAid(int(article.UserID), int(article.ID)); !found {
		return models.Response{Err: common.ErrArticleNoExist, Data: nil}
	} else {
		// 已经发布 保存修改
		if tArticle.State == 0 {
			article.State = tArticle.State
			as.FindCover(article)
			if err := as.repo.Save(article); err != nil {
				return models.Response{Err: common.ErrInternal, Data: nil}
			}
			as.redis.UpdateTitleOfRedisByAid(article)
			return models.Response{Err: common.Err{Msg: common.MsgUpdateArticleSucc}, Data: nil}
		} else {
			//草稿箱中,在根据状态操作
			if state == "publish" {
				article.State = 0
				article.CreateTime = time.Now()
			} else {
				article.State = 2 //继续保存草稿
			}
			as.FindCover(article)
			if err := as.repo.Save(article); err != nil {
				return models.Response{Err: common.ErrInternal, Data: nil}
			}
			if article.State == 0 {
				as.UpdateOrInsertArchive(categoryRepo, article)
				as.redis.InsertToRedis(article)
			}
			return models.Response{Err: common.Err{Msg: common.MsgUpdateArticleSucc}, Data: nil}
		}

	}
}
func (as *ArticleService) GetArticleByUidAndCIdAndAid(uid, cid, aid int) models.Response {
	if resp, found := as.CheckUserExist(int(uid)); !found {
		return resp
	}
	categoryRepo := new(repositories.CateGoryRepository)
	if _, found := categoryRepo.GetCategoryItemById(int(cid)); !found {
		return models.Response{Err: common.ErrCategoryNoExist, Data: nil}
	}

	if article, found := as.repo.GetByUidAndCidAndAidWithState(uid, cid, aid, 0); !found {
		return models.Response{Err: common.ErrArticleNoExist, Data: nil}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgGetArticleSucc}, Data: *article}
	}
}
func (as *ArticleService) GetArticleByUidAndAid(uid, aid int, remoteIP string) models.Response {
	if resp, found := as.CheckUserExist(int(uid)); !found {
		return resp
	}

	if article, found := as.repo.GetByUidAndAidWithState(uid, aid, 0); !found {
		return models.Response{Err: common.ErrArticleNoExist, Data: nil}
	} else {
		as.redis.UpdatePv(uid,aid, remoteIP)
		if pv, err := as.redis.GetPVByAid(uid,int(article.ID)); err == nil {
			article.Views = pv
		}
		return models.Response{Err: common.Err{Msg: common.MsgGetArticleSucc}, Data: *article}
	}
}
func (as *ArticleService) GetArticleWithState(uid, aid int, op string) models.Response {
	if resp, found := as.CheckUserExist(int(uid)); !found {
		return resp
	}
	state := -1

	if op == "publish" {
		state = 0
	} else if op == "draft" {
		state = 2
	}

	if article, found := as.repo.GetByUidAndAidWithState(uid, aid, state); !found {
		return models.Response{Err: common.ErrArticleNoExist, Data: nil}
	} else {
		return models.Response{Err: common.Err{Msg: common.MsgGetArticleSucc}, Data: *article}
	}
}
func (as *ArticleService) DelArticlesByUidAndCid(uid, cid int) models.Response {
	if resp, found := as.CheckUserExist(int(uid)); !found {
		return resp
	}
	categoryRepo := new(repositories.CateGoryRepository)
	if _, found := categoryRepo.GetCategoryItemById(int(cid)); !found {
		return models.Response{Err: common.ErrCategoryNoExist, Data: nil}
	}
	if err := as.repo.DelAllByUidAndCid(uid, cid); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	}
	return models.Response{Err: common.Err{Msg: common.MsgDelArticleSucc}, Data: nil}
}
func (as *ArticleService) DelArticleByUidAndCidAndAid(uid, cid, aid int, op string) models.Response {
	if resp, found := as.CheckUserExist(int(uid)); !found {
		return resp
	}
	categoryRepo := new(repositories.CateGoryRepository)
	if _, found := categoryRepo.GetCategoryItemById(int(cid)); !found {
		return models.Response{Err: common.ErrCategoryNoExist, Data: nil}
	}
	state := -1
	if op == "publish" {
		state = 0
	} else {
		state = 2
	}
	if article, found := as.repo.GetByUidAndCidAndAidWithState(uid, cid, aid, state); !found {
		return models.Response{Err: common.ErrArticleNoExist, Data: nil}
	} else {
		if err := as.repo.DelFor(article); err != nil {
			return models.Response{Err: common.ErrInternal, Data: nil}
		} else {
			//删除redis缓存
			if article.State == 0 {
				as.redis.DelOfRedis(uid,aid)
			}
		}
		return models.Response{Err: common.Err{Msg: common.MsgDelArticleSucc}, Data: nil}
	}
}
func (as *ArticleService) UpdateArticle(uid, cid, aid int, params map[string]interface{}) models.Response {
	if resp, found := as.CheckUserExist(int(uid)); !found {
		return resp
	}
	categoryRepo := new(repositories.CateGoryRepository)
	if _, found := categoryRepo.GetCategoryItemById(int(cid)); !found {
		return models.Response{Err: common.ErrCategoryNoExist, Data: nil}
	}
	if article, found := as.repo.GetByUidAndCidAndAidWithState(uid, cid, aid, 0); !found {
		return models.Response{Err: common.ErrArticleNoExist, Data: nil}
	} else {
		if err := as.repo.Update(article, params); err != nil {
			return models.Response{Err: common.ErrInternal, Data: nil}
		}
		return models.Response{Err: common.Err{Msg: common.MsgUpdateArticleSucc}, Data: nil}
	}
}
