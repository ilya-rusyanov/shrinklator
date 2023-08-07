package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ilya-rusyanov/shrinklator/internal/config"
	"github.com/ilya-rusyanov/shrinklator/internal/handlers"
	"github.com/ilya-rusyanov/shrinklator/internal/server"
	"github.com/ilya-rusyanov/shrinklator/internal/services"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"
)

func newRouter(expandHandler http.HandlerFunc, shortenHandler http.HandlerFunc) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", shortenHandler)
	r.Get("/{id}", expandHandler)
	return r
}

func main() {
	config := config.New()
	config.Parse()

	storage := storage.NewInMemory()

	shortenerService := services.NewShortener(storage)

	shortenHandler := handlers.Shorten(shortenerService, config.BasePath)
	expandHandler := handlers.Expand(shortenerService)

	router := newRouter(shortenHandler, expandHandler)

	err := server.Run(config.ListenAddr, router)
	if err != nil {
		panic(err)
	}
}
