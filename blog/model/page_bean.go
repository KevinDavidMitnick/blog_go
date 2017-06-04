package model

import "github.com/toolkits/cache"

type Page struct {
	Name    string
	Content string
}

func (this *Page) Update(content string) (int64, error) {
	this.Content = content

	affected, err := Orm.Where("name=?", this.Name).Cols("content").Update(this)
	if affected > 0 {
		go cache.Delete(CacheKeyPage(this.Name))
	}

	return affected, err
}
