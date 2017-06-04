package middleware

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/toolkits/web/errors"
	"github.com/ulricqin/blog/http/cookie"
)

func Authentication(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user := cookie.ReadUser(r)
	if user == "" {
		panic(errors.NotLoginError())
	}

	context.Set(r, "User", user)
	next(w, r)
}
