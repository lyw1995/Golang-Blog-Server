package services

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/repositories"
	"time"
)

type ExtendSerivce struct {
	redis *repositories.RedisRepository
}

func NewExtendSerivce() *ExtendSerivce {
	return &ExtendSerivce{redis: new(repositories.RedisRepository)}
}

//其实就只有日活IP 总独立IP, 总文章数,总访问量
func (es *ExtendSerivce) GetAll(uid int) models.Response {
	now := time.Now()
	data := make(map[string]interface{})
	data["date"] = fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
	if count, err := es.redis.CountIps(); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		data["total_ip_count"] = count
	}
	if count, err := es.redis.CountUV(); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		data["today_ip_count"] = count
	}
	if count, err := es.redis.CountArticls(uid); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		data["articles_count"] = count
	}
	if count, err := es.redis.TotalPV(uid); err != nil {
		return models.Response{Err: common.ErrInternal, Data: nil}
	} else {
		data["total_pv"] = count
	}

	return models.Response{Err: common.Err{Msg: common.MsgGetExtSucc}, Data: data}
}
func (es *ExtendSerivce) CollectByCsdn(url string, article *models.Article) models.Response {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return models.Response{Err: common.ErrCollectSource, Data: nil}
	}
	doc.Find(".blog-content-box").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		article.Title = s.Find(".title-article").Text()
		fmt.Println(len(article.Title))
		s.Find("article").Each(func(i int, selection *goquery.Selection) {
			selection.Find(".hide-article-box").Remove()
			if html, err := selection.Html(); err == nil {
				article.Content = html
			}
		})
	})
	article.CreateTime = time.Now()
	as := NewArticleService()
	return as.Create(article)
}
