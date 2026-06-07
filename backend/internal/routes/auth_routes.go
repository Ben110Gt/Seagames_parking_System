package routes

import (
	"seagame/ticket/backend/database"
	"seagame/ticket/backend/internal/handler"
	"seagame/ticket/backend/internal/repository"
	"seagame/ticket/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	userRepo := repository.NewUserRepository(database.GetDB())
	userService := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(userService)

	auth := app.Group("/auth")
	auth.Post("/login", authHandler.Login)       // ✅
	auth.Post("/register", authHandler.Register) // ✅
}
