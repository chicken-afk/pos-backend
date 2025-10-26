package jwt

import (
	"os"
	"pos/auth_service/app/entities"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func InitJWT() {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("⚠️ .env file not found, relying on system environment")
	}
	secretKeyAccessToken = []byte(os.Getenv("JWT_SECRET_ACCESS_TOKEN"))
	expStr := os.Getenv("JWT_EXPIRED_ACCESS_TOKEN") // detik
	val, err := strconv.Atoi(expStr)
	if err != nil || val <= 0 {
		logrus.Warnf("⚠️ JWT_EXPIRED_ACCESS_TOKEN invalid='%s', fallback 3600s", expStr)
		val = 3600
	}
	expiredAccessToken = val
}

var (
	blacklistedTokens    = make(map[string]time.Time)
	userTokens           = make(map[string][]string)
	blacklistMutex       sync.RWMutex
	secretKeyAccessToken []byte
	expiredAccessToken   int
)

// JWTClaims struct: user_email + standard claims
type JWTClaims struct {
	UserEmail string `json:"user_email"`
	SessionID string `json:"session_id"`
	Type      string `json:"typ"`
	jwt.RegisteredClaims
}

// CreateTokenJwks membuat JWT dengan JWKs return accessToken, error
func CreateTokenJwks(user *entities.User, privateKeyJWT []byte, sessionId string) (string, error) {

	var accessToken string

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyJWT)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_email": user.Email,
		"role":       user.Role,
		"outlet":     user.Outlet,
		"exp":        jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		"iss":        "auth-service",
		"typ":        "access",
		"session_id": sessionId,
	})
	token.Header["kid"] = "rsa-key-1" // key identifier

	accessToken, err = token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func CreateRefreshTokenJwks(userEmail string, privateKeyJWT []byte, sessionId string) (string, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyJWT)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_email": userEmail,
		"exp":        jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		"iss":        "auth-service",
		"typ":        "refresh",
		"session_id": sessionId,
	})
	token.Header["kid"] = "rsa-key-1" // key identifier
	refreshToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func ParseTokenJwks(tokenStr string, publicKey []byte) (*JWTClaims, error) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenSignatureInvalid
}

func RemoveToken(tokenStr string) error {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()
	blacklistedTokens[tokenStr] = time.Now()
	return nil
}

func ValidateToken(tokenStr string) (bool, error) {
	blacklistMutex.RLock()
	defer blacklistMutex.RUnlock()
	if _, found := blacklistedTokens[tokenStr]; found {
		return false, nil
	}
	return true, nil
}
