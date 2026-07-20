package grpc

import (
	"context"

	middleware "laundry-service/middlewares"
	pb "laundry-service/proto" // Import file proto yang dihasilkan
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthServiceServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	isValid, message := middleware.ValidateTokenMiddleware(req.Token)

	return &pb.ValidateTokenResponse{
		IsValid: isValid,
		Message: message,
	}, nil
}
