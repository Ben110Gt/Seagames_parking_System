package handler

import (
	"seagame/ticket/backend/internal/models/user"
	"seagame/ticket/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{userService: userService}
}

// Get api/v1/users/:id
func (h *userHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	u, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "User fetched successfully",
		"user":    u,
	})
}

// Get api/v1/users
func (h *userHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch users"})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Users fetched successfully",
		"users":   users,
	})
}

// PUT /api/v1/users
func (h *userHandler) UpdateUser(c *fiber.Ctx) error {
	var req user.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if err := h.userService.UpdateUser(c.Context(), &req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update user"})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

// DELETE /api/v1/users/:id
func (h *userHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if err := h.userService.DeleteUser(c.Context(), id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete user"})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
