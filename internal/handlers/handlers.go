package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Shrinker interface {
	Shrink(string) string
	Expand(string) (string, error)
}

func postHandler(shrinker Shrinker) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		sb := &strings.Builder{}
		io.Copy(sb, r.Body)
		short := shrinker.Shrink(sb.String())

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		self := scheme + "://" + r.Host

		result := self + "/" + short

		rw.WriteHeader(http.StatusCreated)
		io.WriteString(rw, result)
	}
}

func getHandler(shrinker Shrinker) http.HandlerFunc {
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

func NewHandler(s Shrinker) http.Handler {
	router := chi.NewRouter()
	router.Post("/", postHandler(s))
	router.Get("/{id}", getHandler(s))
	return router
}
