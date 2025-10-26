package routes

import (
	"pos/auth_service/app/handlers"
	"pos/auth_service/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router, authController handlers.AuthController) {
	v1Router := router.Group("/v1")
	v1Router.Post("/login", authController.Login)
	v1Router.Post("/logout", middlewares.JWTAuth(), authController.Logout)
}
