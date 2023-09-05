package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

type Shorten struct {
	shrinker shrinker
	basePath string
	log      *logger.Log
}

func NewShorten(log *logger.Log, shrinker shrinker, basePath string) *Shorten {
	return &Shorten{
		shrinker: shrinker,
		basePath: basePath,
		log:      log,
	}
}

func (s *Shorten) Handler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		sb := &strings.Builder{}
		io.Copy(sb, r.Body)
		short, err := s.shrinker.Shrink(r.Context(), sb.String())

		if err != nil {
			s.log.Error("error shortening",
				zap.String("message", err.Error()))
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		result := s.basePath + "/" + short

		rw.WriteHeader(http.StatusCreated)
		s.respondWithString(rw, result)
	}
}

func (s *Shorten) respondWithString(rw http.ResponseWriter, text string) {
	if _, err := io.WriteString(rw, text); err != nil {
		s.log.Error("error writing response",
			zap.String("message", err.Error()))
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
