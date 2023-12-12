package server

import (
	"fmt"
	"net/http"
)

// Server - main application server
type Server struct {
	srv http.Server
}

// New consturcts a server
func New(addr string, handler http.Handler) Server {
	return Server{
		srv: http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

// Run starts main application server
func (s *Server) Run() error {
	err := s.srv.ListenAndServe()

	if err != nil {
		return fmt.Errorf("failed to run the server: %w", err)
	}

	return nil
}
