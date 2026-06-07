package handler

import (
	"seagame/ticket/backend/internal/models/user"
	"seagame/ticket/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *authHandler {
	return &authHandler{userService: userService}
}

// Post /api/auth/login
func (h *authHandler) Login(c *fiber.Ctx) error {
	var req user.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	user, err := h.userService.Login(c.Context(), &req)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "User logged in successfully",
		"user":    user,
	})
}

// Post /api/auth/register
func (h *authHandler) Register(c *fiber.Ctx) error {
	var req user.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	user, err := h.userService.Register(c.Context(), &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}
