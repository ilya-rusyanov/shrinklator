package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
)

var errIPParsing = errors.New("failed to parse IP address")

// StatsService - interface to statistics
type StatsService interface {
	CountUsersAndUrls(context.Context) (entities.Stats, error)
}

// StatsHandler - handler of statistics requests
type StatsHandler struct {
	log           Logger
	statsService  StatsService
	trustedSubnet *net.IPNet
}

// NewStatsHandler constructs StatsHandler
func NewStatsHandler(log Logger, statsService StatsService, trustedSubnet string) (*StatsHandler, error) {
	var sub *net.IPNet

	if len(trustedSubnet) > 0 {
		var err error
		_, sub, err = net.ParseCIDR(trustedSubnet)
		if err != nil {
			return nil, fmt.Errorf("failed to parse subnet")
		}
	}

	return &StatsHandler{
		log:           log,
		trustedSubnet: sub,
		statsService:  statsService,
	}, nil
}

// Handler handles stats requests
func (h *StatsHandler) Handler(rw http.ResponseWriter, r *http.Request) {
	if !h.allowed(r) {
		http.Error(rw, "forbidden", http.StatusForbidden)
		return
	}

	stats, err := h.statsService.CountUsersAndUrls(r.Context())
	if err != nil {
		http.Error(rw, "error retrieving statistics", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	err = json.NewEncoder(rw).Encode(stats)
	if err != nil {
		h.log.Error("error writing response: ", err)
	}
}

func (h *StatsHandler) allowed(r *http.Request) bool {
	ip, err := getIP(r)
	if err != nil {
		h.log.Debug("could not parse IP address of request")
		return false
	}

	if h.trustedSubnet == nil || !h.trustedSubnet.Contains(ip) {
		return false
	}

	return true
}

func getIP(r *http.Request) (net.IP, error) {
	ip := net.ParseIP(r.Header.Get("X-Real-IP"))
	if ip == nil {
		return net.IP{}, errIPParsing
	}

	return ip, nil
}
