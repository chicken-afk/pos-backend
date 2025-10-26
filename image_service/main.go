package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"pos/image_service/app/handlers"
	services "pos/image_service/app/service"
	"pos/image_service/pb"

	"github.com/gorilla/mux"
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
		port = "50053" // Different port from auth service
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

func startHTTPServer() {
	// Create HTTP image handler
	httpHandler := handlers.NewHTTPImageHandler()

	// Create router
	r := mux.NewRouter()

	// Route for serving images
	// This will handle URLs like: /storage/images/2025/10/26/testupload_1761478232.jpg
	r.PathPrefix("/storage/").HandlerFunc(httpHandler.ServeImage).Methods("GET")

	// Add CORS middleware if needed
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Get HTTP port
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8082"
	}

	log.Printf("🌐 HTTP Server starting on port %s", httpPort)
	log.Printf("📸 Image serving endpoint: http://localhost:%s/storage/images/...", httpPort)

	if err := http.ListenAndServe(":"+httpPort, r); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func main() {
	//Init godotenv or config here if needed
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	go startGRPCServer() //Start grpc server on go routine
	startHTTPServer()
}
