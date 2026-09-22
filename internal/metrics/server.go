package metrics

import (
	"context"
	"time"

	pb "github.com/goforge/ai-platform/pkg/proto"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type Server struct {
	pb.UnimplementedMetricsServiceServer
	logger  *zap.Logger
	metrics map[string][]pb.Metric
}

var (
	customMetrics = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "custom_metric",
			Help: "Custom application metrics",
		},
		[]string{"metric_name", "label_key", "label_value"},
	)
)

func init() {
	prometheus.MustRegister(customMetrics)
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{
		logger:  logger,
		metrics: make(map[string][]pb.Metric),
	}
}

func (s *Server) RecordMetric(ctx context.Context, req *pb.MetricRequest) (*pb.MetricResponse, error) {
	s.logger.Debug("Recording metric",
		zap.String("name", req.MetricName),
		zap.Float64("value", req.Value),
	)

	metric := pb.Metric{
		Name:      req.MetricName,
		Value:     req.Value,
		Labels:    req.Labels,
		Timestamp: req.Timestamp,
		Type:      req.Type,
	}

	s.metrics[req.MetricName] = append(s.metrics[req.MetricName], metric)

	for key, value := range req.Labels {
		customMetrics.WithLabelValues(req.MetricName, key, value).Set(req.Value)
	}

	return &pb.MetricResponse{
		Success: true,
		Message: "Metric recorded successfully",
	}, nil
}

func (s *Server) GetMetrics(ctx context.Context, req *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	s.logger.Debug("Getting metrics", zap.String("name", req.MetricName))

	metrics, exists := s.metrics[req.MetricName]
	if !exists {
		return &pb.GetMetricsResponse{
			Metrics: []*pb.Metric{},
		}, nil
	}

	var filteredMetrics []*pb.Metric
	for i := range metrics {
		metric := &metrics[i]
		if req.StartTime > 0 && metric.Timestamp < req.StartTime {
			continue
		}
		if req.EndTime > 0 && metric.Timestamp > req.EndTime {
			continue
		}
		filteredMetrics = append(filteredMetrics, metric)
	}

	return &pb.GetMetricsResponse{
		Metrics: filteredMetrics,
	}, nil
}

func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Healthy: true,
		Status:  "healthy",
	}, nil
}
