package handlers

import (
	"pos/auth_service/app/dto/request"
	"pos/auth_service/app/dto/response"
	"pos/auth_service/app/services"

	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (h *authController) Login(c *fiber.Ctx) error {
	var loginReq request.LoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	loginResp, err := h.authService.Login(loginReq)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(loginResp)
}

func (h *authController) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing Authorization header",
		})
	}
	tokenStr := authHeader[len("Bearer "):]

	var logoutReq request.LogoutRequest
	if err := c.BodyParser(&logoutReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	sessionID := logoutReq.SessionID

	if err := h.authService.Logout(tokenStr, sessionID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Failed to logout",
			"detail": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.LogoutResponse{
		Message: "Successfully logged out",
	})
}

func (h *authController) RefreshToken(c *fiber.Ctx) error {
	//Get refresh token from header
	refreshToken := c.Get("Refresh-Token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing Refresh-Token header",
		})
	}
	res, err := h.authService.RefreshAccessToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}
