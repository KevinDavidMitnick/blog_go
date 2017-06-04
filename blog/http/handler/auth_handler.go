package handler

import (
	"net/http"

	"github.com/toolkits/web/param"
	"github.com/ulricqin/blog/config"
	"github.com/ulricqin/blog/http/cookie"
	"github.com/ulricqin/blog/http/render"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	render.Put(r, "callback", param.String(r, "callback", "/"))
	render.HTML(r, w, "common/login")
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := param.MustString(r, "username")
	password := param.MustString(r, "password")

	message := ""
	if username != config.G.Username || password != config.G.Password {
		message = "username or password error"
	} else {
		cookie.WriteUser(w, username)
	}

	render.Message(w, message)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie.RemoveUser(w)
	http.Redirect(w, r, "/", 302)
}
