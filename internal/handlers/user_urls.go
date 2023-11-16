package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

// URLsService usecase for listing URLs belonging to user
type URLsService interface {
	URLsForUser(context.Context, entities.UserID) (entities.PairArray, error)
}

// UserURLs - lists URLs submitted by user
type UserURLs struct {
	log     *logger.Log
	service URLsService
	baseURL string
}

// NewUserURLs constructs UserURLs object
func NewUserURLs(log *logger.Log, service URLsService,
	baseURL string) *UserURLs {
	return &UserURLs{
		log:     log,
		service: service,
		baseURL: baseURL,
	}
}

// Handler handles HTTP requests
func (u *UserURLs) Handler(rw http.ResponseWriter, r *http.Request) {
	id := getUID(r.Context())

	if id == nil {
		http.Error(rw, "user ID is expected", http.StatusUnauthorized)
		return
	}

	urls, err := u.service.URLsForUser(r.Context(), *id)

	if err != nil {
		u.log.Info("failure to fetch URLs", zap.String("err", err.Error()))
		http.Error(rw, fmt.Sprintf("failed to fetch URLs: %q", err.Error()),
			http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	for i := range urls {
		urls[i].Short = u.baseURL + "/" + urls[i].Short
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(urls); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
