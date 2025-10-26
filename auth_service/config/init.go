package config

import (
	"context"
	"log"
	"pos/auth_service/app/repositories"
	"pos/auth_service/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Initialization struct {
	// Add Resources here
	App   *fiber.App
	DB    *gorm.DB
	Redis *redis.Client
	// JWT Public Data
	privateKeyJWT *jwt.SigningMethodRSA
	keyDataJWT    []byte
	publicDataJWT []byte

	// Repositories
	UserRepo repositories.UserRepository

	// Services
	AuthService services.AuthService
	JwksService services.JwksService
}

func NewInitialization(ctx context.Context) *Initialization {
	userRepo := repositories.NewUserRepository(DB)
	jwksService := services.NewJwksService(publicDataJWT)

	authService := services.NewAuthService(userRepo, Redis, keyDataJWT, publicDataJWT)

	return &Initialization{
		UserRepo:    userRepo,
		AuthService: authService,
		JwksService: jwksService,
	}
}

// FUNCTION CLOSE ALL RESOURCES
func (i *Initialization) Close() {
	// CLOSE DB
	if i.DB != nil {
		sqlDB, err := i.DB.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("✅ DB closed")
		}
	}
	// CLOSE REDIS
	if i.Redis != nil {
		if err := i.Redis.Close(); err == nil {
			log.Println("✅ Redis closed")
		}
	}
}
