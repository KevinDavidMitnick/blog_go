package cookie

import (
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/ulricqin/blog/config"
)

const USER_COOKIE_NAME = "ulricqin-blog"

var SecureCookie *securecookie.SecureCookie

func Init() {
	var hashKey = []byte(config.G.Http.Secret)
	var blockKey = []byte(nil)
	SecureCookie = securecookie.New(hashKey, blockKey)
}

func ReadUser(r *http.Request) string {
	cookie, err := r.Cookie(USER_COOKIE_NAME)
	if err != nil {
		// not found in cookie
		return ""
	}

	var name string
	err = SecureCookie.Decode(USER_COOKIE_NAME, cookie.Value, &name)
	if err != nil {
		// fail to decode secure cookie
		return ""
	}

	return name
}

func WriteUser(w http.ResponseWriter, name string) {
	encoded, _ := SecureCookie.Encode(USER_COOKIE_NAME, name)
	c := &http.Cookie{
		Name:     USER_COOKIE_NAME,
		Value:    encoded,
		Path:     "/",
		MaxAge:   3600 * 24 * 365,
		HttpOnly: true,
	}

	http.SetCookie(w, c)
}

func RemoveUser(w http.ResponseWriter) {
	encoded, _ := SecureCookie.Encode(USER_COOKIE_NAME, "")
	cookie := &http.Cookie{
		Name:     USER_COOKIE_NAME,
		Value:    encoded,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}
