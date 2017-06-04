package model

import (
	"fmt"
	"time"

	"github.com/toolkits/cache"
)

var ArticleRepo *Article = new(Article)

func (this *Article) LoadById(id int64) (*Article, error) {
	obj := new(Article)
	key := CacheKeyArticle(id)
	if cache.Get(key, obj) == nil {
		return obj, nil
	}

	has, err := Orm.Id(id).Get(obj)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	go cache.Set(key, *obj, cache.DEFAULT)
	return obj, nil
}

func (this *Article) ID(url string) (int64, error) {
	var key string = CacheKeyArticleId(url)
	var id int64
	if cache.Get(key, &id) == nil {
		return id, nil
	}

	var obj Article
	has, err := Orm.Cols("id").Where("url=?", url).Get(&obj)
	if err != nil {
		return 0, err
	}

	if !has {
		return 0, nil
	}

	go cache.Set(key, obj.Id, cache.DEFAULT)
	return obj.Id, nil
}

func (this *Article) LoadByUrl(url string) (*Article, error) {
	id, err := this.ID(url)
	if err != nil {
		return nil, err
	}

	return this.LoadById(id)
}

func (this *Article) LoadByIds(ids []int64) ([]*Article, error) {
	count := len(ids)
	objs := make([]*Article, count)
	for i := 0; i < count; i++ {
		obj, err := this.LoadById(ids[i])
		if err != nil {
			return []*Article{}, err
		}
		objs[i] = obj
	}

	return objs, nil
}

func (this *Article) IdsByTag(tag string, limit, offset int) ([]int64, error) {
	objs := make([]*Article, 0)
	err := Orm.Cols("id").Where("tags like ? and status=1", "%,"+tag+",%").Desc("created").Limit(limit, offset).Find(&objs)
	if err != nil {
		return []int64{}, err
	}

	return ParseArticleIds(objs), nil
}

func (this *Article) Ids(limit, offset int) ([]int64, error) {
	objs := make([]*Article, 0)
	err := Orm.Cols("id").Where("status=1").Desc("created").Limit(limit, offset).Find(&objs)
	if err != nil {
		return []int64{}, err
	}

	return ParseArticleIds(objs), nil
}

func (this *Article) Drafts() ([]*Article, error) {
	objs := make([]*Article, 0)
	err := Orm.Where("status=0").Desc("created").Find(&objs)
	return objs, err
}

func ParseArticleIds(objs []*Article) []int64 {
	count := len(objs)
	ids := make([]int64, count)
	for i := 0; i < count; i++ {
		ids[i] = objs[i].Id
	}
	return ids
}

func (this *Article) Create(title, url, tags, content string, status int) (int64, error) {
	id, err := this.ID(url)
	if err != nil {
		return 0, err
	}

	if id > 0 {
		return 0, fmt.Errorf("url already existent")
	}

	obj := new(Article)
	obj.Title = title
	obj.Url = url
	obj.Tags = tags
	obj.Content = content
	obj.Status = status
	obj.Created = int(time.Now().Unix())

	_, err = Orm.Insert(obj)
	if err != nil {
		return 0, err
	}

	if obj.Status == 1 {
		go cache.Delete(CacheKeyArticleTotal())
	}

	return obj.Id, nil
}

func (this *Article) Total() (int64, error) {
	var total int64
	if cache.Get(CacheKeyArticleTotal(), &total) == nil {
		return total, nil
	}

	total, err := Orm.Where("status=1").Count(new(Article))
	if err != nil {
		return 0, err
	}

	go cache.Set(CacheKeyArticleTotal(), total, cache.DEFAULT)
	return total, nil
}

func (this *Article) TotalByTag(tag string) (int64, error) {
	return Orm.Where("tags like ? and status=1", "%,"+tag+",%").Count(new(Article))
}
