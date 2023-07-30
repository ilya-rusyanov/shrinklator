package app

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ilya-rusyanov/shrinklator/internal/shortener"
)

type shrinker interface {
	Shrink(string) (string, error)
	Expand(string) (string, error)
}

type shortenerHandler struct {
	shrinker shrinker
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
		io.WriteString(w, err.Error()+"\n")
	}
}

func (h *shortenerHandler) handlePost(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		return fmt.Errorf("unsupported POST path %q", r.URL.Path)
	}

	/*
		contentType := r.Header.Get("Content-Type")

			if contentType != "text/plain" {
				return fmt.Errorf("unsupported content type %q", contentType)
			}
	*/

	sb := &strings.Builder{}
	io.Copy(sb, r.Body)
	short, err := h.shrinker.Shrink(sb.String())

	if err != nil {
		return fmt.Errorf("error shortening: %w", err)
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	self := scheme + "://" + r.Host

	result := self + "/" + short

	io.WriteString(w, result+"\n")

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

func Run() {
	sh := shortener.New()
	http.ListenAndServe("localhost:8080", &shortenerHandler{sh})
}
