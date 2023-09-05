package main

import (
	"context"
	"net/http"
	"os"
	"time"

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
	pingHandler http.HandlerFunc) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.NewLogger(log).Middleware())
	r.Use(middleware.Gzip)
	r.Post("/", shortenHandler)
	r.Get("/{id}", expandHandler)
	r.Post("/api/shorten", restShortener)
	r.Get("/ping", pingHandler)
	return r
}

func main() {
	config := config.New()
	config.Parse()

	log, err := logger.NewLogger(config.LogLevel)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	db, err := storage.NewPostgres(ctx, log, config.DSN)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		log.Error("failed to ping DB")
		os.Exit(1)
	}

	var repository storage.Interface
	switch {
	case config.StoreInDB:
		repository = db
	case config.StoreInFile:
		file, err := storage.NewFile(log, config.FileStoragePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		repository = file
	default:
		repository = storage.NewInMemory(log)
	}

	shortenerService := services.NewShortener(repository)
	pingService := services.NewPing(db)

	shortenHandler := handlers.NewShorten(log, shortenerService, config.BasePath)
	expandHandler := handlers.NewExpand(shortenerService)
	restShortenerHandler := handlers.NewShortenREST(log, shortenerService,
		config.BasePath)
	pingHandler := handlers.NewPing(log, pingService)

	router := newRouter(
		log,
		shortenHandler.Handler(),
		expandHandler.Handler(),
		restShortenerHandler.Handler(),
		pingHandler.Handler())

	err = server.Run(config.ListenAddr, router)
	if err != nil {
		panic(err)
	}
}
