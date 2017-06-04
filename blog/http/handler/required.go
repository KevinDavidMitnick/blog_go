package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/toolkits/web/errors"
	"github.com/ulricqin/blog/model"
)

func ArticleRequired(w http.ResponseWriter, r *http.Request) *model.Article {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		panic(errors.BadRequestError())
	}

	obj, err := model.ArticleRepo.LoadById(id)
	errors.MaybePanic(err)

	if obj == nil {
		panic(errors.NotFoundError())
	}

	return obj
}

func ArticleUrlRequired(w http.ResponseWriter, r *http.Request) *model.Article {
	vars := mux.Vars(r)
	obj, err := model.ArticleRepo.LoadByUrl(vars["url"])
	errors.MaybePanic(err)

	if obj == nil {
		panic(errors.NotFoundError())
	}

	return obj
}
