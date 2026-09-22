package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/goforge/ai-platform/internal/gateway"
	"github.com/goforge/ai-platform/pkg/config"
	"github.com/goforge/ai-platform/pkg/logger"
	"github.com/goforge/ai-platform/pkg/middleware"
	pb "github.com/goforge/ai-platform/pkg/proto"
	"github.com/goforge/ai-platform/pkg/ratelimit"
)

func main() {
	cfg := config.LoadConfig()

	if err := logger.InitLogger(cfg.Observability.LogLevel); err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting API Gateway",
		zap.String("version", cfg.Service.Version),
		zap.Int("port", cfg.GRPC.Port),
	)

	serviceEndpoints := map[string]string{
		"model-service":   "model-service:50051",
		"auth-service":    "auth-service:50052",
		"metrics-service": "metrics-service:50053",
	}

	rateLimiter := ratelimit.New(cfg.RateLimit.RequestsPerSecond, cfg.RateLimit.BurstSize)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			middleware.MetricsUnaryInterceptor("api-gateway"),
			rateLimiter.UnaryServerInterceptor(),
		)),
	)

	gatewayServer := gateway.NewServer(logger.Log, serviceEndpoints)
	pb.RegisterGatewayServiceServer(grpcServer, gatewayServer)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("api-gateway", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)
	grpc_prometheus.Register(grpcServer)

	go startMetricsServer(cfg.Observability.PrometheusPort)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	go func() {
		logger.Info("API Gateway listening", zap.Int("port", cfg.GRPC.Port))
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down API Gateway")
	gatewayServer.Close()
	grpcServer.GracefulStop()
}

func startMetricsServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	logger.Info("Metrics server listening", zap.Int("port", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Error("Metrics server failed", zap.Error(err))
	}
}
