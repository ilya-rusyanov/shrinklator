package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ilya-rusyanov/shrinklator/internal/config"
	"github.com/ilya-rusyanov/shrinklator/internal/handlers"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"github.com/ilya-rusyanov/shrinklator/internal/server"
	"github.com/ilya-rusyanov/shrinklator/internal/server/middleware"
	"github.com/ilya-rusyanov/shrinklator/internal/services"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"
)

func newRouter(shortenHandler http.HandlerFunc, expandHandler http.HandlerFunc,
	restShortener http.HandlerFunc) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Gzip)
	r.Post("/", shortenHandler)
	r.Get("/{id}", expandHandler)
	r.Post("/api/shorten", restShortener)
	return r
}

func main() {
	config := config.New()
	config.Parse()

	logger.Initialize(config.LogLevel)

	inMemory := storage.NewInMemory()

	persistence := storage.NewNullPersistence()

	hybridStorage := storage.NewHybrid(inMemory, persistence)

	shortenerService := services.NewShortener(hybridStorage)

	shortenHandler := handlers.Shorten(shortenerService, config.BasePath)
	expandHandler := handlers.Expand(shortenerService)
	restShortenerHandler := handlers.ShortenREST(shortenerService,
		config.BasePath)

	router := newRouter(shortenHandler, expandHandler, restShortenerHandler)

	err := server.Run(config.ListenAddr, router)
	if err != nil {
		panic(err)
	}
}
