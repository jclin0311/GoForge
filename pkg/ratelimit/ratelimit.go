package ratelimit

import (
	"context"
	"fmt"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RateLimiter struct {
	limiter *rate.Limiter
}

func New(requestsPerSecond float64, burstSize int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), burstSize),
	}
}

func (rl *RateLimiter) Allow() bool {
	return rl.limiter.Allow()
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	return rl.limiter.Wait(ctx)
}

func (rl *RateLimiter) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if !rl.Allow() {
			return nil, status.Errorf(
				codes.ResourceExhausted,
				fmt.Sprintf("rate limit exceeded for method %s", info.FullMethod),
			)
		}
		return handler(ctx, req)
	}
}
