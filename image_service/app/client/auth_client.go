package client

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient interface for auth service operations
type AuthClient interface {
	ValidateToken(token string) (bool, error)
	Close() error
}

// authClient implements AuthClient
type authClient struct {
	conn *grpc.ClientConn
	// client authpb.AuthServiceClient - will be implemented after copying proto files
}

// NewAuthClient creates a new auth service client
func NewAuthClient(authServiceURL string) (AuthClient, error) {
	if authServiceURL == "" {
		authServiceURL = "localhost:50052" // default auth service port
	}

	conn, err := grpc.NewClient(authServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// client := authpb.NewAuthServiceClient(conn) - will be uncommented after proto setup

	return &authClient{
		conn: conn,
		// client: client,
	}, nil
}

// ValidateToken validates JWT token with auth service
func (c *authClient) ValidateToken(token string) (bool, error) {
	// TODO: Implement actual gRPC call
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// req := &authpb.ValidateTokenRequest{Token: token}
	// resp, err := c.client.ValidateToken(ctx, req)
	// if err != nil {
	//     return false, err
	// }
	// return resp.IsValid, nil

	// Temporary implementation for development
	log.Printf("Validating token (temporary bypass): %s", token)
	return true, nil
}

// Close closes the gRPC connection
func (c *authClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
