package handlers

import (
	"context"
	"log"
	"pos/auth_service/app/services"
	"pos/auth_service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCAuthHandler implements the gRPC AuthServiceServer interface
type GRPCAuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService services.AuthService
}

// NewGRPCAuthHandler creates a new gRPC auth handler
func NewGRPCAuthHandler(authService services.AuthService) *GRPCAuthHandler {
	return &GRPCAuthHandler{
		authService: authService,
	}
}

// ValidateToken implements the gRPC ValidateToken method
func (h *GRPCAuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	log.Printf("gRPC ValidateToken called with token: %s", req.Token)

	// Call the existing auth service ValidateToken method
	isValid, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		return &pb.ValidateTokenResponse{
			IsValid:   false,
			ExpiresAt: nil,
		}, status.Errorf(codes.Internal, "failed to validate token: %v", err)
	}

	response := &pb.ValidateTokenResponse{
		IsValid: isValid,
		// Note: ExpiresAt could be added later if needed
		// For now, we'll keep it simple and just return validation status
		ExpiresAt: nil,
	}

	log.Printf("Token validation result: %v", isValid)
	return response, nil
}
