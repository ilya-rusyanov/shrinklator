package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Shrinker interface {
	Shrink(string) string
	Expand(string) (string, error)
}

type shortenerHandler struct {
	shrinker Shrinker
}

func (h *shortenerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	switch r.Method {
	case http.MethodPost:
		err = h.handlePost(w, r)
	case http.MethodGet:
		err = h.handleGet(w, r)
	default:
		err = fmt.Errorf("unsupported method %q", r.Method)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
	}
}

func (h *shortenerHandler) handlePost(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		return fmt.Errorf("unsupported POST path %q", r.URL.Path)
	}

	sb := &strings.Builder{}
	io.Copy(sb, r.Body)
	short := h.shrinker.Shrink(sb.String())

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	self := scheme + "://" + r.Host

	result := self + "/" + short

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, result)

	return nil
}

func (h *shortenerHandler) handleGet(w http.ResponseWriter, r *http.Request) error {
	id := strings.TrimLeft(r.URL.Path, "/")

	url, err := h.shrinker.Expand(id)

	if err != nil {
		return fmt.Errorf("error expanding url: %w", err)
	}

	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

	return nil
}

func NewHandler(s Shrinker) *shortenerHandler {
    return &shortenerHandler{s}
}
