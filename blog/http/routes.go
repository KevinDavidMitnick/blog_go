package http

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/ulricqin/blog/config"
	"github.com/ulricqin/blog/http/handler"
	"github.com/ulricqin/blog/http/middleware"
)

func ConfigRouter(r *mux.Router) {
	configFrontendRoutes(r)
	configBackendRoutes(r)
	configAuthRoutes(r)
	configCommonRoutes(r)
	configStaticRoutes(r)
}

func configFrontendRoutes(r *mux.Router) {
	r.HandleFunc("/", handler.HomeIndex).Methods("GET")
	r.HandleFunc("/archive", handler.ArchivePosts).Methods("GET")
	r.HandleFunc("/tag/{tag}", handler.TagPosts).Methods("GET")
	r.HandleFunc("/about", handler.AboutPage).Methods("GET")
	r.HandleFunc("/article/{url}.html", handler.ArticleGet).Methods("GET")
}

func configBackendRoutes(r *mux.Router) {
	backendBase := mux.NewRouter()
	r.PathPrefix("/backend").Handler(negroni.New(
		negroni.HandlerFunc(middleware.Authentication),
		negroni.Wrap(backendBase),
	))

	backend := backendBase.PathPrefix("/backend").Subrouter()

	backend.HandleFunc("/", handler.ArticleCreator).Methods("GET")
	backend.HandleFunc("/articles", handler.ArticlesPost).Methods("POST")
	backend.HandleFunc("/draft/{url}.html", handler.DraftGet).Methods("GET")
	backend.HandleFunc("/article/{id:[0-9]+}", handler.ArticleDelete).Methods("DELETE")
	backend.HandleFunc("/article/{id:[0-9]+}/editor", handler.ArticleEditor).Methods("GET")
	backend.HandleFunc("/drafts", handler.DraftsGet).Methods("GET")
	backend.HandleFunc("/about", handler.AboutGet).Methods("GET")
	backend.HandleFunc("/about", handler.AboutPut).Methods("PUT")
	backend.HandleFunc("/upload", handler.Upload).Methods("POST")
}

func configAuthRoutes(r *mux.Router) {
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", handler.LoginPage).Methods("GET")
	auth.HandleFunc("/login", handler.Login).Methods("POST")
	auth.HandleFunc("/logout", handler.Logout).Methods("GET")
}

func configStaticRoutes(r *mux.Router) {
	r.HandleFunc("/favicon.ico", handler.Favicon)
	r.PathPrefix("/img").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/css").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/js").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/" + config.G.Upload).Handler(http.FileServer(http.Dir("./static")))
}

func configCommonRoutes(r *mux.Router) {
	r.HandleFunc("/health", handler.Health).Methods("GET")
	r.HandleFunc("/version", handler.Version).Methods("GET")
}
