package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// ShortenREST - handler for JSON shorten requests
type ShortenREST struct {
	shrinker shrinker
	basePath string
	log      Logger
}

// NewShortenREST constructs ShortenREST objects
func NewShortenREST(log Logger, shrinker shrinker, basePath string) *ShortenREST {
	return &ShortenREST{
		shrinker: shrinker,
		basePath: basePath,
		log:      log,
	}
}

// Handler handles HTTP requests
func (s *ShortenREST) Handler(rw http.ResponseWriter, r *http.Request) {
	buf := bytes.Buffer{}
	var (
		shortenRequest struct {
			URL string `json:"url"`
		}
		result = make(map[string]string, 1)
	)

	if _, err := buf.ReadFrom(r.Body); err != nil {
		s.log.Error("cannot read request body",
			zap.String("message", err.Error()))
		http.Error(rw,
			fmt.Sprintf("error reading request body: %v", err),
			http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &shortenRequest); err != nil {
		http.Error(rw, err.Error(),
			http.StatusBadRequest)
		s.log.Error("error marshaling JSON",
			zap.String("message", err.Error()))
		return
	}

	status := http.StatusCreated
	uid := getUID(r.Context())
	short, err := s.shrinker.Shrink(r.Context(), shortenRequest.URL, uid)
	if err != nil {
		if short, err = handleAlreadyExists(err, &status); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			s.log.Error("error shortening URL",
				zap.String("message", err.Error()))
			return
		}
	}

	result["result"] = s.basePath + "/" + short

	resultJSON, err := json.Marshal(result)
	if err != nil {
		http.Error(rw,
			fmt.Sprintf("error serializing response: %v", err),
			http.StatusInternalServerError)
		s.log.Error("error marshaling JSON",
			zap.String("message", err.Error()))
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)

	if _, err = rw.Write(resultJSON); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		s.log.Error("unable to write response",
			zap.String("error", err.Error()))
	}
}
