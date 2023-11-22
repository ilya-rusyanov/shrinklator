package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

// DeleteService - usecase for URLs deletion
type DeleteService interface {
	Delete(context.Context, entities.DeleteRequest) error
}

// Delete - asynchroniously deletes shotened URLs
type Delete struct {
	log     Logger
	service DeleteService
}

// NewDeleteHandler constructs Delete handler
func NewDeleteHandler(log Logger, service DeleteService) *Delete {
	return &Delete{
		log:     log,
		service: service,
	}
}

// Handler handles HTTP requests
func (d *Delete) Handler(rw http.ResponseWriter, r *http.Request) {
	var urls []string

	err := json.NewDecoder(r.Body).Decode(&urls)

	if err != nil {
		http.Error(rw,
			fmt.Sprintf("error decoding request: %q", err.Error()),
			http.StatusBadRequest)
		return
	}

	uid := getUID(r.Context())

	if uid == nil {
		http.Error(rw, "only authorized users are allowed",
			http.StatusUnauthorized)
		return
	}

	var deleteRequest entities.DeleteRequest
	for _, url := range urls {
		deleteRequest = append(deleteRequest, entities.UserAndShort{
			URL: url,
			UID: *uid,
		})
	}

	err = d.service.Delete(r.Context(), deleteRequest)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}
