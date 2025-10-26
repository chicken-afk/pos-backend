package services

import (
	"fmt"
	"log"
	"os"
	"pos/auth_service/app/dto/request"
	"pos/auth_service/app/dto/response"
	jwtPkg "pos/auth_service/app/pkg/jwt"
	redisPkg "pos/auth_service/app/pkg/redis"
	"pos/auth_service/app/repositories"
	"pos/auth_service/app/utils"
	"time"

	"pos/shared/jwtlib"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	// Define authentication related methods here
	Login(request request.LoginRequest) (response.LoginResponse, error)
	ValidateToken(token string) (bool, error)
	Logout(tokenStr string, sessionID string) error
	RefreshAccessToken(refreshToken string) (response.LoginResponse, error)
}

type authService struct {
	userRepo      repositories.UserRepository
	redisClient   *redis.Client
	privateKeyJWT []byte
	publicKeyJWT  []byte
}

func NewAuthService(userRepo repositories.UserRepository, redisClient *redis.Client, privateKeyJWT []byte, publicKeyJWT []byte) AuthService {
	return &authService{
		userRepo:      userRepo,
		redisClient:   redisClient,
		privateKeyJWT: privateKeyJWT,
		publicKeyJWT:  publicKeyJWT,
	}
}

func (s *authService) Login(request request.LoginRequest) (response.LoginResponse, error) {
	var loginResp response.LoginResponse

	user, err := s.userRepo.FindByEmail(request.Email)
	if err != nil {
		return loginResp, err
	}

	if !utils.CheckPasswordHash(request.Password, user.Password) {
		return loginResp, fmt.Errorf("invalid credentials")
	}

	sessionID := uuid.New().String()

	accessToken, err := jwtPkg.CreateTokenJwks(user, s.privateKeyJWT, sessionID)
	if err != nil {
		log.Println("Failed to create token:", err)
		return loginResp, err
	}

	refreshToken, err := jwtPkg.CreateRefreshTokenJwks(user.Email, s.privateKeyJWT, sessionID)
	if err != nil {
		log.Println("Failed to create refresh token:", err)
		return loginResp, err
	}

	claims, err := jwtPkg.ParseTokenJwks(refreshToken, s.publicKeyJWT)
	if err != nil {
		log.Println("Failed to parse token:", err)
		return loginResp, err
	}
	log.Printf("Token claims: %+v\n", claims)
	refreshTokenExpiryAt := claims.ExpiresAt.Unix()

	log.Printf("Refresh token expiry (seconds): %d\n", refreshTokenExpiryAt)

	// Store refresh token in Redis
	loginResp.SessionID = sessionID
	err = redisPkg.SetRefreshToken(s.redisClient, redisPkg.SetRefreshTokenParams{
		SessionID:    sessionID,
		RefreshToken: refreshToken,
		Email:        user.Email,
		ExpiryTime:   int(refreshTokenExpiryAt - time.Now().Unix()),
	})
	if err != nil {
		log.Println("Failed to store refresh token in Redis:", err)
		return loginResp, err
	}

	loginResp.User.Email = user.Email
	loginResp.User.Name = user.Name
	loginResp.User.Role = user.Role
	loginResp.User.Outlet = user.Outlet
	loginResp.AccessToken = accessToken
	loginResp.RefreshToken = refreshToken
	loginResp.TokenType = "Bearer"
	loginResp.ExpiresIn = utils.ConvertEpochToDateTimeJakarta(int64(claims.ExpiresAt.Unix()))

	return loginResp, nil
}

func (s *authService) ValidateToken(token string) (bool, error) {
	// isValid, err := jwtPkg.ValidateToken(s.redisClient, token)
	// if err != nil {
	// 	return false, err
	// }
	// return isValid, nil
	return true, nil
}

func (s *authService) Logout(tokenStr string, sessionID string) error {

	//parse token to get expiry time
	claims, err := jwtPkg.ParseTokenJwks(tokenStr, s.publicKeyJWT)
	if err != nil {
		log.Println("Failed to parse token:", err)
		return err
	}
	log.Printf("Token claims: %+v\n", claims)
	refreshTokenExpiryAt := claims.ExpiresAt.Unix()

	log.Printf("Refresh token expiry (seconds): %d\n", refreshTokenExpiryAt)

	err = redisPkg.BlacklistAccessToken(s.redisClient, tokenStr, int(refreshTokenExpiryAt-time.Now().Unix())) // Blacklist for remaining time
	if err != nil {
		log.Println("Failed to blacklist access token:", err)
		return err
	}

	err = redisPkg.RemoveRefreshToken(s.redisClient, sessionID, claims.UserEmail)
	if err != nil {
		log.Println("Failed to remove refresh token from Redis:", err)
		return err
	}

	return nil
}

func (s *authService) RefreshAccessToken(refreshToken string) (response.LoginResponse, error) {
	var loginResp response.LoginResponse
	// ✅ Verifikasi refresh token
	publicJWKSUrl := os.Getenv("PUBLIC_JWKS_URL")
	if publicJWKSUrl == "" {
		publicJWKSUrl = "http://localhost:8080/api/.well-known/jwks.json"
	}
	token, err := jwtlib.VerifyToken(refreshToken, publicJWKSUrl)
	if err != nil || !token.Valid {
		return loginResp, fmt.Errorf("invalid refresh token")
	}
	claims, err := jwtPkg.ParseTokenJwks(refreshToken, s.publicKeyJWT)
	if err != nil {
		log.Println("Failed to parse token:", err)
		return loginResp, err
	}
	email := claims.UserEmail
	sessionID := claims.SessionID
	tokenType := claims.Type

	if tokenType != "refresh" {
		return loginResp, fmt.Errorf("invalid token type")
	}

	// ✅ Cek refresh token di Redis
	storedRefreshToken, err := redisPkg.GetRefreshToken(s.redisClient, sessionID, email)
	if err != nil {
		return loginResp, fmt.Errorf("refresh token not found")
	}
	if storedRefreshToken != refreshToken {
		return loginResp, fmt.Errorf("refresh token mismatch")
	}

	// ✅ Buat access token baru
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return loginResp, err
	}
	accessToken, err := jwtPkg.CreateTokenJwks(user, s.privateKeyJWT, sessionID)
	if err != nil {
		log.Println("Failed to create token:", err)
		return loginResp, err
	}
	loginResp.User.Email = user.Email
	loginResp.User.Name = user.Name
	loginResp.User.Role = user.Role
	loginResp.User.Outlet = user.Outlet
	loginResp.AccessToken = accessToken
	loginResp.RefreshToken = refreshToken
	loginResp.SessionID = sessionID
	loginResp.TokenType = "Bearer"

	return loginResp, nil
}
