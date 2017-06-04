package model

import "github.com/toolkits/cache"

var PageRepo *Page = new(Page)

func (this *Page) LoadByName(name string) (*Page, error) {
	obj := new(Page)
	key := CacheKeyPage(name)
	if cache.Get(key, obj) == nil {
		return obj, nil
	}

	has, err := Orm.Where("name=?", name).Get(obj)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	go cache.Set(key, *obj, cache.DEFAULT)
	return obj, nil
}

func (this *Page) Save(name, content string) (int64, error) {
	obj, err := this.LoadByName(name)
	if err != nil {
		return 0, err
	}

	if obj == nil {
		obj = new(Page)
		obj.Name = name
		obj.Content = content
		return Orm.Insert(obj)
	} else {
		return obj.Update(content)
	}
}
