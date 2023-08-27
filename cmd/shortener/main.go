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

type Persistence interface {
	Append(short, long string) error
	ReadAll() (values map[string]string, err error)
}

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

	persistence := Persistence(storage.NewNullPersistence())
	if config.StoreInFile {
		file, err := storage.NewFile(config.FileStoragePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		persistence = storage.NewRWpersistence(file, file)
	}

	values, err := persistence.ReadAll()
	if err != nil {
		panic(err)
	}

	inMemory := storage.NewInMemory(values)

	hybridStorage := storage.NewHybrid(inMemory, persistence)

	shortenerService := services.NewShortener(hybridStorage)

	shortenHandler := handlers.Shorten(shortenerService, config.BasePath)
	expandHandler := handlers.Expand(shortenerService)
	restShortenerHandler := handlers.ShortenREST(shortenerService,
		config.BasePath)

	router := newRouter(shortenHandler, expandHandler, restShortenerHandler)

	err = server.Run(config.ListenAddr, router)
	if err != nil {
		panic(err)
	}
}
