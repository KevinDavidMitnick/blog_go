package handler

import (
	"net/http"
	"strings"

	"github.com/ulricqin/blog/http/render"
	"github.com/ulricqin/blog/model"

	"github.com/toolkits/str"
	"github.com/toolkits/web/errors"
	"github.com/toolkits/web/param"
)

func ArticleCreator(w http.ResponseWriter, r *http.Request) {
	render.Put(r, "Title", "博文撰写")
	render.HTML(r, w, "backend/editor")
}

func ArticlesPost(w http.ResponseWriter, r *http.Request) {
	id := param.Int64(r, "id", 0)
	title := param.MustString(r, "title")
	url := param.MustString(r, "url")
	content := param.MustString(r, "content")
	tags := param.String(r, "tags", "")
	status := param.MustInt(r, "status")

	if !str.IsMatch(url, `^[a-z0-9\-]+$`) {
		errors.Panic("URL格式不合法，只允许小写字母、数字、中划线")
	}

	// 防止使用了中文逗号
	tags = strings.Replace(tags, "，", ",", -1)

	// ,linux,go, for sql like
	if tags != "" {
		tags = strings.TrimLeft(tags, ",")
		tags = strings.TrimRight(tags, ",")
		tags = "," + tags + ","
	}

	if id == 0 {
		_, err := model.ArticleRepo.Create(title, url, tags, content, status)
		errors.MaybePanic(err)
	} else {
		a, err := model.ArticleRepo.LoadById(id)
		errors.MaybePanic(err)
		_, err = a.Update(title, url, content, tags, status)
		errors.MaybePanic(err)
	}

	render.Message(w, "")
}

func DraftsGet(w http.ResponseWriter, r *http.Request) {
	drafts, err := model.ArticleRepo.Drafts()
	errors.MaybePanic(err)

	render.Put(r, "Title", "草稿箱")
	render.Put(r, "Drafts", drafts)
	render.HTML(r, w, "backend/drafts")
}

func DraftGet(w http.ResponseWriter, r *http.Request) {
	a := ArticleUrlRequired(w, r)
	render.Put(r, "Title", a.Url)
	render.Put(r, "Draft", a)
	render.HTML(r, w, "backend/draft")
}

func ArticleDelete(w http.ResponseWriter, r *http.Request) {
	_, err := ArticleRequired(w, r).Delete()
	render.Error(w, err)
}

func ArticleEditor(w http.ResponseWriter, r *http.Request) {
	render.Put(r, "Title", "编辑博文")
	render.Put(r, "Article", ArticleRequired(w, r))
	render.HTML(r, w, "backend/editor")
}

func AboutGet(w http.ResponseWriter, r *http.Request) {
	page, err := model.PageRepo.LoadByName("about")
	errors.MaybePanic(err)

	// 防止渲染一个nil给html造成crash
	if page == nil {
		page = new(model.Page)
	}

	render.Put(r, "Title", "自我介绍")
	render.Put(r, "Page", page)
	render.HTML(r, w, "backend/about")
}

func AboutPut(w http.ResponseWriter, r *http.Request) {
	_, err := model.PageRepo.Save("about", param.MustString(r, "content"))
	render.Error(w, err)
}
