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

func newRouter(log *logger.Log, shortenHandler http.HandlerFunc,
	expandHandler http.HandlerFunc,
	restShortener http.HandlerFunc,
	pingHandler http.HandlerFunc,
	batchHandler http.HandlerFunc) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.NewLogger(log).Middleware())
	r.Use(middleware.Gzip)
	r.Post("/", shortenHandler)
	r.Get("/{id}", expandHandler)
	r.Post("/api/shorten", restShortener)
	r.Get("/ping", pingHandler)
	r.Post("/api/shorten/batch", batchHandler)
	return r
}

func main() {
	config := config.New()
	config.Parse()

	log, err := logger.NewLogger(config.LogLevel)
	if err != nil {
		panic(err)
	}

	factory := storage.Factory{}

	repository := factory.MustInitStorage(*config, log)
	defer repository.MustClose()
	algorithm := services.MD5Algo

	shortenerService := services.NewShortener(repository, algorithm)
	pingService := services.NewPing(repository)
	batchService := services.NewBatch(repository, algorithm)

	shortenHandler := handlers.NewShorten(log, shortenerService, config.BasePath)
	expandHandler := handlers.NewExpand(shortenerService)
	restShortenerHandler := handlers.NewShortenREST(log, shortenerService,
		config.BasePath)
	pingHandler := handlers.NewPing(log, pingService)
	batchHandler := handlers.NewBatchShorten(log, batchService, config.BasePath)

	router := newRouter(
		log,
		shortenHandler.Handler,
		expandHandler.Handler,
		restShortenerHandler.Handler,
		pingHandler.Handler,
		batchHandler.Handler)

	err = server.Run(config.ListenAddr, router)
	if err != nil {
		panic(err)
	}
}
