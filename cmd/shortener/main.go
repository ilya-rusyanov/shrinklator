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
	restShortener http.HandlerFunc) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.NewLogger(log).Middleware())
	r.Use(middleware.Gzip)
	r.Post("/", shortenHandler)
	r.Get("/{id}", expandHandler)
	r.Post("/api/shorten", restShortener)
	return r
}

func main() {
	config := config.New()
	config.Parse()

	log, err := logger.NewLogger(config.LogLevel)
	if err != nil {
		panic(err)
	}

	var repository storage.Interface
	if config.StoreInFile {
		file, err := storage.NewFile(log, config.FileStoragePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		repository = file
	} else {
		repository = storage.NewInMemory(log)
	}

	shortenerService := services.NewShortener(repository)

	shortenHandler := handlers.NewShorten(log, shortenerService, config.BasePath)
	expandHandler := handlers.NewExpand(shortenerService)
	restShortenerHandler := handlers.NewShortenREST(log, shortenerService,
		config.BasePath)

	router := newRouter(
		log,
		shortenHandler.Handler(),
		expandHandler.Handler(),
		restShortenerHandler.Handler())

	err = server.Run(config.ListenAddr, router)
	if err != nil {
		panic(err)
	}
}
