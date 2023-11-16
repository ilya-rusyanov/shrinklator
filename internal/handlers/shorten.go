package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

// Shorten - shortens URL to plain text
type Shorten struct {
	shrinker shrinker
	basePath string
	log      *logger.Log
}

// NewShorten constructs Shorten handler
func NewShorten(log *logger.Log, shrinker shrinker, basePath string) *Shorten {
	return &Shorten{
		shrinker: shrinker,
		basePath: basePath,
		log:      log,
	}
}

// Handler handles HTTP requests
func (s *Shorten) Handler(rw http.ResponseWriter, r *http.Request) {
	sb := &strings.Builder{}
	io.Copy(sb, r.Body)
	status := http.StatusCreated
	uid := getUID(r.Context())
	s.log.Sugar().Infof("request to shorten %q with headers %#v", sb.String(), r.Header)
	short, err := s.shrinker.Shrink(r.Context(), sb.String(), uid)
	if err != nil {
		if short, err = handleAlreadyExists(err, &status); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			s.log.Error("error shortening URL",
				zap.String("message", err.Error()))
			return
		}
	}

	result := s.basePath + "/" + short

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(status)
	s.respondWithString(rw, result)
}

func (s *Shorten) respondWithString(rw http.ResponseWriter, text string) {
	if _, err := io.WriteString(rw, text); err != nil {
		s.log.Error("error writing response",
			zap.String("message", err.Error()))
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
