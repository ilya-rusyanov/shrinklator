package handlers

import (
	"net/http"
	"strings"
)

// Expand - expands shortened URL in plain text
type Expand struct {
	shrinker shrinker
}

// NewExpand constructs Expand objext
func NewExpand(shrinker shrinker) *Expand {
	return &Expand{
		shrinker: shrinker,
	}
}

// Handler handles HTTP requests
func (e *Expand) Handler(rw http.ResponseWriter, r *http.Request) {
	id := strings.TrimLeft(r.URL.Path, "/")

	expandResult, err := e.shrinker.Expand(r.Context(), id)

	if err != nil {
		http.Error(rw, "not found", http.StatusBadRequest)
		return
	}

	if expandResult.Removed {
		http.Error(rw, "entry is removed", http.StatusGone)
		return
	}

	rw.Header().Set("Content-Type", "text/plain")
	rw.Header().Add("Location", expandResult.URL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
