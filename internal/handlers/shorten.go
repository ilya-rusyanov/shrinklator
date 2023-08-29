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
}

func NewShorten(shrinker shrinker, basePath string) *Shorten {
	return &Shorten{
		shrinker: shrinker,
		basePath: basePath,
	}
}

func (s *Shorten) Handler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		sb := &strings.Builder{}
		io.Copy(sb, r.Body)
		short, err := s.shrinker.Shrink(sb.String())

		if err != nil {
			logger.Log.Error("error shortening",
				zap.String("message", err.Error()))
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		result := s.basePath + "/" + short

		rw.WriteHeader(http.StatusCreated)
		respondWithString(rw, result)
	}
}

func respondWithString(rw http.ResponseWriter, text string) {
	if _, err := io.WriteString(rw, text); err != nil {
		logger.Log.Error("error writing response",
			zap.String("message", err.Error()))
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
