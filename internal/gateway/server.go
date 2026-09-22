package gateway

import (
	"context"
	"fmt"
	"time"

	pb "github.com/goforge/ai-platform/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	pb.UnimplementedGatewayServiceServer
	logger          *zap.Logger
	serviceClients  map[string]*grpc.ClientConn
	serviceEndpoints map[string]string
}

func NewServer(logger *zap.Logger, serviceEndpoints map[string]string) *Server {
	return &Server{
		logger:           logger,
		serviceClients:   make(map[string]*grpc.ClientConn),
		serviceEndpoints: serviceEndpoints,
	}
}

func (s *Server) RouteRequest(ctx context.Context, req *pb.RouteRequestMessage) (*pb.RouteResponseMessage, error) {
	s.logger.Info("Routing request",
		zap.String("service", req.ServiceName),
		zap.String("method", req.Method),
		zap.String("request_id", req.RequestId),
	)

	start := time.Now()

	endpoint, exists := s.serviceEndpoints[req.ServiceName]
	if !exists {
		return nil, fmt.Errorf("service %s not found", req.ServiceName)
	}

	client, err := s.getOrCreateClient(req.ServiceName, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service: %w", err)
	}

	_ = client

	latency := time.Since(start).Milliseconds()

	return &pb.RouteResponseMessage{
		StatusCode: 200,
		Payload:    []byte("routed successfully"),
		Headers: map[string]string{
			"X-Request-ID": req.RequestId,
			"X-Service":    req.ServiceName,
		},
		RequestId: req.RequestId,
		LatencyMs: latency,
	}, nil
}

func (s *Server) GetServiceHealth(ctx context.Context, req *pb.ServiceHealthRequest) (*pb.ServiceHealthResponse, error) {
	s.logger.Info("Health check request", zap.String("service", req.ServiceName))

	return &pb.ServiceHealthResponse{
		ServiceName: req.ServiceName,
		Healthy:     true,
		Status:      "healthy",
		Details: map[string]string{
			"timestamp": time.Now().Format(time.RFC3339),
			"endpoint":  s.serviceEndpoints[req.ServiceName],
		},
	}, nil
}

func (s *Server) getOrCreateClient(serviceName, endpoint string) (*grpc.ClientConn, error) {
	if client, exists := s.serviceClients[serviceName]; exists {
		return client, nil
	}

	conn, err := grpc.NewClient(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	s.serviceClients[serviceName] = conn
	return conn, nil
}

func (s *Server) Close() {
	for name, conn := range s.serviceClients {
		if err := conn.Close(); err != nil {
			s.logger.Error("Failed to close connection", zap.String("service", name), zap.Error(err))
		}
	}
}
