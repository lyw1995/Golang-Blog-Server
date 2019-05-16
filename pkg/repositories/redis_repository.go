package repositories

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/track/blogserver/pkg/common"
	"github.com/track/blogserver/pkg/models"
	"github.com/track/blogserver/pkg/persistence"
	"time"
)

//Redis缓存管理, 文章之类的 不细分
type RedisRepository struct {
}

//获取Redis实例
func (rr *RedisRepository) redis() redis.Conn {
	return persistence.GetR()
}

//插入一篇文章到redis
func (rr *RedisRepository) InsertToRedis(article *models.Article) error {
	r := rr.redis()
	defer r.Close()
	//article::aid::uid 插入一篇文章
	key := fmt.Sprintf("%s::%d::%d", common.ArticleKey, article.ID, article.UserID)
	//article::uid 用户所有文章
	key1 := fmt.Sprintf("%s::%d", common.ArticleKey, article.UserID)
	fmt.Println(key, key1)
	_, err := r.Do("HMSET", key, "ctime", article.CreateTime.Unix(), "aid", article.ID, "pv", article.Views, "title", article.Title)
	_, err = r.Do("SADD", key1, key)
	return err
}

//更新redis文章标题 更改的 根据用户id和aid
func (rr *RedisRepository) UpdateTitleOfRedisByAid(article *models.Article) error {
	r := rr.redis()
	defer r.Close()
	key := fmt.Sprintf("%s::%d::%d", common.ArticleKey, article.ID, article.UserID)
	_, err := r.Do("HSET", key, "title", article.Title)
	return err
}

//获取文章访问量
func (rr *RedisRepository) GetPVByAid(uid, aid int) (int, error) {
	r := rr.redis()
	defer r.Close()
	key := fmt.Sprintf("%s::%d::%d", common.ArticleKey, aid, uid)
	return redis.Int(r.Do("HGET", key, "pv"))
}

//更新访问量
func (rr *RedisRepository) UpdatePv(uid, aid int, remoteIP string) {
	r := rr.redis()
	defer r.Close()
	key := fmt.Sprintf("%s::%d", remoteIP, aid)
	value := fmt.Sprintf("%s::%d::%d", common.ArticleKey, aid, uid)
	if reply, err := r.Do("TTL", key); err == nil {
		//不存在这个key获取已经过期
		if reply.(int64) < 0 {
			r.Do("SET", key, value, "EX", common.ExipreSecond)
			r.Do("HINCRBY", value, "pv", "1")
		}
	}
}

//获取发布文章个数
func (rr *RedisRepository) CountArticls(uid int) (int, error) {
	r := rr.redis()
	defer r.Close()
	//article::uid 用户所有文章
	key := fmt.Sprintf("%s::%d", common.ArticleKey, uid)
	return redis.Int(r.Do("SCARD", key, ))
}

//删除redis缓存文章
func (rr *RedisRepository) DelOfRedis(uid, aid int) {
	r := rr.redis()
	defer r.Close()
	//先删除文章hash 在从用户article里面删除关联hash
	key := fmt.Sprintf("%s::%d::%d", common.ArticleKey, aid, uid)
	r.Do("DEL", key)
	key1 := fmt.Sprintf("%s::%d", common.ArticleKey, uid)
	r.Do("SREM", key1, key)
}

//获取总独立ip个数
func (rr *RedisRepository) CountIps() (int, error) {
	r := rr.redis()
	defer r.Close()
	return redis.Int(r.Do("SCARD", common.IPKey))
}

//插入日活ip,独立ip
func (rr *RedisRepository) InsertIp(ip string) {
	if len(ip) <= 0 {
		return
	}
	r := rr.redis()
	defer r.Close()
	now := time.Now()
	date := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
	key := fmt.Sprintf("%s::%s", common.IPKey, date)
	//独立
	r.Do("SADD", common.IPKey, ip)
	//日活
	r.Do("SADD", key, ip)
}

//获取日活ip
func (rr *RedisRepository) CountUV() (int, error) {
	r := rr.redis()
	defer r.Close()
	now := time.Now()
	date := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
	key := fmt.Sprintf("%s::%s", common.IPKey, date)
	return redis.Int(r.Do("SCARD", key))
}

//获取某个用户所有文章访问量
func (rr *RedisRepository) TotalPV(uid int) (int, error) {
	r := rr.redis()
	defer r.Close()
	key := fmt.Sprintf("%s::%d", common.ArticleKey, uid)
	keys, err := redis.Strings(r.Do("SMEMBERS", key))
	total := 0
	for _, key := range keys {
		if v, err := redis.Int(r.Do("HGET", key, "pv")); err == nil {
			total += v
		}
	}
	return total, err
}

//解析文章
func (rr *RedisRepository) scanArticle(replys []interface{}, err error) []models.Article {
	articles := make([]models.Article, 0)
	if err != nil {
		return articles
	}
	argsLen := 4 //sort get个数aid title ctime pv
	for i := 0; i < len(replys)/argsLen; i++ {
		article := models.Article{}
		for index, item := range replys[(i * argsLen) : (i*argsLen)+argsLen] {
			if item == nil {
				continue
			}
			switch index {
			case 0:
				if id, e := redis.Int(item, err); e == nil {
					article.ID = uint(id)
				}
			case 1:
				if title, e := redis.String(item, err); e == nil {
					article.Title = title
				}
			case 2:
				if ctime, e := redis.Int64(item, err); e == nil {
					article.CreateTime = time.Unix(ctime, 0)
				}
			case 3:
				if pv, e := redis.Int(item, err); e == nil {
					article.Views = pv
				}

			}
		}
		if article.IsValid() {
			articles = append(articles, article)
		}
	}
	return articles
}

//获取热门文章与最新文章(前五篇)
func (rr *RedisRepository) SelectHotAndNewArticle(uid int) map[string][]models.Article {
	r := rr.redis()
	defer r.Close()
	resp := make(map[string][]models.Article)
	key := fmt.Sprintf("%s::%d", common.ArticleKey, uid)
	resp["hot_articles"] = rr.scanArticle(redis.Values(r.Do("SORT", key, "BY", "*->pv", "DESC", "LIMIT", "0", common.LenLimit/2,
		"GET", "*->aid", "GET", "*->title", "GET", "*->ctime", "GET", "*->pv")))
	resp["new_articles"] = rr.scanArticle(redis.Values(r.Do("SORT", key, "BY", "*->ctime", "DESC", "LIMIT", "0", common.LenLimit/2,
		"GET", "*->aid", "GET", "*->title", "GET", "*->ctime", "GET", "*->pv")))
	return resp
}
