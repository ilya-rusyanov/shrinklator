package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
)

// BatchServicer - usecase for shortening URLs in bulk
type BatchServicer interface {
	BatchShorten(context.Context, []entities.BatchRequest) (
		[]entities.BatchResponse, error)
}

// BatchShorten - bulk shortens URLs
type BatchShorten struct {
	log     *logger.Log
	service BatchServicer
	baseURL string
}

// NewBatchShorten - constructor for BatchShorten
func NewBatchShorten(log *logger.Log, service BatchServicer, baseURL string) *BatchShorten {
	return &BatchShorten{
		log:     log,
		service: service,
		baseURL: baseURL,
	}
}

// Handler handles HTTP request
func (s *BatchShorten) Handler(rw http.ResponseWriter, r *http.Request) {
	var batchRequest []entities.BatchRequest
	if err := json.NewDecoder(r.Body).Decode(&batchRequest); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	shortened, err := s.service.BatchShorten(r.Context(), batchRequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range shortened {
		shortened[i].ShortURL = s.baseURL + "/" + shortened[i].ShortURL
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(shortened); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
