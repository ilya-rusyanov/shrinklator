package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

func Shorten(shrinker shrinker, basePath string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		sb := &strings.Builder{}
		io.Copy(sb, r.Body)
		short, err := shrinker.Shrink(sb.String())

		if err != nil {
			logger.Log.Error("error shortening",
				zap.String("message", err.Error()))
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		result := basePath + "/" + short

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
