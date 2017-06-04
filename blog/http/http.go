package http

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/toolkits/file"
	"github.com/ulricqin/blog/config"
	"github.com/ulricqin/blog/http/cookie"
	"github.com/ulricqin/blog/http/middleware"
	"github.com/ulricqin/blog/http/render"
)

func Start() {
	render.Init()
	cookie.Init()

	r := mux.NewRouter().StrictSlash(false)
	ConfigRouter(r)

	n := negroni.New()
	n.Use(middleware.NewRecovery(file.MustOpenLogFile(config.G.Log.Error)))
	n.Use(middleware.NewLogger(file.MustOpenLogFile(config.G.Log.Access)))

	n.UseHandler(r)

	log.Println("listening on", config.G.Http.Listen)
	log.Fatal(http.ListenAndServe(config.G.Http.Listen, n))
}
