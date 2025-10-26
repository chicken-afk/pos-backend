package middlewares

import (
	"context"
	"strings"

	"pos/auth_service/config"
	"pos/shared/jwtlib"

	redisPkg "pos/auth_service/app/pkg/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing token",
			})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwtlib.VerifyToken(tokenStr, "http://localhost:8080/api/.well-known/jwks.json")
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":  "invalid token",
				"detail": err.Error(),
			})
		}

		//Check if token is blacklisted
		isBlacklisted, err := redisPkg.IsTokenBlacklisted(config.Redis, tokenStr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":  "failed to check token blacklist",
				"detail": err.Error(),
			})
		}

		if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token is blacklisted",
			})
		}

		claims := token.Claims.(jwt.MapClaims)
		ctx := context.WithValue(c.Context(), "user_email", claims["user_email"])
		c.SetUserContext(ctx)
		return c.Next()
	}
}
