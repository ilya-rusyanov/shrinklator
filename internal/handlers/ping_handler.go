package handlers

import (
	"context"
	"net/http"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
)

// Pinger - Ping use case
type Pinger interface {
	Ping(context.Context) error
}

// PingHandler - pings the service to check it's working capacity
type PingHandler struct {
	log    *logger.Log
	pinger Pinger
}

// NewPing constructs Ping handler
func NewPing(log *logger.Log, pinger Pinger) *PingHandler {
	return &PingHandler{
		log:    log,
		pinger: pinger,
	}
}

// Handler handles HTTP requests
func (h *PingHandler) Handler(rw http.ResponseWriter, r *http.Request) {
	err := h.pinger.Ping(r.Context())

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
