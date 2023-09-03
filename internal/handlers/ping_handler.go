package handlers

import (
	"context"
	"net/http"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
)

type Pinger interface {
	Ping(context.Context) error
}

type PingHandler struct {
	log    *logger.Log
	pinger Pinger
}

func NewPing(log *logger.Log, pinger Pinger) *PingHandler {
	return &PingHandler{
		log:    log,
		pinger: pinger,
	}
}

func (h *PingHandler) Handler() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := h.pinger.Ping(r.Context())

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}
