package grpcsrv

import (
	"context"
	"strings"

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

type Pinger interface {
	Ping(context.Context) error
}

// ShrotenerServer is shortener implementation in gRPC
type Service struct {
	pb.UnimplementedShortenerServer

	basePath       string
	shortenService ShortenerService
	pinger         Pinger
}

// NewService constructs gRPC service
func NewService(
	base string,
	shortenService ShortenerService,
	pinger Pinger,
) *Service {
	return &Service{
		basePath:       base,
		shortenService: shortenService,
		pinger:         pinger,
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
func (s *Service) Expand(ctx context.Context, url *pb.URL) (*pb.URL, error) {
	urlPart := strings.TrimLeft(url.Link, "/")

	expandResult, err := s.shortenService.Expand(ctx, urlPart)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "entry not found")
	}

	if expandResult.Removed {
		return nil, status.Errorf(codes.DataLoss, "entry is removed")
	}

	return &pb.URL{Link: expandResult.URL}, nil
}

// Ping pings database
func (s *Service) Ping(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	err := s.pinger.Ping(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

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
