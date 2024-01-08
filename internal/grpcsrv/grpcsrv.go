package grpcsrv

import (
	"context"
	"fmt"
	"net"

	pb "github.com/ilya-rusyanov/shrinklator/proto"
	"google.golang.org/grpc"
)

// Server - gRPC server
type Server struct {
	service *Service
}

// New constructs new gRPC server
func New(s *Service) (*Server, error) {
	return &Server{
		service: s,
	}, nil
}

// Run runs the server
func (s *Server) Run() error {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		return err
	}
	// создаём gRPC-сервер без зарегистрированной службы
	srv := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterShortenerServer(srv, s.service)

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := srv.Serve(listen); err != nil {
		return err
	}

	return nil
}

// Stop stops the server
func (s *Server) Stop(ctx context.Context) error {
	// TODO: implement
	return nil
}
