package routes

import (
	"pos/auth_service/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupJWKSRoutes(router fiber.Router, jwksHandler handlers.JwksController) {
	v1Router := router.Group("/.well-known")
	v1Router.Get("/jwks.json", jwksHandler.GetJwks)
}
