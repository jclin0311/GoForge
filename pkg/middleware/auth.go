package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	authServiceAddr string
}

func NewAuthInterceptor(authServiceAddr string) *AuthInterceptor {
	return &AuthInterceptor{
		authServiceAddr: authServiceAddr,
	}
}

func (a *AuthInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if strings.HasSuffix(info.FullMethod, "HealthCheck") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata not found")
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token not found")
		}

		token := values[0]
		if !strings.HasPrefix(token, "Bearer ") {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization format")
		}

		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Errorf(codes.Unauthenticated, "metadata not found")
		}

		values := md["authorization"]
		if len(values) == 0 {
			return status.Errorf(codes.Unauthenticated, "authorization token not found")
		}

		return handler(srv, ss)
	}
}
