package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"pos/auth_service/app/handlers"
	"pos/auth_service/app/pkg/jwt"
	"pos/auth_service/config"
	"pos/auth_service/pb"
	"pos/auth_service/routes"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	// Initialization code here
	//Load godotenv, setup database, redis, fiber app, routes, etc.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading env file in main.go: ", err)
	}

	//Init JWT Config
	config.InitJWT()

	//Init JWT
	jwt.InitJWT()

	//Connect to DB, Redis, Fiber app, etc.
	config.ConnectDatabase()
	config.InitRedis()
}

func main() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}
	time.Local = loc

	fmt.Println("Auth Service is running...")

	// Init dependencies
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	init := config.NewInitialization(ctx)

	//setup routes and start server
	server := routes.Init(init)

	init.App = server

	//Start GRPC Server
	// startGRPCServer(init)

	// Graceful Shutdown
	// GRACEFUL SHUTDOWN
	go func() {
		config.GracefulShutdown(init, nil, 30*time.Second)
	}()

	//Start server Fiber
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	if err := server.Listen(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

	log.Println("✅ Auth Service stopped")
}

func startGRPCServer(init *config.Initialization) {
	//Create grpc Handler
	grpcHandler := handlers.NewGRPCAuthHandler(init.AuthService)

	//Create grpc server
	grpcServer := grpc.NewServer()

	//Register the auth service
	pb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	//Enable reflection for testing with tools like grpcurl
	reflection.Register(grpcServer)

	//Create listener
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	log.Println("🚀 gRPC Server starting on port 50052...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
