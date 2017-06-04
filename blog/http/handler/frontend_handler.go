package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/toolkits/web"
	"github.com/toolkits/web/errors"
	"github.com/toolkits/web/param"
	"github.com/ulricqin/blog/config"
	"github.com/ulricqin/blog/http/render"
	"github.com/ulricqin/blog/model"
	"github.com/ulricqin/blog/utils"
)

func HomeIndex(w http.ResponseWriter, r *http.Request) {
	limit := param.Int(r, "limit", 20)
	total, err := model.ArticleRepo.Total()
	errors.MaybePanic(err)

	pager := web.NewPaginator(r, limit, total)
	ids, err := model.ArticleRepo.Ids(limit, pager.Offset())
	errors.MaybePanic(err)

	list, err := model.ArticleRepo.LoadByIds(ids)
	errors.MaybePanic(err)

	render.Put(r, "List", list)
	render.Put(r, "Pager", pager)
	render.Put(r, "Title", fmt.Sprintf("%s's Blog", config.G.Username))
	render.HTML(r, w, "frontend/index")
}

func TagPosts(w http.ResponseWriter, r *http.Request) {
	tag := mux.Vars(r)["tag"]
	limit := param.Int(r, "limit", 20)
	total, err := model.ArticleRepo.TotalByTag(tag)
	errors.MaybePanic(err)

	pager := web.NewPaginator(r, limit, total)
	ids, err := model.ArticleRepo.IdsByTag(tag, limit, pager.Offset())
	errors.MaybePanic(err)

	list, err := model.ArticleRepo.LoadByIds(ids)
	errors.MaybePanic(err)

	render.Put(r, "List", list)
	render.Put(r, "Pager", pager)
	render.Put(r, "Tag", tag)
	render.Put(r, "Title", tag)
	render.HTML(r, w, "frontend/tag")
}

func ArchivePosts(w http.ResponseWriter, r *http.Request) {
	limit := 1000
	total, err := model.ArticleRepo.Total()
	errors.MaybePanic(err)

	pager := web.NewPaginator(r, limit, total)
	ids, err := model.ArticleRepo.Ids(limit, pager.Offset())
	errors.MaybePanic(err)

	list, err := model.ArticleRepo.LoadByIds(ids)
	errors.MaybePanic(err)

	render.Put(r, "List", list)
	render.Put(r, "Pager", pager)
	render.Put(r, "Title", "Archive Posts")
	render.HTML(r, w, "frontend/archive")
}

func ArticleGet(w http.ResponseWriter, r *http.Request) {
	a := ArticleUrlRequired(w, r)
	if a.Status == 0 {
		errors.Panic("it's a draft")
	}

	render.Put(r, "Title", a.Url)
	render.Put(r, "Article", a)
	render.HTML(r, w, "frontend/article")
}

func AboutPage(w http.ResponseWriter, r *http.Request) {
	page, err := model.PageRepo.LoadByName("about")
	errors.MaybePanic(err)

	// 防止渲染一个nil给html造成crash
	if page == nil {
		page = new(model.Page)
	}

	render.Put(r, "Content", utils.Markdown2HTML(page.Content))
	render.HTML(r, w, "frontend/about")
}
