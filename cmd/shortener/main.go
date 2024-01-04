package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiware "github.com/go-chi/chi/v5/middleware"
	"github.com/ilya-rusyanov/shrinklator/internal/config"
	"github.com/ilya-rusyanov/shrinklator/internal/handlers"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"github.com/ilya-rusyanov/shrinklator/internal/server"
	"github.com/ilya-rusyanov/shrinklator/internal/server/cert"
	"github.com/ilya-rusyanov/shrinklator/internal/server/middleware"
	"github.com/ilya-rusyanov/shrinklator/internal/services"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"
)

const tokenKey string = "this is security flaw"
const accessCookieName string = "access_token"

func newRouter(log Logger, shortenHandler http.HandlerFunc,
	expandHandler http.HandlerFunc,
	restShortener http.HandlerFunc,
	pingHandler http.HandlerFunc,
	batchHandler http.HandlerFunc,
	userURLs http.HandlerFunc,
	del http.HandlerFunc,
	stats http.HandlerFunc,
) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.NewLogger(log).Middleware())
	r.Use(middleware.NewPseudoAuth(log, tokenKey, accessCookieName).Middleware)
	r.Use(middleware.NewGzip(log).Middleware)
	r.Mount("/debug", chiware.Profiler())
	r.Post("/", shortenHandler)
	r.Get("/{id}", expandHandler)
	r.Post("/api/shorten", restShortener)
	r.Get("/ping", pingHandler)
	r.Post("/api/shorten/batch", batchHandler)
	r.Get("/api/user/urls", userURLs)
	r.Delete("/api/user/urls", del)
	r.Get("/api/internal/stats", stats)
	return r
}

func main() {
	printBuildInfo(os.Stdout)

	config := config.New()
	config.MustParse()

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
	userURLsService, deleteErrorsCh := services.NewUserURLs(ctx, repository, config.DelBufSize)
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

	statsHandler, err := handlers.NewStatsHandler(log, repository, config.TrustedSubnet)
	if err != nil {
		panic(err)
	}

	router := newRouter(
		log,
		shortenHandler.Handler,
		expandHandler.Handler,
		restShortenerHandler.Handler,
		pingHandler.Handler,
		batchHandler.Handler,
		userURLsHandler.Handler,
		delHandler.Handler,
		statsHandler.Handler,
	)

	var srvOpts []server.Opt
	if config.Secure {
		srvOpts = append(srvOpts, cert.New())
	}

	server, err := server.New(
		log,
		config.ListenAddr,
		router,
		srvOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}

	allConnsClosed := gracefulShutdown(ctx, log, server)

	err = server.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	<-allConnsClosed
	log.Debug("graceful shutdown")
}
