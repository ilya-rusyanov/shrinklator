package handlers

import (
	"context"
	"net/http"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
)

type BatchServicer interface {
	BatchShorten(context.Context, []entities.CorrelationPair) (
		[]entities.CorrelationPair, error)
}

type BatchShorten struct {
	log     *logger.Log
	service BatchServicer
}

func NewBatchShorten(log *logger.Log, service BatchServicer) *BatchShorten {
	return &BatchShorten{
		log:     log,
		service: service,
	}
}

func (s *BatchShorten) Handler(rw http.ResponseWriter, r *http.Request) {
}
