package render

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/context"
	"github.com/ulricqin/blog/config"
	"github.com/ulricqin/blog/http/cookie"
	"github.com/unrolled/render"
)

var Render *render.Render
var fm = template.FuncMap{
	"safe": func(raw string) template.HTML {
		return template.HTML(raw)
	},
}

func Init() {
	debug := config.G.Debug
	Render = render.New(render.Options{
		Directory:     "views",
		Extensions:    []string{".html"},
		Delims:        render.Delims{"{{", "}}"},
		IndentJSON:    false,
		Funcs:         []template.FuncMap{fm},
		IsDevelopment: debug,
	})
}

func Put(r *http.Request, key string, val interface{}) {
	m, ok := context.GetOk(r, "DATA_MAP")
	if ok {
		mm := m.(map[string]interface{})
		mm[key] = val
		context.Set(r, "DATA_MAP", mm)
	} else {
		context.Set(r, "DATA_MAP", map[string]interface{}{key: val})
	}
}

func HTML(r *http.Request, w http.ResponseWriter, name string, htmlOpt ...render.HTMLOptions) {
	username := cookie.ReadUser(r)
	Put(r, "Username", config.G.Username)
	Put(r, "Login", username == config.G.Username)
	Put(r, "Portrait", config.G.Portrait)
	Put(r, "Motto", config.G.Motto)
	Put(r, "Contact", config.G.Contact)
	Render.HTML(w, http.StatusOK, name, context.Get(r, "DATA_MAP"), htmlOpt...)
}

func Error(w http.ResponseWriter, err error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}

	Render.JSON(w, http.StatusOK, map[string]string{"msg": msg})
}

func Message(w http.ResponseWriter, format string, args ...interface{}) {
	Render.JSON(w, http.StatusOK, map[string]string{"msg": fmt.Sprintf(format, args...)})
}

func Data(w http.ResponseWriter, v interface{}, err ...error) {
	m := ""
	if len(err) > 0 && err[0] != nil {
		m = err[0].Error()
	}

	Render.JSON(w, http.StatusOK, map[string]interface{}{"msg": m, "data": v})
}

func Text(w http.ResponseWriter, v string) {
	Render.Text(w, http.StatusOK, v)
}
