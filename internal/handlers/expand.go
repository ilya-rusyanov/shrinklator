package handlers

import (
	"net/http"
	"strings"
)

type Expand struct {
	shrinker shrinker
}

func NewExpand(shrinker shrinker) *Expand {
	return &Expand{
		shrinker: shrinker,
	}
}

func (e *Expand) Handler(rw http.ResponseWriter, r *http.Request) {
	id := strings.TrimLeft(r.URL.Path, "/")

	url, err := e.shrinker.Expand(r.Context(), id)

	if err != nil {
		http.Error(rw, "not found", http.StatusBadRequest)
		return
	}

	rw.Header().Add("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
