package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ilya-rusyanov/shrinklator/internal/config"
	"github.com/ilya-rusyanov/shrinklator/internal/handlers"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"github.com/ilya-rusyanov/shrinklator/internal/server"
	"github.com/ilya-rusyanov/shrinklator/internal/services"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"
)

func newRouter(shortenHandler http.HandlerFunc, expandHandler http.HandlerFunc) chi.Router {
	r := chi.NewRouter()
	r.Use(logger.Middleware)
	r.Post("/", shortenHandler)
	r.Get("/{id}", expandHandler)
	return r
}

func main() {
	config := config.New()
	config.Parse()

	logger.Initialize(config.LogLevel)

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
