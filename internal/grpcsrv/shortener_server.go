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

// Pinger is UC for pinging DB
type Pinger interface {
	Ping(context.Context) error
}

// BatchServicer shortens URLs in bulk
type BatchServicer interface {
	BatchShorten(context.Context, []entities.BatchRequest) ([]entities.BatchResponse, error)
}

// UserLister lists URLs for certain user
type UserLister interface {
	URLsForUser(context.Context, entities.UserID) (entities.PairArray, error)
}

// Deleter asynchroniously deletes user URLs
type Deleter interface {
	Delete(context.Context, entities.DeleteRequest) error
}

// ShrotenerServer is shortener implementation in gRPC
type Service struct {
	pb.UnimplementedShortenerServer

	basePath       string
	shortenService ShortenerService
	pinger         Pinger
	batchServicer  BatchServicer
	userLister     UserLister
	deleter        Deleter
}

// NewService constructs gRPC service
func NewService(
	base string,
	shortenService ShortenerService,
	pinger Pinger,
	batcher BatchServicer,
	lister UserLister,
	deleter Deleter,
) *Service {
	return &Service{
		basePath:       base,
		shortenService: shortenService,
		pinger:         pinger,
		batchServicer:  batcher,
		userLister:     lister,
		deleter:        deleter,
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

	return &pb.Empty{}, nil
}

// Batch bulk shortens URLs
func (s *Service) Batch(ctx context.Context, req *pb.BatchPayload) (*pb.BatchPayload, error) {
	var res pb.BatchPayload

	in := make([]entities.BatchRequest, len(req.Units))
	for i, u := range req.Units {
		in[i].ID = u.CorrelationId
		in[i].LongURL = u.Url
	}

	shortened, err := s.batchServicer.BatchShorten(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res.Units = make([]*pb.BatchPayload_Unit, len(shortened))
	for i, u := range shortened {
		var unit pb.BatchPayload_Unit

		unit.CorrelationId = u.ID
		unit.Url = s.basePath + "/" + u.ShortURL

		res.Units[i] = &unit
	}

	return &res, nil
}

// List list URLs for user
func (s *Service) List(ctx context.Context, empty *pb.Empty) (*pb.URLs, error) {
	user := getUID(ctx)
	if user == nil {
		return nil, status.Errorf(codes.PermissionDenied, "user ID is expected")
	}

	urls, err := s.userLister.URLsForUser(ctx, *user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var res pb.URLs
	res.Links = make([]string, len(urls))

	for i, u := range urls {
		res.Links[i] = s.basePath + "/" + u.Short
	}

	return &res, nil
}

// Delete deletes URLs for user
func (s *Service) Delete(ctx context.Context, req *pb.URLs) (*pb.Empty, error) {
	user := getUID(ctx)
	if user == nil {
		return nil, status.Errorf(codes.PermissionDenied, "only authorized users are allowed")
	}

	deleteRequest := make([]entities.UserAndShort, len(req.Links))
	for i, l := range req.Links {
		deleteRequest[i].URL = l
		deleteRequest[i].UID = *user
	}

	err := s.deleter.Delete(ctx, deleteRequest)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}
