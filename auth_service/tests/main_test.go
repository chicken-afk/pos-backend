package tests

import (
	"log"
	"pos/auth_service/app/pkg/jwt"
	"pos/auth_service/config"
	"testing"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var testRedis *redis.Client
var grpcConn *grpc.ClientConn

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	jwt.InitJWT()

	config.ConnectDatabase()
	testDB = config.DB
	config.InitRedis()
	testRedis = config.Redis
	_, err := createGRPCConnection()
	if err != nil {
		log.Fatal("Error creating gRPC client:", err)
	}
	m.Run()
}

func createGRPCConnection() (*grpc.ClientConn, error) {
	address := "localhost:50052"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	grpcConn = conn
	return conn, nil
}
