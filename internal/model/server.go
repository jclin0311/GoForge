package model

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	pb "github.com/goforge/ai-platform/pkg/proto"
	"go.uber.org/zap"
)

type Server struct {
	pb.UnimplementedModelServiceServer
	logger *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Predict(ctx context.Context, req *pb.PredictRequest) (*pb.PredictResponse, error) {
	s.logger.Info("Predict request received",
		zap.String("model_id", req.ModelId),
		zap.String("model_version", req.ModelVersion),
		zap.String("request_id", req.RequestId),
	)

	start := time.Now()

	time.Sleep(time.Duration(50+rand.Intn(200)) * time.Millisecond)

	inferenceTime := time.Since(start).Milliseconds()

	outputData := []byte(fmt.Sprintf("prediction_result_for_%s", req.ModelId))
	confidence := 0.85 + rand.Float32()*0.15

	response := &pb.PredictResponse{
		RequestId:        req.RequestId,
		OutputData:       outputData,
		Confidence:       confidence,
		InferenceTimeMs:  inferenceTime,
		Metadata: map[string]string{
			"model_id":      req.ModelId,
			"model_version": req.ModelVersion,
			"timestamp":     time.Now().Format(time.RFC3339),
		},
	}

	s.logger.Info("Predict request completed",
		zap.String("request_id", req.RequestId),
		zap.Int64("inference_time_ms", inferenceTime),
		zap.Float32("confidence", confidence),
	)

	return response, nil
}

func (s *Server) BatchPredict(ctx context.Context, req *pb.BatchPredictRequest) (*pb.BatchPredictResponse, error) {
	s.logger.Info("BatchPredict request received",
		zap.String("model_id", req.ModelId),
		zap.Int("batch_size", len(req.InputBatch)),
	)

	start := time.Now()
	var predictions []*pb.PredictResponse

	for i, input := range req.InputBatch {
		prediction := &pb.PredictResponse{
			RequestId:        fmt.Sprintf("%s_%d", req.RequestId, i),
			OutputData:       []byte(fmt.Sprintf("batch_prediction_%d", i)),
			Confidence:       0.80 + rand.Float32()*0.20,
			InferenceTimeMs:  int64(30 + rand.Intn(100)),
			Metadata: map[string]string{
				"batch_index": fmt.Sprintf("%d", i),
				"input_size":  fmt.Sprintf("%d", len(input)),
			},
		}
		predictions = append(predictions, prediction)
	}

	totalTime := time.Since(start).Milliseconds()

	return &pb.BatchPredictResponse{
		RequestId:            req.RequestId,
		Predictions:          predictions,
		TotalInferenceTimeMs: totalTime,
	}, nil
}

func (s *Server) GetModelInfo(ctx context.Context, req *pb.ModelInfoRequest) (*pb.ModelInfoResponse, error) {
	s.logger.Info("GetModelInfo request received",
		zap.String("model_id", req.ModelId),
		zap.String("model_version", req.ModelVersion),
	)

	return &pb.ModelInfoResponse{
		ModelId:      req.ModelId,
		ModelVersion: req.ModelVersion,
		ModelType:    "neural-network",
		Framework:    "tensorflow",
		Metadata: map[string]string{
			"architecture": "transformer",
			"parameters":   "1.5B",
			"dataset":      "custom",
		},
		Status: pb.ModelStatus_MODEL_STATUS_READY,
	}, nil
}

func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Healthy: true,
		Status:  "healthy",
		Details: map[string]string{
			"service":   "model-service",
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    "active",
		},
	}, nil
}
