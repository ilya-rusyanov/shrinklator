package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ilya-rusyanov/shrinklator/internal/config"
	"github.com/ilya-rusyanov/shrinklator/internal/handlers"
	"github.com/ilya-rusyanov/shrinklator/internal/models"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"
)

func newRouter(expandHandler http.HandlerFunc, shortenHandler http.HandlerFunc) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", shortenHandler)
	r.Get("/{id}", expandHandler)
	return r
}

func Run() {
	sh := models.New(storage.New())

	router := newRouter(
		handlers.Shorten(sh, config.Values.BasePath),
		handlers.Expand(sh))

	err := http.ListenAndServe(config.Values.ListenAddr,
		router)

	if err != nil {
		panic(err)
	}
}
