package server

import (
	"net/http"

	"github.com/ilya-rusyanov/shrinklator/internal/handlers"
	"github.com/ilya-rusyanov/shrinklator/internal/models"
	"github.com/ilya-rusyanov/shrinklator/internal/storage"
)

func Run() {
	sh := models.New(storage.New())
	err := http.ListenAndServe("localhost:8080", handlers.NewHandler(sh))
	if err != nil {
		panic(err)
	}
}
