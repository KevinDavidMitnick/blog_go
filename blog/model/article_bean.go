package model

import (
	"strings"
	"time"

	"github.com/toolkits/cache"
	tt "github.com/toolkits/time"
	"github.com/ulricqin/blog/utils"
)

type Article struct {
	Id      int64
	Title   string
	Url     string
	Tags    string
	Content string
	Status  int
	Created int
}

func (this *Article) TagsReadable() string {
	if this.Tags == "" {
		return ""
	}

	tags := strings.TrimLeft(this.Tags, ",")
	tags = strings.TrimRight(tags, ",")

	return tags
}

func (this *Article) TagsArray() []string {
	tags := this.TagsReadable()
	if tags == "" {
		return []string{}
	}
	return strings.Split(tags, ",")
}

func (this *Article) CreatedReadable() string {
	now := int(time.Now().Unix())
	return tt.HumanDurationInt(now, this.Created)
}

func (this *Article) HtmlContent() string {
	return utils.Markdown2HTML(this.Content)
}

func (this *Article) HtmlSummary() string {
	index := strings.Index(this.Content, "<!-- more -->")
	if index > 0 {
		return utils.Markdown2HTML(this.Content[:index])
	} else {
		return utils.Markdown2HTML(this.Content)
	}
}

func (this *Article) Delete() (int64, error) {
	affected, err := Orm.Where("id=?", this.Id).Delete(new(Article))

	if affected > 0 {
		go cache.Delete(CacheKeyArticle(this.Id))
		go cache.Delete(CacheKeyArticleId(this.Url))
		if this.Status == 1 {
			go cache.Delete(CacheKeyArticleTotal())
		}
	}

	return affected, err
}

func (this *Article) Update(title, url, content, tags string, status int) (int64, error) {
	oldUrl := this.Url
	oldStatus := this.Status

	this.Title = title
	this.Url = url
	this.Content = content
	this.Tags = tags
	this.Status = status

	affected, err := Orm.Where("id=?", this.Id).Cols("title", "url", "content", "tags", "status").Update(this)

	if affected > 0 {
		go cache.Delete(CacheKeyArticle(this.Id))
		if oldUrl != url {
			go cache.Delete(CacheKeyArticleId(oldUrl))
		}
		if oldStatus != status {
			go cache.Delete(CacheKeyArticleTotal())
		}
	}

	return affected, err
}
