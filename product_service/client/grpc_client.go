package client

import (
	"context"
	"pos/grpc/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn            *grpc.ClientConn
	imageGrpcClient pb.ImageServiceClient
}

func NewGRPCClient(address string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewImageServiceClient(conn)

	return &GRPCClient{
		conn:            conn,
		imageGrpcClient: client,
	}, nil
}

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

// UploadImage uploads an image via gRPC
func (c *GRPCClient) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.imageGrpcClient.UploadImage(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.ImageUrl, nil
}
