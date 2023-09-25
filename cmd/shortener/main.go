package main

import (
	"context"
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

const tokenKey string = "this is security flaw"
const accessCookieName string = "access_token"

func newRouter(log *logger.Log, shortenHandler http.HandlerFunc,
	expandHandler http.HandlerFunc,
	restShortener http.HandlerFunc,
	pingHandler http.HandlerFunc,
	batchHandler http.HandlerFunc,
	userURLs http.HandlerFunc,
	del http.HandlerFunc) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.NewLogger(log).Middleware())
	r.Use(middleware.Gzip)
	r.Use(middleware.NewPseudoAuth(log, tokenKey, accessCookieName).Middleware)
	r.Post("/", shortenHandler)
	r.Get("/{id}", expandHandler)
	r.Post("/api/shorten", restShortener)
	r.Get("/ping", pingHandler)
	r.Post("/api/shorten/batch", batchHandler)
	r.Get("/api/user/urls", userURLs)
	r.Delete("/api/user/urls", del)
	return r
}

func main() {
	config := config.New()
	config.Parse()

	log, err := logger.NewLogger(config.LogLevel)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	factory := storage.Factory{}

	repository := factory.MustInitStorage(*config, log)
	defer repository.MustClose()
	algorithm := services.MD5Algo

	shortenerService := services.NewShortener(log, repository, algorithm)
	pingService := services.NewPing(repository)
	batchService := services.NewBatch(repository, algorithm)
	userURLsService, deleteErrorsCh := services.NewUserURLs(repository, ctx)
	defer userURLsService.Close()
	go printDeleteErrors(log, deleteErrorsCh)

	shortenHandler := handlers.NewShorten(log, shortenerService, config.BasePath)
	expandHandler := handlers.NewExpand(shortenerService)
	restShortenerHandler := handlers.NewShortenREST(log, shortenerService,
		config.BasePath)
	pingHandler := handlers.NewPing(log, pingService)
	batchHandler := handlers.NewBatchShorten(log, batchService, config.BasePath)
	userURLsHandler := handlers.NewUserURLs(log, userURLsService, config.BasePath)
	delHandler := handlers.NewDeleteHandler(log, userURLsService)

	router := newRouter(
		log,
		shortenHandler.Handler,
		expandHandler.Handler,
		restShortenerHandler.Handler,
		pingHandler.Handler,
		batchHandler.Handler,
		userURLsHandler.Handler,
		delHandler.Handler)

	err = server.Run(config.ListenAddr, router)
	if err != nil {
		panic(err)
	}
}
