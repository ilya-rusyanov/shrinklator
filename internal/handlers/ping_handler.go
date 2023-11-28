package handlers

import (
	"context"
	"net/http"
)

// Pinger - Ping use case
type Pinger interface {
	Ping(context.Context) error
}

// PingHandler - pings the service to check it's working capacity
type PingHandler struct {
	log    Logger
	pinger Pinger
}

// NewPing constructs Ping handler
func NewPing(log Logger, pinger Pinger) *PingHandler {
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
