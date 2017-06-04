package handler

import (
	"fmt"
	"net/http"

	"github.com/ulricqin/blog/config"
)

func Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func Version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, config.VERSION)
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
}
