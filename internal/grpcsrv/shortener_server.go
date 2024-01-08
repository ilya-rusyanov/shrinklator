package grpcsrv

import (
	"context"

	pb "github.com/ilya-rusyanov/shrinklator/proto"
)

// ShrotenerServer is shortener implementation in gRPC
type Service struct {
	pb.UnimplementedShortenerServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Shorten(context.Context, *pb.URL) (*pb.URL, error) {
	return nil, nil
}

func (s *Service) Expand(context.Context, *pb.URL) (*pb.URL, error) {
	return nil, nil
}

func (s *Service) Ping(context.Context, *pb.Empty) (*pb.Empty, error) {
	return nil, nil
}

func (s *Service) Batch(context.Context, *pb.BatchPayload) (*pb.BatchPayload, error) {
	return nil, nil
}

func (s *Service) List(context.Context, *pb.Empty) (*pb.URLs, error) {
	return nil, nil
}

func (s *Service) DeleteRequest(context.Context, *pb.URLs) (*pb.Empty, error) {
	return nil, nil
}
