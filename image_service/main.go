package main

import (
	"log"
	"net"
	"os"
	"pos/image_service/app/handlers"
	services "pos/image_service/app/service"
	"pos/image_service/pb"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startGRPCServer() {
	//Create grpc Handler
	grpcHandler := handlers.NewImageHandler(services.NewImageService())

	//Create grpc server
	grpcServer := grpc.NewServer()

	//Register the image service
	pb.RegisterImageServiceServer(grpcServer, grpcHandler)

	//Enable reflection for testing with tools like grpcurl
	reflection.Register(grpcServer)

	//port
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}

	//Create listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Println("🚀 gRPC Server starting on port", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}

func main() {
	//Init godotenv or config here if needed
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	startGRPCServer()
}
