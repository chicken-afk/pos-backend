package routes

import (
	"pos/auth_service/app/handlers"
	"pos/auth_service/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Init(init *config.Initialization) *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
	})
	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))

	app.Use(func(c *fiber.Ctx) error {
		// Call next handlers first
		err := c.Next()
		// After handler finishes, add trace id handler
		c.Set("X-Trace-Id", c.Get("x-trace-id", ""))
		return err
	})
	// HOME ROUTE
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "OK",
			"message": "Welcome to POS Auth Service",
			"version": "1.0.0",
		})
	})

	api := app.Group("/api")

	authHandler := handlers.NewAuthController(init.AuthService)
	jwksHandler := handlers.NewJwksController(init.JwksService)

	SetupAuthRoutes(api, authHandler)
	SetupJWKSRoutes(api, jwksHandler)

	return app
}
