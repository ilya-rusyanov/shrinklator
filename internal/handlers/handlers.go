package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"go.uber.org/zap"
)

type shrinker interface {
	Shrink(string) string
	Expand(string) (string, error)
}

func Shorten(shrinker shrinker, basePath string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		sb := &strings.Builder{}
		io.Copy(sb, r.Body)
		short := shrinker.Shrink(sb.String())

		result := basePath + "/" + short

		rw.WriteHeader(http.StatusCreated)
		_, err := io.WriteString(rw, result)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(os.Stderr, "unable to write response: %v", err)
		}
	}
}

func ShortenREST(shrinker shrinker, basePath string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		buf := bytes.Buffer{}
		var (
			shortenRequest struct {
				URL string `json:"url"`
			}
			result = make(map[string]string, 1)
		)

		_, err := buf.ReadFrom(r.Body)

		if err != nil {
			http.Error(rw,
				fmt.Sprintf("error reading request body: %v", err),
				http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(buf.Bytes(), &shortenRequest)

		if err != nil {
			http.Error(rw, err.Error(),
				http.StatusBadRequest)
			return
		}

		result["result"] = basePath + "/" + shrinker.Shrink(shortenRequest.URL)

		resultJSON, err := json.Marshal(result)
		if err != nil {
			http.Error(rw,
				fmt.Sprintf("error serializing response: %v", err),
				http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusCreated)
		_, err = rw.Write(resultJSON)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			logger.Log.Error("unable to write response",
				zap.String("error", err.Error()))
		}
	}
}

func Expand(shrinker shrinker) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id := strings.TrimLeft(r.URL.Path, "/")

		url, err := shrinker.Expand(id)

		if err != nil {
			http.Error(rw, "not found", http.StatusBadRequest)
			return
		}

		rw.Header().Add("Location", url)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}
