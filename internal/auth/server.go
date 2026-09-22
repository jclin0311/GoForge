package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	pb "github.com/goforge/ai-platform/pkg/proto"
	"go.uber.org/zap"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	logger *zap.Logger
	tokens map[string]*TokenInfo
}

type TokenInfo struct {
	UserID    string
	Scopes    []string
	ExpiresAt int64
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{
		logger: logger,
		tokens: make(map[string]*TokenInfo),
	}
}

func (s *Server) Authenticate(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	s.logger.Info("Authentication request received",
		zap.String("username", req.Username),
		zap.String("client_id", req.ClientId),
	)

	if req.Username == "" || req.Password == "" {
		return &pb.AuthResponse{
			Success: false,
		}, fmt.Errorf("invalid credentials")
	}

	accessToken, err := generateToken()
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken()
	if err != nil {
		return nil, err
	}

	expiresIn := int64(3600)
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second).Unix()

	s.tokens[accessToken] = &TokenInfo{
		UserID:    req.Username,
		Scopes:    []string{"read", "write"},
		ExpiresAt: expiresAt,
	}

	return &pb.AuthResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		Scopes:       []string{"read", "write"},
	}, nil
}

func (s *Server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	s.logger.Debug("Token validation request", zap.String("token", req.Token[:10]+"..."))

	tokenInfo, exists := s.tokens[req.Token]
	if !exists {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	if time.Now().Unix() > tokenInfo.ExpiresAt {
		delete(s.tokens, req.Token)
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:     true,
		UserId:    tokenInfo.UserID,
		Scopes:    tokenInfo.Scopes,
		ExpiresAt: tokenInfo.ExpiresAt,
	}, nil
}

func (s *Server) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	s.logger.Info("Token refresh request received")

	newAccessToken, err := generateToken()
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := generateToken()
	if err != nil {
		return nil, err
	}

	expiresIn := int64(3600)
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second).Unix()

	s.tokens[newAccessToken] = &TokenInfo{
		UserID:    "refreshed_user",
		Scopes:    []string{"read", "write"},
		ExpiresAt: expiresAt,
	}

	return &pb.RefreshTokenResponse{
		Success:      true,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *Server) RevokeToken(ctx context.Context, req *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	s.logger.Info("Token revocation request received")

	delete(s.tokens, req.Token)

	return &pb.RevokeTokenResponse{
		Success: true,
	}, nil
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
