package grpcsrv

import (
	"context"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	pb "github.com/ilya-rusyanov/shrinklator/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShortenerService are UC for shortening and expanding
type ShortenerService interface {
	Shrink(context.Context, string, *entities.UserID) (string, error)
	Expand(context.Context, string) (entities.ExpandResult, error)
}

// ShrotenerServer is shortener implementation in gRPC
type Service struct {
	pb.UnimplementedShortenerServer

	basePath       string
	shortenService ShortenerService
}

// NewService constructs gRPC service
func NewService(base string, shortenService ShortenerService) *Service {
	return &Service{
		basePath:       base,
		shortenService: shortenService,
	}
}

// Shorten shortens URL
func (s *Service) Shorten(ctx context.Context, url *pb.URL) (*pb.URL, error) {
	var (
		code     error
		response pb.URL
		err      error
	)

	user := getUID(ctx)

	response.Link, err = s.shortenService.Shrink(ctx, url.Link, user)
	if err != nil {
		if response.Link, err = handleAlreadyExists(err, &code); err != nil {
			return nil, status.Errorf(codes.Internal, "cannot shorten URL: %w", err)
		}
	}

	return &response, code
}

// Expand expands URL
func (s *Service) Expand(context.Context, *pb.URL) (*pb.URL, error) {
	return nil, nil
}

// Ping pings database
func (s *Service) Ping(context.Context, *pb.Empty) (*pb.Empty, error) {
	return nil, nil
}

// Batch bulk shortens URLs
func (s *Service) Batch(context.Context, *pb.BatchPayload) (*pb.BatchPayload, error) {
	return nil, nil
}

// List list URLs for user
func (s *Service) List(context.Context, *pb.Empty) (*pb.URLs, error) {
	return nil, nil
}

// Delete deletes URLs for user
func (s *Service) Delete(context.Context, *pb.URLs) (*pb.Empty, error) {
	return nil, nil
}
