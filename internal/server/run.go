package server

import (
	"fmt"
	"net/http"
)

func Run(listenAddr string, handler http.Handler) error {
	err := http.ListenAndServe(listenAddr, handler)

	if err != nil {
		return fmt.Errorf("failed to run the server: %w", err)
	}

	return nil
}
