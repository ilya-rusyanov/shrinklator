package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// NewAuth constructs authenticating interceptor
func NewAuth(tokenKey string, cookieName string) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var token string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get(cookieName)
			if len(values) > 0 {
				token = values[0]
			}
		}
		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}
		if token != tokenKey {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		return handler(ctx, req)
	})
}
