package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

var errNoUserID = errors.New("user ID is not specified")

type URLsService interface {
	URLsForUser(context.Context, entities.UserID) (entities.PairArray, error)
}

type UserURLs struct {
	log     *logger.Log
	service URLsService
	baseURL string
}

func NewUserURLs(log *logger.Log, service URLsService,
	baseURL string) *UserURLs {
	return &UserURLs{
		log:     log,
		service: service,
		baseURL: baseURL,
	}
}

func (u *UserURLs) Handler(rw http.ResponseWriter, r *http.Request) {
	id, err := u.getUID(r.Context())

	if err != nil {
		if errors.Is(err, errNoUserID) {
			http.Error(rw, "user ID is expected", http.StatusUnauthorized)
			return
		} else {
			u.log.Error("unexpected failure to retrieve user id",
				zap.String("err", err.Error()))
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}

	urls, err := u.service.URLsForUser(r.Context(), id)

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
	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(urls); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u *UserURLs) getUID(ctx context.Context) (id entities.UserID, err error) {
	if id := ctx.Value("uid"); id != nil {
		return id.(entities.UserID), nil
	}

	return entities.UserID(0), errNoUserID
}
