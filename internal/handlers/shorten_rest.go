package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

func ShortenREST(shrinker shrinker, basePath string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		buf := bytes.Buffer{}
		var (
			shortenRequest struct {
				URL string `json:"url"`
			}
			result = make(map[string]string, 1)
		)

		if _, err := buf.ReadFrom(r.Body); err != nil {
			logger.Log.Error("cannot read request body",
				zap.String("message", err.Error()))
			http.Error(rw,
				fmt.Sprintf("error reading request body: %v", err),
				http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(buf.Bytes(), &shortenRequest); err != nil {
			http.Error(rw, err.Error(),
				http.StatusBadRequest)
			logger.Log.Error("error marshaling JSON",
				zap.String("message", err.Error()))
			return
		}

		short, err := shrinker.Shrink(shortenRequest.URL)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			logger.Log.Error("error shortening URL",
				zap.String("message", err.Error()))
			return
		}

		result["result"] = basePath + "/" + short

		resultJSON, err := json.Marshal(result)
		if err != nil {
			http.Error(rw,
				fmt.Sprintf("error serializing response: %v", err),
				http.StatusInternalServerError)
			logger.Log.Error("error marshaling JSON",
				zap.String("message", err.Error()))
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)

		if _, err = rw.Write(resultJSON); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			logger.Log.Error("unable to write response",
				zap.String("error", err.Error()))
		}
	}
}
