package handlers

import (
	"pos/auth_service/app/services"

	"github.com/gofiber/fiber/v2"
)

type JwksController interface {
	GetJwks(c *fiber.Ctx) error
}

type jwksController struct {
	jwksService services.JwksService
}

func NewJwksController(jwksService services.JwksService) JwksController {
	return &jwksController{
		jwksService: jwksService,
	}
}

func (ctrl *jwksController) GetJwks(c *fiber.Ctx) error {
	jwks, err := ctrl.jwksService.GetJwks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve JWKS",
		})
	}
	return c.JSON(jwks)
}
