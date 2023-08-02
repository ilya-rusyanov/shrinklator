package server

import (
	"net/http"

	"github.com/ilya-rusyanov/shrinklator/internal/config"
	"github.com/ilya-rusyanov/shrinklator/internal/handlers"
	"github.com/ilya-rusyanov/shrinklator/internal/models"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"
)

func Run() {
	sh := models.New(storage.New())

	err := http.ListenAndServe(config.Values.ListenAddr,
		handlers.NewHandler(sh))

	if err != nil {
		panic(err)
	}
}
