package server

import (
	"context"
	"fmt"
	"net/http"
)

// Opt is a funcopts for Server
type Opt func(*http.Server) error

// Logger logs messages
type Logger interface {
	Debug(...any)
}

// Server - main application server
type Server struct {
	srv http.Server
	log Logger
}

// New consturcts a server
func New(
	logger Logger,
	addr string,
	handler http.Handler,
	opts ...Opt,
) (*Server, error) {
	res := Server{
		srv: http.Server{
			Addr:    addr,
			Handler: handler,
		},
		log: logger,
	}

	for _, opt := range opts {
		err := opt(&res.srv)
		if err != nil {
			return nil, fmt.Errorf("failed to set server opt: %w", err)
		}
	}

	return &res, nil
}

// Run starts main application server
func (s *Server) Run() error {
	if s.srv.TLSConfig == nil {
		s.log.Debug("starting insecure server on ", s.srv.Addr)
		err := s.srv.ListenAndServe()
		if err != nil {
			return fmt.Errorf("failed to run insecure server: %w", err)
		}
	} else {
		s.log.Debug("starting https server on ", s.srv.Addr)
		err := s.srv.ListenAndServeTLS("", "")
		if err != nil {
			return fmt.Errorf("failed to run secure server: %w", err)
		}
	}

	return nil
}

// Stop stops main application server
func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
