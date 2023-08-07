package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/ilya-rusyanov/shrinklator/internal/config"
)

type shrinker interface {
	Shrink(string) string
	Expand(string) (string, error)
}

func NewHandler(s shrinker) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", shorten(s, config.Values.BasePath))
	r.Get("/{id}", expand(s))
	return r
}

func shorten(shrinker shrinker, basePath string) http.HandlerFunc {
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

func expand(shrinker shrinker) http.HandlerFunc {
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
